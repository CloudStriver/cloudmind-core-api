package service

import (
	"context"
	"fmt"
	"github.com/CloudStriver/cloudmind-core-api/biz/adaptor"
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/cloudmind-core-api/biz/domain/service"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/consts"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/kq"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_content"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_sts"
	platformservice "github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/platform"
	"github.com/CloudStriver/cloudmind-mq/app/util/message"
	"github.com/CloudStriver/go-pkg/utils/pconvertor"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/basic"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/sts"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/platform"
	"github.com/bytedance/sonic"
	"github.com/google/wire"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/mr"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type IPostService interface {
	CreatePost(ctx context.Context, req *core_api.CreatePostReq) (resp *core_api.CreatePostResp, err error)
	UpdatePost(ctx context.Context, req *core_api.UpdatePostReq) (resp *core_api.UpdatePostResp, err error)
	DeletePost(ctx context.Context, req *core_api.DeletePostReq) (resp *core_api.DeletePostResp, err error)
	GetPosts(ctx context.Context, c *core_api.GetPostsReq) (*core_api.GetPostsResp, error)
	GetPost(ctx context.Context, c *core_api.GetPostReq) (*core_api.GetPostResp, error)
}

var PostServiceSet = wire.NewSet(
	wire.Struct(new(PostService), "*"),
	wire.Bind(new(IPostService), new(*PostService)),
)

type PostService struct {
	Config                *config.Config
	CloudMindContent      cloudmind_content.ICloudMindContent
	PostDomainService     service.IPostDomainService
	Platform              platformservice.IPlatForm
	CloudMindSts          cloudmind_sts.ICloudMindSts
	RelationDomainService service.IRelationDomainService
	UserDomainService     service.IUserDomainService
	CreateItemKq          *kq.CreateItemKq
	UpdateItemKq          *kq.UpdateItemKq
	DeleteItemKq          *kq.DeleteItemKq
	Redis                 *redis.Redis
}

func (s *PostService) FilterContent(ctx context.Context, IsSure bool, contents []*string) ([]string, error) {
	cts := lo.Map[*string, string](contents, func(item *string, index int) string {
		return *item
	})
	if IsSure {
		replaceContentResp, err := s.CloudMindSts.ReplaceContent(ctx, &sts.ReplaceContentReq{
			Contents: cts,
		})
		if err != nil {
			return nil, err
		}
		for i, c := range replaceContentResp.Content {
			*contents[i] = c
		}
		return nil, nil
	} else {
		// 内容检测
		findAllContentResp, err := s.CloudMindSts.FindAllContent(ctx, &sts.FindAllContentReq{
			Contents: cts,
		})
		if err != nil {
			return nil, err
		}

		return findAllContentResp.Keywords, nil
	}
}

func (s *PostService) CreatePost(ctx context.Context, req *core_api.CreatePostReq) (resp *core_api.CreatePostResp, err error) {
	resp = new(core_api.CreatePostResp)
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	if req.Status == int64(core_api.PostStatus_PublicPostStatus) {
		resp.Keywords, err = s.FilterContent(ctx, req.IsSure, []*string{&req.Title, &req.Text})
		if err != nil {
			return resp, err
		}
		if resp.Keywords != nil {
			return resp, nil
		}
	}

	createPostResp, err := s.CloudMindContent.CreatePost(ctx, &content.CreatePostReq{
		UserId:   userData.UserId,
		Title:    req.Title,
		Text:     req.Text,
		LabelIds: req.LabelIds,
		Status:   req.Status,
		Url:      req.Url,
	})
	if err != nil {
		return resp, err
	}
	resp.PostId = createPostResp.PostId

	if err = mr.Finish(func() error {
		if _, err1 := s.Platform.CreateCommentSubject(ctx, &platform.CreateCommentSubjectReq{
			Id:     createPostResp.PostId,
			UserId: userData.UserId,
		}); err1 != nil {
			return err1
		}
		return nil
	}, func() error {
		if _, err2 := s.CloudMindContent.CreateHot(ctx, &content.CreateHotReq{
			HotId: createPostResp.PostId,
		}); err2 != nil {
			return err2
		}
		return nil
	}, func() error {
		data, _ := sonic.Marshal(&message.CreateItemMessage{
			ItemId:   createPostResp.PostId,
			IsHidden: req.Status == int64(core_api.PostStatus_DraftPostStatus),
			Labels:   req.LabelIds,
			Category: core_api.Category_name[int32(core_api.Category_PostCategory)],
		})
		if err3 := s.CreateItemKq.Push(pconvertor.Bytes2String(data)); err != nil {
			return err3
		}
		return nil
	}, func() error {
		if req.Status == int64(core_api.PostStatus_PublicPostStatus) {
			if err4 := s.RelationDomainService.CreateRelation(ctx, &core_api.Relation{
				FromType:     core_api.TargetType_UserType,
				FromId:       userData.UserId,
				ToType:       core_api.TargetType_PostType,
				ToId:         createPostResp.PostId,
				RelationType: core_api.RelationType_PublishRelationType,
			}); err4 != nil {
				return err4
			}
		}
		return nil
	}); err != nil {
		return resp, err
	}

	return resp, nil
}

func (s *PostService) UpdatePost(ctx context.Context, req *core_api.UpdatePostReq) (resp *core_api.UpdatePostResp, err error) {
	resp = new(core_api.UpdatePostResp)
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	post, err := s.CheckIsMyPost(ctx, req.PostId, userData.UserId)
	if err != nil {
		return resp, err
	}
	if post.Status == int64(core_api.PostStatus_PublicPostStatus) || req.Status == int64(core_api.PostStatus_PublicPostStatus) {
		if req.Title == "" {
			req.Title = post.Title
		}
		if req.Text == "" {
			req.Text = post.Text
		}
		resp.Keywords, err = s.FilterContent(ctx, req.IsSure, []*string{&req.Title, &req.Text})
		if err != nil {
			return resp, err
		}
		if resp.Keywords != nil {
			return resp, nil
		}
	}

	if _, err = s.CloudMindContent.UpdatePost(ctx, &content.UpdatePostReq{
		PostId:   req.PostId,
		Title:    req.Title,
		Text:     req.Text,
		LabelIds: req.LabelIds,
		Status:   req.Status,
		Url:      req.Url,
	}); err != nil {
		return resp, err
	}

	if req.Status != 0 {
		if err = mr.Finish(func() error {
			data, _ := sonic.Marshal(&message.UpdateItemMessage{
				ItemId:   req.PostId,
				IsHidden: lo.ToPtr(req.Status != int64(core_api.PostStatus_PublicPostStatus)),
				Labels:   req.LabelIds,
			})
			if err1 := s.UpdateItemKq.Push(pconvertor.Bytes2String(data)); err != nil {
				return err1
			}
			return nil
		}, func() error {
			if err2 := s.RelationDomainService.CreateRelation(ctx, &core_api.Relation{
				FromType:     core_api.TargetType_UserType,
				FromId:       userData.UserId,
				ToType:       core_api.TargetType_PostType,
				ToId:         req.PostId,
				RelationType: core_api.RelationType_PublishRelationType,
			}); err2 != nil {
				return err2
			}
			return nil
		}); err != nil {
			return resp, err
		}
	}
	return resp, nil
}

func (s *PostService) DeletePost(ctx context.Context, req *core_api.DeletePostReq) (resp *core_api.DeletePostResp, err error) {
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	// 只能删除自己的帖子
	if err = s.CheckIsMyPosts(ctx, req.PostIds, userData.UserId); err != nil {
		return resp, err
	}

	if _, err = s.CloudMindContent.DeletePost(ctx, &content.DeletePostReq{
		PostId: req.PostIds,
	}); err != nil {
		return resp, err
	}
	if err = mr.Finish(lo.Map(req.PostIds, func(item string, index int) func() (err error) {
		return func() (err error) {
			if err = mr.Finish(func() error {
				data, _ := sonic.Marshal(&message.DeleteItemMessage{
					ItemId: item,
				})
				if err1 := s.DeleteItemKq.Push(pconvertor.Bytes2String(data)); err != nil {
					return err1
				}
				return nil
			}, func() error {
				if _, err2 := s.Platform.DeleteRelation(ctx, &platform.DeleteRelationReq{
					FromType:     int64(core_api.TargetType_UserType),
					FromId:       userData.UserId,
					ToType:       int64(core_api.TargetType_PostType),
					ToId:         item,
					RelationType: int64(core_api.RelationType_PublishRelationType),
				}); err2 != nil {
					return err2
				}
				return nil
			}); err != nil {
				return err
			}
			return nil
		}
	})...); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *PostService) GetPost(ctx context.Context, req *core_api.GetPostReq) (resp *core_api.GetPostResp, err error) {
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil {
		return resp, consts.ErrNotAuthentication
	}
	var res *content.GetPostResp
	if res, err = s.CloudMindContent.GetPost(ctx, &content.GetPostReq{
		PostId: req.PostId,
	}); err != nil {
		return resp, err
	}

	// 如果该帖子非公开，并且不是他的，那么他没有权限查看
	if res.Status != int64(core_api.PostStatus_PublicPostStatus) && res.UserId != userData.GetUserId() {
		return resp, consts.ErrForbidden
	}

	_, _ = s.Redis.Pfadd(fmt.Sprintf("%s:%s", consts.ViewCountKey, req.PostId), 1)

	resp = &core_api.GetPostResp{
		Title:  res.Title,
		Text:   res.Text,
		Status: res.Status,
		Url:    res.Url,
		Author: &core_api.PostUser{
			UserId: res.UserId,
		},
		Labels:     make([]*core_api.Label, 0),
		CreateTime: res.CreateTime,
		UpdateTime: res.UpdateTime,
	}

	if err = mr.Finish(func() error {
		s.PostDomainService.LoadAuthor(ctx, resp.Author, userData.GetUserId()) // 作者
		return nil
	}, func() error {
		s.PostDomainService.LoadLikeCount(ctx, &resp.LikeCount, req.PostId) // 点赞量
		return nil
	}, func() error {
		s.PostDomainService.LoadCollectCount(ctx, &resp.CollectCount, req.PostId) // 收藏量
		return nil
	}, func() error {
		s.PostDomainService.LoadShareCount(ctx, &resp.ShareCount, req.PostId) // 分享量
		return nil
	}, func() error {
		s.PostDomainService.LoadLiked(ctx, &resp.Liked, userData.GetUserId(), req.PostId) // 是否点赞
		return nil
	}, func() error {
		s.PostDomainService.LoadCollected(ctx, &resp.Collected, userData.GetUserId(), req.PostId) // 是否收藏
		return nil
	}, func() error {
		s.PostDomainService.LoadCommentCount(ctx, &resp.CommentCount, req.PostId)
		return nil
	}, func() error {
		if userData.GetUserId() != "" {
			_ = s.RelationDomainService.CreateRelation(ctx, &core_api.Relation{
				FromType:     core_api.TargetType_UserType,
				FromId:       userData.UserId,
				ToType:       core_api.TargetType_PostType,
				ToId:         req.PostId,
				RelationType: core_api.RelationType_ViewRelationType,
			})
		}
		return nil
	}, func() error {
		s.PostDomainService.LoadLabels(ctx, &resp.Labels, res.LabelIds)
		return nil
	}, func() error {
		resp.ViewCount, _ = s.Redis.PfcountCtx(ctx, fmt.Sprintf("%s:%s", consts.ViewCountKey, req.PostId)) // 浏览量级
		return nil
	}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *PostService) GetPosts(ctx context.Context, req *core_api.GetPostsReq) (resp *core_api.GetPostsResp, err error) {
	resp = new(core_api.GetPostsResp)
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil {
		return resp, consts.ErrNotAuthentication
	}
	var (
		getPostsResp *content.GetPostsResp
	)

	filter := &content.PostFilterOptions{
		OnlyUserId:  req.OnlyUserId,
		OnlyLabelId: req.OnlyLabelId,
	}

	// 查看所有人的，或者查看的不是自己的
	if req.OnlyUserId == nil || req.GetOnlyUserId() != userData.GetUserId() {
		filter.OnlyStatus = lo.ToPtr(int64(core_api.PostStatus_PublicPostStatus))
	}

	// 查看的自己的
	if req.GetOnlyUserId() != "" && req.GetOnlyUserId() == userData.GetUserId() {
		filter.OnlyStatus = req.OnlyStatus
	}

	var search *content.SearchOption
	if req.SearchKeyword != nil {
		search = &content.SearchOption{
			SearchKeyword:  req.SearchKeyword,
			SearchSortType: content.SearchSortType(*req.SearchType),
			SearchTimeType: content.SearchTimeType(*req.SearchTimerType),
		}
	}
	if getPostsResp, err = s.CloudMindContent.GetPosts(ctx, &content.GetPostsReq{
		SearchOption:      search,
		PostFilterOptions: filter,
		PaginationOptions: &basic.PaginationOptions{
			Limit:     req.Limit,
			LastToken: req.LastToken,
			Backward:  req.Backward,
			Offset:    req.Offset,
		},
	}); err != nil {
		return resp, err
	}
	resp.Posts = make([]*core_api.Post, len(getPostsResp.Posts))
	if err = mr.Finish(lo.Map(getPostsResp.Posts, func(item *content.Post, i int) func() error {
		return func() error {
			resp.Posts[i] = &core_api.Post{
				PostId:   item.PostId,
				Title:    item.Title,
				Text:     item.Text,
				Url:      item.Url,
				Labels:   make([]*core_api.Label, 0),
				UserName: item.UserId,
			}
			if err = mr.Finish(func() error {
				s.PostDomainService.LoadLikeCount(ctx, &resp.Posts[i].LikeCount, item.PostId) // 点赞量
				return nil
			}, func() error {
				// 加载评论量
				return nil
			}, func() error {
				if userData.GetUserId() != "" {
					s.PostDomainService.LoadLiked(ctx, &resp.Posts[i].Liked, userData.GetUserId(), item.PostId)
				}
				return nil
			}, func() error {
				return nil
			}, func() error {
				s.PostDomainService.LoadLabels(ctx, &resp.Posts[i].Labels, item.LabelIds)
				return nil
			}, func() error {
				resp.Posts[i].ViewCount, _ = s.Redis.PfcountCtx(ctx, fmt.Sprintf("%s:%s", consts.ViewCountKey, item.PostId))
				return nil
			}); err != nil {
				return err
			}
			return nil
		}
	})...); err != nil {
		return resp, err
	}
	resp.Total = getPostsResp.Total
	resp.Token = getPostsResp.Token
	return resp, nil
}

func (s *PostService) CheckIsMyPost(ctx context.Context, postId, userId string) (*content.GetPostResp, error) {
	getPostResp, err := s.CloudMindContent.GetPost(ctx, &content.GetPostReq{
		PostId: postId,
	})
	if err != nil {
		return nil, err
	}
	if getPostResp.UserId != userId {
		return nil, consts.ErrForbidden
	}
	return getPostResp, nil
}

func (s *PostService) CheckIsMyPosts(ctx context.Context, postIds []string, userId string) (err error) {
	post, err := s.CloudMindContent.GetPostsByPostIds(ctx, &content.GetPostsByPostIdsReq{
		PostIds: postIds,
	})
	if err != nil {
		return err
	}
	for i := range post.Posts {
		if post.Posts[i].UserId != userId {
			return consts.ErrForbidden
		}
	}
	return nil
}
