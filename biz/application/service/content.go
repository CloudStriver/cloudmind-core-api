package service

import (
	"context"
	"github.com/CloudStriver/cloudmind-core-api/biz/adaptor"
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/consts"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/convertor"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_content"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
	"github.com/google/wire"
)

type IContentService interface {
	SearchUser(ctx context.Context, req *core_api.SearchUserReq) (resp *core_api.SearchUserResp, err error)
	UpdateUser(ctx context.Context, req *core_api.UpdateUserReq) (resp *core_api.UpdateUserResp, err error)
}

var ContentServiceSet = wire.NewSet(
	wire.Struct(new(ContentService), "*"),
	wire.Bind(new(IContentService), new(*ContentService)),
)

type ContentService struct {
	Config           *config.Config
	CloudMindContent cloudmind_content.ICloudMindContent
}

func (s *ContentService) SearchUser(ctx context.Context, req *core_api.SearchUserReq) (resp *core_api.SearchUserResp, err error) {
	resp = new(core_api.SearchUserResp)
	users, err := s.CloudMindContent.SearchUser(ctx, &content.SearchUserReq{
		Keyword:           req.Keyword,
		PaginationOptions: convertor.PaginationOptionsToPaginationOptions(req.PaginationOptions),
	})
	if err != nil {
		return resp, err
	}
	resp.Users = make([]*core_api.User, 0, len(users.Users))
	for _, user := range users.Users {
		resp.Users = append(resp.Users, convertor.UserDetailToUser(user))
	}
	resp.Total = users.Total
	resp.LastToken = users.LastToken
	return resp, nil
}

func (s *ContentService) UpdateUser(ctx context.Context, req *core_api.UpdateUserReq) (resp *core_api.UpdateUserResp, err error) {
	resp = new(core_api.UpdateUserResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	user := convertor.UserToUserDetailInfo(req.UserDetail)
	user.UserId = userData.UserId

	if _, err = s.CloudMindContent.UpdateUser(ctx, &content.UpdateUserReq{UserDetailInfo: user}); err != nil {
		return resp, err
	}
	return resp, nil
}
