package service

import (
	"context"
	"github.com/CloudStriver/cloudmind-core-api/biz/adaptor"
	//"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/basic"
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/cloudmind-core-api/biz/domain/service"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/consts"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/convertor"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_content"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
	"github.com/google/wire"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/mr"
)

type IPostService interface {
	CreatePost(ctx context.Context, req *core_api.CreatePostReq) (resp *core_api.CreatePostResp, err error)
	GetPosts(ctx context.Context, req *core_api.GetPostsReq) (resp *core_api.GetPostsResp, err error)
	UpdatePost(ctx context.Context, req *core_api.UpdatePostReq) (resp *core_api.UpdatePostResp, err error)
	DeletePost(ctx context.Context, req *core_api.DeletePostReq) (resp *core_api.DeletePostResp, err error)
	GetPost(ctx context.Context, c *core_api.GetPostReq) (resp *core_api.GetPostResp, err error)
}

var PostServiceSet = wire.NewSet(
	wire.Struct(new(PostService), "*"),
	wire.Bind(new(IPostService), new(*PostService)),
)

type PostService struct {
	Config            *config.Config
	CloudMindContent  cloudmind_content.ICloudMindContent
	PostDomainService service.IPostDomainService
}

func (s *PostService) GetPost(ctx context.Context, c *core_api.GetPostReq) (resp *core_api.GetPostResp, err error) {
	resp = new(core_api.GetPostResp)
	userData := adaptor.ExtractUserMeta(ctx)

	var res *content.GetPostResp
	if res, err = s.CloudMindContent.GetPost(ctx, &content.GetPostReq{
		PostFilterOptions: &content.PostFilterOptions{
			OnlyPostId: lo.ToPtr(c.PostId),
		},
	}); err != nil {
		return resp, err
	}
	resp.Post = convertor.PostToCorePost(res.Post)
	if err = mr.Finish(func() error {
		s.PostDomainService.LoadLikeCount(ctx, resp.Post) // 点赞量
		return nil
	}, func() error {
		s.PostDomainService.LoadAuthor(ctx, resp.Post, res.Post.UserId) // 作者
		return nil
	}, func() error {
		s.PostDomainService.LoadViewCount(ctx, resp.Post) // 浏览量
		return nil
	}, func() error {
		s.PostDomainService.LoadLiked(ctx, resp.Post, userData.UserId) // 是否点赞
		return nil
	}, func() error {
		s.PostDomainService.LoadCollected(ctx, resp.Post, userData.GetUserId()) // 是否收藏
		return nil
	}, func() error {
		s.PostDomainService.LoadCollectCount(ctx, resp.Post) // 收藏量
		return nil
	}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *PostService) CreatePost(ctx context.Context, req *core_api.CreatePostReq) (resp *core_api.CreatePostResp, err error) {
	resp = new(core_api.CreatePostResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	if userData.UserId != req.PostInfo.UserId {
		return resp, consts.ErrNotPermission
	}

	if _, err = s.CloudMindContent.CreatePost(ctx, &content.CreatePostReq{
		Post: convertor.CorePostInfoToPostInfo(req.PostInfo),
	}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *PostService) GetPosts(ctx context.Context, req *core_api.GetPostsReq) (resp *core_api.GetPostsResp, err error) {
	resp = new(core_api.GetPostsResp)
	var getPostsResp *content.GetPostsResp
	if getPostsResp, err = s.CloudMindContent.GetPosts(ctx, &content.GetPostsReq{
		SearchOptions:     convertor.SearchOptionsToFileSearchOptions(req.GetSearchOptions()),
		PostFilterOptions: convertor.PostFilterOptionsToPostFilterOptions(req.GetPostFilterOptions()),
		PaginationOptions: convertor.PaginationOptionsToPaginationOptions(req.GetPaginationOptions()),
	}); err != nil {
		return resp, err
	}

	resp.Posts = make([]*core_api.Post, len(getPostsResp.Posts))
	if err = mr.Finish(lo.Map(getPostsResp.Posts, func(item *content.Post, i int) func() error {
		return func() error {
			p := convertor.PostToCorePost(item)
			if err = mr.Finish(func() error {
				s.PostDomainService.LoadLikeCount(ctx, p) // 点赞量
				return nil
			}, func() error {
				s.PostDomainService.LoadAuthor(ctx, p, item.UserId) // 作者
				return nil
			}, func() error {
				s.PostDomainService.LoadViewCount(ctx, p) // 浏览量
				return nil
			}); err != nil {
				return err
			}
			resp.Posts[i] = p
			return nil
		}
	})...); err != nil {
		return resp, err
	}
	resp.Total = getPostsResp.Total
	resp.Token = getPostsResp.Token
	return resp, nil
}

func (s *PostService) UpdatePost(ctx context.Context, req *core_api.UpdatePostReq) (resp *core_api.UpdatePostResp, err error) {
	resp = new(core_api.UpdatePostResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	if userData.UserId != req.PostInfo.UserId {
		return resp, consts.ErrNotPermission
	}

	if _, err = s.CloudMindContent.UpdatePost(ctx, &content.UpdatePostReq{
		Post: convertor.CorePostInfoToPostInfo(req.PostInfo),
	}); err != nil {
		return resp, err
	}

	return resp, nil
}

func (s *PostService) DeletePost(ctx context.Context, req *core_api.DeletePostReq) (resp *core_api.DeletePostResp, err error) {
	resp = new(core_api.DeletePostResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	if userData.UserId != req.PostId {
		return resp, consts.ErrNotPermission
	}

	if _, err = s.CloudMindContent.DeletePost(ctx, &content.DeletePostReq{
		PostId: req.PostId,
	}); err != nil {
		return resp, err
	}
	return resp, nil
}
