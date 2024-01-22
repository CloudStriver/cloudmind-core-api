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

type IUserService interface {
	GetUser(ctx context.Context, req *core_api.GetUserReq) (resp *core_api.GetUserResp, err error)
	CreateUser(ctx context.Context, req *core_api.CreateUserReq) (resp *core_api.CreateUserResp, err error)
	UpdateUser(ctx context.Context, req *core_api.UpdateUserReq) (resp *core_api.UpdateUserResp, err error)
	GetUserDetail(ctx context.Context, req *core_api.GetUserDetailReq) (resp *core_api.GetUserDetailResp, err error)
	SearchUser(ctx context.Context, req *core_api.SearchUserReq) (resp *core_api.SearchUserResp, err error)
}

var UserServiceSet = wire.NewSet(
	wire.Struct(new(UserService), "*"),
	wire.Bind(new(IUserService), new(*UserService)),
)

type UserService struct {
	Config           *config.Config
	CloudMindContent cloudmind_content.ICloudMindContent
}

func (s *UserService) GetUser(ctx context.Context, req *core_api.GetUserReq) (resp *core_api.GetUserResp, err error) {
	resp = new(core_api.GetUserResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	if _, err = s.CloudMindContent.GetUser(ctx, &content.GetUserReq{UserId: req.UserId}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *UserService) CreateUser(ctx context.Context, req *core_api.CreateUserReq) (resp *core_api.CreateUserResp, err error) {
	resp = new(core_api.CreateUserResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	if _, err = s.CloudMindContent.CreateUser(ctx, &content.CreateUserReq{
		User: convertor.CoreUserInfoToUser(req.UserInfo),
	}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *core_api.UpdateUserReq) (resp *core_api.UpdateUserResp, err error) {
	resp = new(core_api.UpdateUserResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	user := convertor.UserToUserDetailInfo(req.UserDetail)
	user.UserId = userData.UserId

	if _, err = s.CloudMindContent.UpdateUser(ctx, &content.UpdateUserReq{
		User: convertor.CoreUserDetailToUser(req.UserDetail),
	}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *UserService) GetUserDetail(ctx context.Context, req *core_api.GetUserDetailReq) (resp *core_api.GetUserDetailResp, err error) {
	resp = new(core_api.GetUserDetailResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	if _, err = s.CloudMindContent.GetUser(ctx, &content.GetUserReq{
		UserId: req.UserId,
	}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *UserService) SearchUser(ctx context.Context, req *core_api.SearchUserReq) (resp *core_api.SearchUserResp, err error) {
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
