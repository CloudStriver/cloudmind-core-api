package service

import (
	"context"
	"github.com/CloudStriver/cloudmind-core-api/biz/adaptor"
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/cloudmind-core-api/biz/domain/service"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/consts"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/kq"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_content"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/platform_relation"
	"github.com/CloudStriver/cloudmind-mq/app/util/message"
	"github.com/CloudStriver/go-pkg/utils/pconvertor"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/basic"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/platform/relation"
	"github.com/bytedance/sonic"
	"github.com/google/wire"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/mr"
)

type IRelationService interface {
	GetRelation(ctx context.Context, req *core_api.GetRelationReq) (resp *core_api.GetRelationResp, err error)
	CreateRelation(ctx context.Context, req *core_api.CreateRelationReq) (resp *core_api.CreateRelationResp, err error)
	GetFromRelations(ctx context.Context, c *core_api.GetFromRelationsReq) (resp *core_api.GetFromRelationsResp, err error)
	GetToRelations(ctx context.Context, c *core_api.GetToRelationsReq) (resp *core_api.GetToRelationsResp, err error)
	DeleteRelation(ctx context.Context, c *core_api.DeleteRelationReq) (resp *core_api.DeleteRelationResp, err error)
}

var RelationServiceSet = wire.NewSet(
	wire.Struct(new(RelationService), "*"),
	wire.Bind(new(IRelationService), new(*RelationService)),
)

type RelationService struct {
	Config               *config.Config
	PlatFormRelation     platform_relation.IPlatFormRelation
	CloudMindContent     cloudmind_content.ICloudMindContent
	PostDomainService    service.IPostDomainService
	CreateNotificationKq *kq.CreateNotificationsKq
	CreateFeedBackKq     *kq.CreateFeedBackKq
}

func (s *RelationService) GetFromRelations(ctx context.Context, req *core_api.GetFromRelationsReq) (resp *core_api.GetFromRelationsResp, err error) {
	resp = new(core_api.GetFromRelationsResp)
	getFromRelationsResp, err := s.PlatFormRelation.GetRelations(ctx, &relation.GetRelationsReq{
		RelationFilterOptions: &relation.GetRelationsReq_FromFilterOptions{
			FromFilterOptions: &relation.FromFilterOptions{
				FromType: req.FromType,
				FromId:   req.FromId,
				ToType:   int64(req.ToType),
			},
		},
		RelationType: req.RelationType,
		PaginationOptions: &basic.PaginationOptions{
			Limit:  req.Limit,
			Offset: req.Offset,
		},
	})
	if err != nil {
		return resp, err
	}

	switch req.ToType {
	case core_api.TargetType_UserType:
		resp.Users = make([]*core_api.User, len(getFromRelationsResp.Relations))
		err = mr.Finish(lo.Map[*relation.Relation](getFromRelationsResp.Relations, func(relation *relation.Relation, i int) func() error {
			return func() error {
				user, err := s.CloudMindContent.GetUser(ctx, &content.GetUserReq{
					UserId: relation.ToId,
				})
				if err != nil {
					return err
				}

				resp.Users[i] = &core_api.User{
					UserId: relation.ToId,
					Name:   user.Name,
					Url:    user.Url,
				}
				return nil
			}
		})...)
		if err != nil {
			return resp, err
		}
	case core_api.TargetType_PostType:
		resp.Posts = make([]*core_api.Post, len(getFromRelationsResp.Relations))
		if err = mr.Finish(lo.Map[*relation.Relation](getFromRelationsResp.Relations, func(relation *relation.Relation, i int) func() error {
			return func() error {
				resp.Posts[i] = &core_api.Post{}
				if err = mr.Finish(func() error {
					post, err := s.CloudMindContent.GetPost(ctx, &content.GetPostReq{
						PostId: relation.ToId,
					})
					resp.Posts[i].PostId = relation.ToId
					resp.Posts[i].Title = post.Title
					resp.Posts[i].Text = post.Text
					resp.Posts[i].Url = post.Url
					resp.Posts[i].Tags = post.Tags
					return err
				}, func() error {
					s.PostDomainService.LoadLikeCount(ctx, &resp.Posts[i].LikeCount, relation.ToId)
					return nil
				}); err != nil {
					return err
				}
				return nil
			}
		})...); err != nil {
			return resp, err
		}
	default:
		//return resp, consts.ErrNotSupportTargetType
	}
	return resp, nil
}

func (s *RelationService) GetToRelations(ctx context.Context, req *core_api.GetToRelationsReq) (resp *core_api.GetToRelationsResp, err error) {
	resp = new(core_api.GetToRelationsResp)
	getFromRelationsResp, err := s.PlatFormRelation.GetRelations(ctx, &relation.GetRelationsReq{
		RelationFilterOptions: &relation.GetRelationsReq_ToFilterOptions{
			ToFilterOptions: &relation.ToFilterOptions{
				ToType:   req.ToType,
				ToId:     req.ToId,
				FromType: int64(req.FromType),
			},
		},
		RelationType: req.RelationType,
		PaginationOptions: &basic.PaginationOptions{
			Limit:  req.Limit,
			Offset: req.Offset,
		},
	})
	if err != nil {
		return resp, err
	}

	switch req.FromType {
	case core_api.TargetType_UserType:
		resp.Users = make([]*core_api.User, len(getFromRelationsResp.Relations))
		err = mr.Finish(lo.Map[*relation.Relation](getFromRelationsResp.Relations, func(relation *relation.Relation, i int) func() error {
			return func() error {
				user, err := s.CloudMindContent.GetUser(ctx, &content.GetUserReq{
					UserId: relation.FromId,
				})
				if err != nil {
					return err
				}
				resp.Users[i] = &core_api.User{
					UserId: relation.FromId,
					Name:   user.Name,
					Url:    user.Url,
				}
				return nil
			}
		})...)
		if err != nil {
			return resp, err
		}
	default:
		//return resp, consts.ErrNotSupportTargetType
	}
	return resp, nil
}

func (s *RelationService) DeleteRelation(ctx context.Context, req *core_api.DeleteRelationReq) (resp *core_api.DeleteRelationResp, err error) {
	user := adaptor.ExtractUserMeta(ctx)
	if user.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	if _, err = s.PlatFormRelation.DeleteRelation(ctx, &relation.DeleteRelationReq{
		FromType:     int64(core_api.TargetType_UserType),
		FromId:       user.UserId,
		ToType:       req.ToType,
		ToId:         req.ToId,
		RelationType: req.RelationType,
	}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *RelationService) GetRelation(ctx context.Context, req *core_api.GetRelationReq) (resp *core_api.GetRelationResp, err error) {
	resp = new(core_api.GetRelationResp)
	getRelationResp, err := s.PlatFormRelation.GetRelation(ctx, &relation.GetRelationReq{
		FromType:     req.FromType,
		FromId:       req.FromId,
		ToType:       req.ToType,
		ToId:         req.ToId,
		RelationType: req.RelationType,
	})
	if err != nil {
		return resp, err
	}

	resp.Ok = getRelationResp.Ok
	return resp, nil
}

func (s *RelationService) CreateRelation(ctx context.Context, req *core_api.CreateRelationReq) (resp *core_api.CreateRelationResp, err error) {
	user := adaptor.ExtractUserMeta(ctx)
	if user.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	ok, err := s.PlatFormRelation.CreateRelation(ctx, &relation.CreateRelationReq{
		FromType:     int64(core_api.TargetType_UserType),
		FromId:       user.UserId,
		ToType:       int64(req.ToType),
		ToId:         req.ToId,
		RelationType: int64(req.RelationType),
	})
	if err != nil {
		return resp, err
	}

	if !ok.Ok {
		return resp, nil
	}

	if req.RelationType == core_api.RelationType_HateType {
		req.RelationType = core_api.RelationType(content.Action_HateType)
	}

	userId := ""
	toName := ""
	var reqs *content.IncrHotValueReq
	switch req.ToType {
	case core_api.TargetType_UserType:
		userId = req.ToId
		reqs = &content.IncrHotValueReq{
			Action:     content.Action(req.RelationType),
			HotId:      req.ToId,
			TargetType: content.TargetType_UserType,
		}
	case core_api.TargetType_FileType:
		getFileResp, err := s.CloudMindContent.GetFile(ctx, &content.GetFileReq{
			FileId: req.ToId,
		})
		if err != nil {
			return resp, err
		}

		toName = getFileResp.File.Name
		userId = getFileResp.File.UserId

		reqs = &content.IncrHotValueReq{
			Action:     content.Action(req.RelationType),
			HotId:      req.ToId,
			TargetType: content.TargetType_FileType,
		}
	case core_api.TargetType_PostType:
		getPostResp, err := s.CloudMindContent.GetPost(ctx, &content.GetPostReq{
			PostId: req.ToId,
		})
		if err != nil {
			return resp, err
		}
		toName = getPostResp.Title
		userId = getPostResp.UserId
		reqs = &content.IncrHotValueReq{
			Action:     content.Action(req.RelationType),
			HotId:      req.ToId,
			TargetType: content.TargetType_PostType,
		}
	}
	if _, err = s.CloudMindContent.IncrHotValue(ctx, reqs); err != nil {
		return resp, err
	}

	userinfo, err := s.CloudMindContent.GetUser(ctx, &content.GetUserReq{
		UserId: user.UserId,
	})
	if err != nil {
		return resp, err
	}

	// 创建通知
	msg, _ := sonic.Marshal(&Msg{
		FromName: userinfo.Name,
		ToName:   toName,
	})
	data, _ := sonic.Marshal(&message.CreateNotificationMessage{
		TargetUserId:    userId,
		SourceUserId:    user.UserId,
		SourceContentId: req.ToId,
		TargetType:      int64(req.ToType),
		Type:            int64(req.RelationType),
		Text:            pconvertor.Bytes2String(msg),
	})
	if err = s.CreateNotificationKq.Push(pconvertor.Bytes2String(data)); err != nil {
		return resp, err
	}

	data, _ = sonic.Marshal(&message.CreateFeedBackMessage{
		FeedbackType: core_api.RelationType_name[int32(req.RelationType)],
		UserId:       user.UserId,
		ItemId:       req.ToId,
	})
	if err = s.CreateFeedBackKq.Push(pconvertor.Bytes2String(data)); err != nil {
		return resp, err
	}

	// 更新热度

	return resp, nil
}

type Msg struct {
	FromName string
	ToName   string
}
