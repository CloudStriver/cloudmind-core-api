package service

import (
	"context"
	"fmt"
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/consts"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/kq"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_content"
	platformservice "github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/platform"
	"github.com/CloudStriver/cloudmind-mq/app/util/message"
	"github.com/CloudStriver/go-pkg/utils/pconvertor"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/platform"
	"github.com/bytedance/sonic"
	"github.com/google/wire"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/mr"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type IRelationDomainService interface {
	CreateRelation(ctx context.Context, r *core_api.Relation) (err error)
	GetUserByRelations(ctx context.Context, relations []*platform.Relation, users []*core_api.User, userId string, query string) (err error)
	GetPostByRelations(ctx context.Context, relations []*platform.Relation, posts []*core_api.Post, userId string) (err error)
}
type RelationDomainService struct {
	Config               *config.Config
	Platform             platformservice.IPlatForm
	CloudMindContent     cloudmind_content.ICloudMindContent
	UserDomainService    IUserDomainService
	PostDomainService    IPostDomainService
	CreateNotificationKq *kq.CreateNotificationsKq
	CreateFeedBackKq     *kq.CreateFeedBackKq
	Redis                *redis.Redis
}

func (s *RelationDomainService) GetUserByRelations(ctx context.Context, relations []*platform.Relation, users []*core_api.User, userId string, query string) (err error) {
	err = mr.Finish(lo.Map[*platform.Relation](relations, func(r *platform.Relation, i int) func() error {
		return func() error {
			users[i] = &core_api.User{}
			if query == "from" {
				users[i].UserId = r.ToId
			} else {
				users[i].UserId = r.FromId
			}
			if err = mr.Finish(func() error {
				user, err := s.CloudMindContent.GetUser(ctx, &content.GetUserReq{
					UserId: users[i].UserId,
				})
				if err != nil {
					return err
				}
				users[i].Name = user.Name
				users[i].Url = user.Url
				users[i].Labels = user.Labels
				users[i].Description = user.Description
				s.UserDomainService.LoadLabel(ctx, users[i].Labels)
				return nil
			}, func() error {
				if userId != "" && userId != users[i].UserId {
					s.UserDomainService.LoadFollowed(ctx, &users[i].Followed, userId, users[i].UserId)
				}
				return nil
			}, func() error {
				s.UserDomainService.LoadFollowedCount(ctx, &users[i].FollowedCount, users[i].UserId)
				return nil
			}); err != nil {
				return err
			}
			return nil
		}
	})...)
	return nil
}

func (s *RelationDomainService) GetPostByRelations(ctx context.Context, relations []*platform.Relation, posts []*core_api.Post, userId string) (err error) {
	if err = mr.Finish(lo.Map[*platform.Relation](relations, func(relation *platform.Relation, i int) func() error {
		return func() error {
			posts[i] = &core_api.Post{}
			if err = mr.Finish(func() error {
				post, err1 := s.CloudMindContent.GetPost(ctx, &content.GetPostReq{
					PostId: relation.ToId,
				})
				if err1 != nil {
					return err1
				}

				posts[i].PostId = relation.ToId
				posts[i].Title = post.Title
				posts[i].Text = post.Text
				posts[i].Url = post.Url
				s.PostDomainService.LoadLabels(ctx, &posts[i].Labels, post.LabelIds)
				user, err1 := s.CloudMindContent.GetUser(ctx, &content.GetUserReq{
					UserId: post.UserId,
				})
				if err1 != nil {
					return err1
				}
				posts[i].UserName = user.Name
				return nil
			}, func() error {
				s.PostDomainService.LoadLikeCount(ctx, &posts[i].LikeCount, relation.ToId)
				return nil
			}, func() error {
				if userId != "" {
					s.PostDomainService.LoadLiked(ctx, &posts[i].Liked, userId, relation.ToId)
				}
				return nil
			}, func() error {
				getCommentSubjectResp, err2 := s.Platform.GetCommentSubject(ctx, &platform.GetCommentSubjectReq{SubjectId: relation.ToId})
				if err2 != nil {
					return err2
				}
				posts[i].CommentCount = getCommentSubjectResp.AllCount
				return nil
			}, func() error {
				posts[i].ViewCount, _ = s.Redis.PfcountCtx(ctx, fmt.Sprintf("%s:%s", consts.ViewCountKey, relation.ToId))
				return nil
			}); err != nil {
				return err
			}
			return nil
		}
	})...); err != nil {
		return err
	}
	return nil
}

var RelationDomainServiceSet = wire.NewSet(
	wire.Struct(new(RelationDomainService), "*"),
	wire.Bind(new(IRelationDomainService), new(*RelationDomainService)),
)

func (s *RelationDomainService) CreateRelation(ctx context.Context, r *core_api.Relation) (err error) {
	if r.ToId == r.FromId && r.RelationType == core_api.RelationType_FollowRelationType {
		return nil
	}

	ok, err := s.Platform.CreateRelation(ctx, &platform.CreateRelationReq{
		FromType:     int64(r.FromType),
		FromId:       r.FromId,
		ToType:       int64(r.ToType),
		ToId:         r.ToId,
		RelationType: int64(r.RelationType),
	})
	if err != nil {
		return err
	}
	if !ok.Ok {
		return nil
	}

	act := r.RelationType
	if r.RelationType == core_api.RelationType_HateRelationType {
		act = core_api.RelationType(content.Action_HateType)
	}

	userId := ""
	toName := ""
	subjectId := ""
	var reqs *content.IncrHotValueReq
	switch r.ToType {
	case core_api.TargetType_UserType:
		reqs = &content.IncrHotValueReq{
			Action:     content.Action(act),
			HotId:      r.ToId,
			TargetType: content.TargetType_UserType,
		}
	case core_api.TargetType_FileType:
		reqs = &content.IncrHotValueReq{
			Action:     content.Action(act),
			HotId:      r.ToId,
			TargetType: content.TargetType_FileType,
		}
	case core_api.TargetType_PostType:
		reqs = &content.IncrHotValueReq{
			Action:     content.Action(act),
			HotId:      r.ToId,
			TargetType: content.TargetType_PostType,
		}
	case core_api.TargetType_CommentContentType:
		reqs = &content.IncrHotValueReq{
			Action:     content.Action(act),
			HotId:      r.ToId,
			TargetType: content.TargetType_CommentContentType,
		}
	}

	if _, err = s.CloudMindContent.IncrHotValue(ctx, reqs); err != nil {
		return err
	}
	// 发布，上传，浏览不需要通知
	if r.RelationType == core_api.RelationType_PublishRelationType || r.RelationType == core_api.RelationType_UploadRelationType || r.RelationType == core_api.RelationType_ViewRelationType {
		return nil
	}

	switch r.ToType {
	case core_api.TargetType_UserType:
		userId = r.ToId
	case core_api.TargetType_FileType:
		getFileResp, err := s.CloudMindContent.GetFile(ctx, &content.GetFileReq{
			FileId: r.ToId,
		})
		if err != nil {
			return err
		}

		toName = getFileResp.Name
		userId = getFileResp.UserId

	case core_api.TargetType_PostType:
		getPostResp, err := s.CloudMindContent.GetPost(ctx, &content.GetPostReq{
			PostId: r.ToId,
		})
		if err != nil {
			return err
		}
		toName = getPostResp.Title
		userId = getPostResp.UserId

	case core_api.TargetType_CommentContentType:
		getCommentResp, err := s.Platform.GetComment(ctx, &platform.GetCommentReq{
			CommentId: r.ToId,
		})
		if err != nil {
			return err
		}
		subjectId = getCommentResp.SubjectId
		toName = getCommentResp.Content
		userId = getCommentResp.AtUserId
	}

	userinfo, err := s.CloudMindContent.GetUser(ctx, &content.GetUserReq{
		UserId: r.FromId,
	})
	if err != nil {
		return err
	}

	if r.RelationType == core_api.RelationType_LikeRelationType || r.RelationType == core_api.RelationType_CollectRelationType || r.RelationType == core_api.RelationType_CommentRelationType {
		if _, err = s.CloudMindContent.IncrHotValue(ctx, &content.IncrHotValueReq{
			Action:     content.Action(r.RelationType),
			HotId:      userId,
			TargetType: content.TargetType_UserType,
		}); err != nil {
			return err
		}
		if r.RelationType == core_api.RelationType_CommentRelationType {
			if _, err = s.CloudMindContent.IncrHotValue(ctx, &content.IncrHotValueReq{
				Action:     content.Action(r.RelationType),
				HotId:      subjectId,
				TargetType: content.TargetType_PostType,
			}); err != nil {
				return err
			}
		}
	}
	if r.ToType == core_api.TargetType_CommentContentType && r.RelationType == core_api.RelationType_LikeRelationType {
		return nil
	}
	// 创建通知
	msg, _ := sonic.Marshal(Msg{
		FromName: userinfo.Name,
		ToName:   toName,
	})
	data, _ := sonic.Marshal(&message.CreateNotificationMessage{
		TargetUserId:    userId,
		SourceUserId:    r.FromId,
		SourceContentId: r.ToId,
		TargetType:      int64(r.ToType),
		Type:            int64(r.RelationType),
		Text:            pconvertor.Bytes2String(msg),
	})
	if err = s.CreateNotificationKq.Push(pconvertor.Bytes2String(data)); err != nil {
		return err
	}

	data, _ = sonic.Marshal(&message.CreateFeedBackMessage{
		FeedbackType: core_api.RelationType_name[int32(r.RelationType)],
		UserId:       r.FromId,
		ItemId:       r.ToId,
	})
	if err = s.CreateFeedBackKq.Push(pconvertor.Bytes2String(data)); err != nil {
		return err
	}

	return nil
}

type Msg struct {
	FromName string
	ToName   string
}
