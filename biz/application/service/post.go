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
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
	"github.com/google/wire"
	"github.com/samber/lo"
)

type IPostService interface {
	CreatePost(ctx context.Context, req *core_api.CreatePostReq) (resp *core_api.CreatePostResp, err error)
	UpdatePost(ctx context.Context, req *core_api.UpdatePostReq) (resp *core_api.UpdatePostResp, err error)
	DeletePost(ctx context.Context, req *core_api.DeletePostReq) (resp *core_api.DeletePostResp, err error)
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
