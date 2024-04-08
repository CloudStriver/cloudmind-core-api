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
	"github.com/CloudStriver/cloudmind-mq/app/util/message"
	"github.com/CloudStriver/go-pkg/utils/pconvertor"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
	"github.com/bytedance/sonic"
	"github.com/google/wire"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/mr"
)

type IRecommendService interface {
	GetLatestRecommend(ctx context.Context, req *core_api.GetLatestRecommendReq) (resp *core_api.GetLatestRecommendResp, err error)
	GetPopularRecommend(ctx context.Context, req *core_api.GetPopularRecommendReq) (resp *core_api.GetPopularRecommendResp, err error)
	CreateFeedBack(ctx context.Context, req *core_api.CreateFeedBackReq) (resp *core_api.CreateFeedBackResp, err error)
	GetRecommendByItem(ctx context.Context, req *core_api.GetRecommendByItemReq) (resp *core_api.GetRecommendByItemResp, err error)
	GetRecommendByUser(ctx context.Context, req *core_api.GetRecommendByUserReq) (resp *core_api.GetRecommendByUserResp, err error)
}

var RecommendServiceSet = wire.NewSet(
	wire.Struct(new(RecommendService), "*"),
	wire.Bind(new(IRecommendService), new(*RecommendService)),
)

type RecommendService struct {
	Config            *config.Config
	CloudMindContent  cloudmind_content.ICloudMindContent
	PostDomainService service.IPostDomainService
	CreateFeedBackKq  *kq.CreateFeedBackKq
	UserDomainService service.IUserDomainService
}

func (s *RecommendService) GetLatestRecommend(ctx context.Context, req *core_api.GetLatestRecommendReq) (resp *core_api.GetLatestRecommendResp, err error) {
	resp = new(core_api.GetLatestRecommendResp)
	resp.Recommends = new(core_api.Recommends)
	user, err := adaptor.ExtractUserMeta(ctx)
	if err != nil {
		return resp, consts.ErrNotAuthentication
	}
	userId := user.GetUserId()
	if userId == "" {
		userId = "default"
	}
	getLatestRecommendResp, err := s.CloudMindContent.GetLatestRecommend(ctx, &content.GetLatestRecommendReq{
		UserId:   userId,
		Limit:    req.Limit,
		Category: core_api.Category_name[int32(req.Category)],
	})
	if err != nil {
		return resp, err
	}
	if err = s.GetItemByItemId(ctx, userId, req.Category, getLatestRecommendResp.ItemIds, resp.Recommends); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *RecommendService) GetPopularRecommend(ctx context.Context, req *core_api.GetPopularRecommendReq) (resp *core_api.GetPopularRecommendResp, err error) {
	resp = new(core_api.GetPopularRecommendResp)
	resp.Recommends = new(core_api.Recommends)
	user, err := adaptor.ExtractUserMeta(ctx)
	if err != nil {
		return resp, consts.ErrNotAuthentication
	}
	userId := user.GetUserId()
	if userId == "" {
		userId = "default"
	}
	getPopularRecommendResp, err := s.CloudMindContent.GetPopularRecommend(ctx, &content.GetPopularRecommendReq{
		UserId:   userId,
		Limit:    req.Limit,
		Category: core_api.Category_name[int32(req.Category)],
	})
	if err != nil {
		return resp, err
	}
	if err = s.GetItemByItemId(ctx, userId, req.Category, getPopularRecommendResp.ItemIds, resp.Recommends); err != nil {
		return resp, err
	}

	return resp, nil
}

func (s *RecommendService) CreateFeedBack(ctx context.Context, req *core_api.CreateFeedBackReq) (resp *core_api.CreateFeedBackResp, err error) {
	user, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || user.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	data, _ := sonic.Marshal(&message.CreateFeedBackMessage{
		FeedbackType: req.FeedbackType,
		UserId:       user.GetUserId(),
		ItemId:       req.ItemId,
	})
	if err = s.CreateFeedBackKq.Push(pconvertor.Bytes2String(data)); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *RecommendService) GetRecommendByItem(ctx context.Context, req *core_api.GetRecommendByItemReq) (resp *core_api.GetRecommendByItemResp, err error) {
	resp = new(core_api.GetRecommendByItemResp)
	resp.Recommends = new(core_api.Recommends)
	user, err := adaptor.ExtractUserMeta(ctx)
	if err != nil {
		return resp, consts.ErrNotAuthentication
	}
	userId := user.GetUserId()
	if userId == "" {
		userId = "default"
	}
	getRecommendByItemResp, err := s.CloudMindContent.GetRecommendByItem(ctx, &content.GetRecommendByItemReq{
		ItemId:   req.ItemId,
		Limit:    req.Limit,
		Category: core_api.Category_name[int32(req.Category)],
	})
	if err != nil {
		return resp, err
	}

	if err = s.GetItemByItemId(ctx, userId, req.Category, getRecommendByItemResp.ItemIds, resp.Recommends); err != nil {
		return resp, err
	}

	return resp, nil
}

func (s *RecommendService) GetRecommendByUser(ctx context.Context, req *core_api.GetRecommendByUserReq) (resp *core_api.GetRecommendByUserResp, err error) {
	resp = new(core_api.GetRecommendByUserResp)
	resp.Recommends = new(core_api.Recommends)
	user, err := adaptor.ExtractUserMeta(ctx)
	if err != nil {
		return resp, consts.ErrNotAuthentication
	}
	userId := user.GetUserId()
	if userId == "" {
		userId = "default"
	}
	getRecommendByItemResp, err := s.CloudMindContent.GetRecommendByUser(ctx, &content.GetRecommendByUserReq{
		UserId:   userId,
		Limit:    req.Limit,
		Category: core_api.Category_name[int32(req.Category)],
	})
	if err != nil {
		return resp, err
	}

	if err = s.GetItemByItemId(ctx, userId, req.Category, getRecommendByItemResp.ItemIds, resp.Recommends); err != nil {
		return resp, err
	}

	return resp, nil
}

func (s *RecommendService) GetItemByItemId(ctx context.Context, userId string, category core_api.Category, itemIds []string, recommends *core_api.Recommends) (err error) {
	switch category {
	case core_api.Category_UserCategory:
		getUsersResp, err := s.CloudMindContent.GetUsersByUserIds(ctx, &content.GetUsersByUserIdsReq{
			UserIds: itemIds,
		})
		if err != nil {
			return err
		}
		recommends.Users = make([]*core_api.User, len(getUsersResp.Users))
		if err = mr.Finish(lo.Map(getUsersResp.Users, func(user *content.User, i int) func() error {
			return func() error {
				recommends.Users[i] = &core_api.User{
					UserId:      user.UserId,
					Name:        user.Name,
					Url:         user.Url,
					Labels:      user.Labels,
					Description: user.Description,
				}
				_ = mr.Finish(func() error {
					s.UserDomainService.LoadFollowedCount(ctx, &recommends.Users[i].FollowedCount, user.UserId)
					return nil
				}, func() error {
					s.UserDomainService.LoadLabel(ctx, recommends.Users[i].Labels)
					return nil
				}, func() error {
					s.UserDomainService.LoadFollowed(ctx, &recommends.Users[i].Followed, userId, user.UserId)
					return nil
				})
				return nil
			}
		})...); err != nil {
			return err
		}
	case core_api.Category_FileCategory:
	case core_api.Category_PostCategory:
		getPostsResp, err := s.CloudMindContent.GetPostsByPostIds(ctx, &content.GetPostsByPostIdsReq{
			PostIds: itemIds,
		})
		if err != nil {
			return err
		}
		getPostsResp.Posts = lo.Filter[*content.Post](getPostsResp.Posts, func(item *content.Post, index int) bool {
			return item.Status == int64(core_api.PostStatus_PublicPostStatus)
		})
		recommends.Posts = make([]*core_api.Post, len(getPostsResp.Posts))
		if err = mr.Finish(lo.Map(getPostsResp.Posts, func(post *content.Post, i int) func() error {
			return func() error {
				recommends.Posts[i] = &core_api.Post{
					PostId:   post.PostId,
					Title:    post.Title,
					Text:     post.Text,
					Url:      post.Url,
					LabelIds: post.LabelIds,
				}
				if err = mr.Finish(func() error {
					s.PostDomainService.LoadLikeCount(ctx, &recommends.Posts[i].LikeCount, post.PostId) // 点赞量
					return nil
				}, func() error {
					s.PostDomainService.LoadCommentCount(ctx, &recommends.Posts[i].CommentCount, post.PostId) // 评论量
					return nil
				}, func() error {
					s.PostDomainService.LoadLiked(ctx, &recommends.Posts[i].Liked, userId, post.PostId)
					return nil
				}, func() error {
					getUser, _ := s.CloudMindContent.GetUser(ctx, &content.GetUserReq{UserId: post.UserId})
					recommends.Posts[i].UserName = getUser.Name
					return nil
				}, func() error {
					s.PostDomainService.LoadLabels(ctx, post.LabelIds)
					return nil
				}); err != nil {
					return err
				}
				return nil
			}
		})...); err != nil {
			return err
		}

	default:
	}
	return nil
}
