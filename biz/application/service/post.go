package service

import (
	"context"
	"github.com/CloudStriver/cloudmind-core-api/biz/adaptor"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/kq"
	"github.com/CloudStriver/cloudmind-mq/app/util/message"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/basic"

	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/cloudmind-core-api/biz/domain/service"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/consts"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_content"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
	"github.com/google/wire"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/mr"
)

type IPostService interface {
	CreatePost(ctx context.Context, req *core_api.CreatePostReq) (resp *core_api.CreatePostResp, err error)
	UpdatePost(ctx context.Context, req *core_api.UpdatePostReq) (resp *core_api.UpdatePostResp, err error)
	DeletePost(ctx context.Context, req *core_api.DeletePostReq) (resp *core_api.DeletePostResp, err error)
	GetOwnPost(ctx context.Context, req *core_api.GetOwnPostReq) (resp *core_api.GetOwnPostResp, err error)
	GetOwnPosts(ctx context.Context, req *core_api.GetOwnPostsReq) (resp *core_api.GetOwnPostsResp, err error)
	GetOtherPost(ctx context.Context, req *core_api.GetOtherPostReq) (resp *core_api.GetOtherPostResp, err error)
	GetOtherPosts(ctx context.Context, req *core_api.GetOtherPostsReq) (resp *core_api.GetOtherPostsResp, err error)
}

var PostServiceSet = wire.NewSet(
	wire.Struct(new(PostService), "*"),
	wire.Bind(new(IPostService), new(*PostService)),
)

type PostService struct {
	Config            *config.Config
	CloudMindContent  cloudmind_content.ICloudMindContent
	PostDomainService service.IPostDomainService
	CreateItemsKq     *kq.CreateItemsKq
	UpdateItemKq      *kq.UpdateItemKq
	DeleteItemKq      *kq.DeleteItemKq
}

func (s *PostService) CreatePost(ctx context.Context, req *core_api.CreatePostReq) (resp *core_api.CreatePostResp, err error) {
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	createPostResp, err := s.CloudMindContent.CreatePost(ctx, &content.CreatePostReq{
		UserId: userData.UserId,
		Title:  req.Title,
		Text:   req.Text,
		Tags:   req.Tags,
		Status: req.Status,
		Url:    req.Url,
	})
	if err != nil {
		return resp, err
	}

	if err = s.CreateItemsKq.Add(userData.UserId, &message.CreateItemsMessage{
		Item: &content.Item{
			ItemId:   createPostResp.PostId,
			IsHidden: req.Status == consts.PostPrivateStatus,
			Labels:   req.Tags,
			Category: core_api.Category_name[int32(core_api.Category_PostCategory)],
			Comment:  req.Title,
		},
	}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *PostService) UpdatePost(ctx context.Context, req *core_api.UpdatePostReq) (resp *core_api.UpdatePostResp, err error) {
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	if err = s.CheckIsMyPost(ctx, req.PostId, userData.UserId); err != nil {
		return resp, err
	}

	if _, err = s.CloudMindContent.UpdatePost(ctx, &content.UpdatePostReq{
		PostId: req.PostId,
		Title:  req.Title,
		Text:   req.Text,
		Tags:   req.Tags,
		Status: req.Status,
		Url:    req.Url,
	}); err != nil {
		return resp, err
	}

	if req.Status != 0 || req.Tags != nil || req.Title != "" {
		var isHidden *bool
		if req.Status != 0 {
			isHidden = lo.ToPtr(req.Status == consts.PostPrivateStatus)
		}
		if err = s.UpdateItemKq.Add(userData.UserId, &message.UpdateItemMessage{
			ItemId:   req.PostId,
			IsHidden: isHidden,
			Labels:   req.Tags,
			Comment:  lo.ToPtr(req.Title),
		}); err != nil {
			return resp, err
		}
	}
	return resp, nil
}

func (s *PostService) DeletePost(ctx context.Context, req *core_api.DeletePostReq) (resp *core_api.DeletePostResp, err error) {
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	// 只能删除自己的帖子
	if err = s.CheckIsMyPost(ctx, req.PostId, userData.UserId); err != nil {
		return resp, err
	}

	if _, err = s.CloudMindContent.DeletePost(ctx, &content.DeletePostReq{
		PostId: req.PostId,
	}); err != nil {
		return resp, err
	}

	if err = s.DeleteItemKq.Add(userData.UserId, &message.DeleteItemMessage{
		ItemId: req.PostId,
	}); err != nil {
		return resp, err
	}

	return resp, nil
}

func (s *PostService) GetOwnPost(ctx context.Context, req *core_api.GetOwnPostReq) (resp *core_api.GetOwnPostResp, err error) {
	userData := adaptor.ExtractUserMeta(ctx)
	var res *content.GetPostResp
	if res, err = s.CloudMindContent.GetPost(ctx, &content.GetPostReq{
		PostId: req.PostId,
	}); err != nil {
		return resp, err
	}

	if err = s.CheckIsMyPost(ctx, req.PostId, userData.UserId); err != nil {
		return resp, err
	}

	resp = &core_api.GetOwnPostResp{
		Title:      res.Title,
		Text:       res.Text,
		Url:        res.Url,
		Tags:       res.Tags,
		CreateTime: res.CreateTime,
		UpdateTime: res.UpdateTime,
		Author:     &core_api.User{},
	}
	if err = mr.Finish(func() error {
		s.PostDomainService.LoadAuthor(ctx, resp.Author, res.UserId) // 作者
		return nil
	}, func() error {
		s.PostDomainService.LoadLikeCount(ctx, &resp.LikeCount, req.PostId) // 点赞量
		return nil
	}, func() error {
		s.PostDomainService.LoadViewCount(ctx, &resp.ViewCount, req.PostId) // 浏览量
		return nil
	}, func() error {
		s.PostDomainService.LoadCollectCount(ctx, &resp.CollectCount, req.PostId) // 收藏量
		return nil
	}, func() error {
		s.PostDomainService.LoadLiked(ctx, &resp.Liked, userData.GetUserId(), req.PostId) // 是否点赞
		return nil
	}, func() error {
		s.PostDomainService.LoadCollected(ctx, &resp.Collected, userData.GetUserId(), req.PostId) // 是否收藏
		return nil
	}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *PostService) GetOwnPosts(ctx context.Context, req *core_api.GetOwnPostsReq) (resp *core_api.GetOwnPostsResp, err error) {
	resp = new(core_api.GetOwnPostsResp)
	var (
		getPostsResp  *content.GetPostsResp
		searchOptions *content.SearchOptions
	)

	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	if req.AllFieldsKey != nil {
		searchOptions = &content.SearchOptions{
			Type: &content.SearchOptions_AllFieldsKey{
				AllFieldsKey: *req.AllFieldsKey,
			},
		}
	}
	if req.Text != nil || req.Title != nil || req.Tag != nil {
		searchOptions = &content.SearchOptions{
			Type: &content.SearchOptions_MultiFieldsKey{
				MultiFieldsKey: &content.SearchField{
					Tag:   req.Tag,
					Text:  req.Text,
					Title: req.Title,
				},
			},
		}
	}
	if getPostsResp, err = s.CloudMindContent.GetPosts(ctx, &content.GetPostsReq{
		SearchOptions: searchOptions,
		PostFilterOptions: &content.PostFilterOptions{
			OnlyUserId:      lo.ToPtr(userData.UserId),
			OnlyTags:        req.OnlyTags,
			OnlySetRelation: req.OnlySetRelation,
		},
		PaginationOptions: &basic.PaginationOptions{
			Limit:     req.Limit,
			LastToken: req.LastToken,
			Backward:  req.Backward,
			Offset:    req.Offset,
		},
	}); err != nil {
		return resp, err
	}

	resp.Posts = make([]*core_api.OwnPost, len(getPostsResp.Posts))
	if err = mr.Finish(lo.Map(getPostsResp.Posts, func(item *content.Post, i int) func() error {
		return func() error {
			resp.Posts[i] = &core_api.OwnPost{
				PostId: item.PostId,
				Title:  item.Title,
				Text:   item.Text,
				Url:    item.Url,
				Tags:   item.Tags,
			}
			author := &core_api.User{}
			if err = mr.Finish(func() error {
				s.PostDomainService.LoadLikeCount(ctx, &resp.Posts[i].LikeCount, item.PostId) // 点赞量
				return nil
			}, func() error {
				// 加载评论量
				return nil
			}, func() error {
				s.PostDomainService.LoadLiked(ctx, &resp.Posts[i].Liked, userData.UserId, item.PostId)
				return nil
			}, func() error {
				s.PostDomainService.LoadAuthor(ctx, author, item.UserId)
				resp.Posts[i].UserName = author.Name
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

func (s *PostService) GetOtherPost(ctx context.Context, req *core_api.GetOtherPostReq) (resp *core_api.GetOtherPostResp, err error) {
	userData := adaptor.ExtractUserMeta(ctx)
	var res *content.GetPostResp
	if res, err = s.CloudMindContent.GetPost(ctx, &content.GetPostReq{
		PostId: req.PostId,
	}); err != nil {
		return resp, err
	}

	if res.Status != consts.PostPublicStatus {
		return resp, consts.ErrForbidden
	}

	resp = &core_api.GetOtherPostResp{
		Title:      res.Title,
		Text:       res.Text,
		Url:        res.Url,
		Tags:       res.Tags,
		CreateTime: res.CreateTime,
		UpdateTime: res.UpdateTime,
		Author:     &core_api.User{},
	}
	if err = mr.Finish(func() error {
		s.PostDomainService.LoadAuthor(ctx, resp.Author, res.UserId) // 作者
		return nil
	}, func() error {
		s.PostDomainService.LoadLikeCount(ctx, &resp.LikeCount, req.PostId) // 点赞量
		return nil
	}, func() error {
		s.PostDomainService.LoadViewCount(ctx, &resp.ViewCount, req.PostId) // 浏览量
		return nil
	}, func() error {
		s.PostDomainService.LoadCollectCount(ctx, &resp.CollectCount, req.PostId) // 收藏量
		return nil
	}, func() error {
		s.PostDomainService.LoadLiked(ctx, &resp.Liked, userData.GetUserId(), req.PostId) // 是否点赞
		return nil
	}, func() error {
		s.PostDomainService.LoadCollected(ctx, &resp.Collected, userData.GetUserId(), req.PostId) // 是否收藏
		return nil
	}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *PostService) GetOtherPosts(ctx context.Context, req *core_api.GetOtherPostsReq) (resp *core_api.GetOtherPostsResp, err error) {
	resp = new(core_api.GetOtherPostsResp)
	userData := adaptor.ExtractUserMeta(ctx)
	var (
		getPostsResp  *content.GetPostsResp
		searchOptions *content.SearchOptions
	)

	if req.AllFieldsKey != nil {
		searchOptions = &content.SearchOptions{
			Type: &content.SearchOptions_AllFieldsKey{
				AllFieldsKey: *req.AllFieldsKey,
			},
		}
	}
	if req.Text != nil || req.Title != nil || req.Tag != nil {
		searchOptions = &content.SearchOptions{
			Type: &content.SearchOptions_MultiFieldsKey{
				MultiFieldsKey: &content.SearchField{
					Tag:   req.Tag,
					Text:  req.Text,
					Title: req.Title,
				},
			},
		}
	}
	if getPostsResp, err = s.CloudMindContent.GetPosts(ctx, &content.GetPostsReq{
		SearchOptions: searchOptions,
		PostFilterOptions: &content.PostFilterOptions{
			OnlyUserId:      req.OnlyUserId,
			OnlyTags:        req.OnlyTags,
			OnlySetRelation: req.OnlySetRelation,
			OnlyStatus:      lo.ToPtr(consts.PostPublicStatus),
		},
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
				PostId: item.PostId,
				Title:  item.Title,
				Text:   item.Text,
				Url:    item.Url,
				Tags:   item.Tags,
			}
			author := &core_api.User{}
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
				s.PostDomainService.LoadAuthor(ctx, author, item.UserId)
				resp.Posts[i].UserName = author.Name
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

func (s *PostService) CheckIsMyPost(ctx context.Context, postId, userId string) (err error) {
	post, err := s.CloudMindContent.GetPost(ctx, &content.GetPostReq{
		PostId: postId,
	})
	if err != nil {
		return err
	}
	if post.UserId != userId {
		return consts.ErrForbidden
	}
	return nil
}
