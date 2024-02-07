package service

import (
	"context"
	"github.com/CloudStriver/cloudmind-core-api/biz/adaptor"
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/consts"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/convertor"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_content"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_sts"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/sts"
	"github.com/google/wire"
)

type IUserService interface {
	GetUser(ctx context.Context, req *core_api.GetUserReq) (resp *core_api.GetUserResp, err error)
	UpdateUser(ctx context.Context, req *core_api.UpdateUserReq) (resp *core_api.UpdateUserResp, err error)
	GetUserDetail(ctx context.Context, req *core_api.GetUserDetailReq) (resp *core_api.GetUserDetailResp, err error)
	SearchUser(ctx context.Context, req *core_api.SearchUserReq) (resp *core_api.SearchUserResp, err error)
	AskUploadAvatar(ctx context.Context, req *core_api.AskUploadAvatarReq) (resp *core_api.AskUploadAvatarResp, err error)
}

var UserServiceSet = wire.NewSet(
	wire.Struct(new(UserService), "*"),
	wire.Bind(new(IUserService), new(*UserService)),
)

type UserService struct {
	Config           *config.Config
	CloudMindContent cloudmind_content.ICloudMindContent
	PlatformSts      cloudmind_sts.ICloudMindSts
}

func (s *UserService) AskUploadAvatar(ctx context.Context, req *core_api.AskUploadAvatarReq) (resp *core_api.AskUploadAvatarResp, err error) {
	resp = new(core_api.AskUploadAvatarResp)
	user := adaptor.ExtractUserMeta(ctx)
	if user.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	userId := user.GetUserId()
	if _, err = s.CloudMindContent.UpdateUser(ctx, &content.UpdateUserReq{
		UserId: userId,
		Url:    s.Config.GetUrl(req.Name),
	}); err != nil {
		return resp, err
	}

	genCosStsResp, err := s.PlatformSts.GenCosSts(ctx, &sts.GenCosStsReq{
		Path:   "users/" + req.Name,
		IsFile: false,
		Time:   req.AvatarSize / (1024 * 1024),
	})
	if err != nil {
		return resp, err
	}
	resp.SessionToken = genCosStsResp.SessionToken
	resp.TmpSecretId = genCosStsResp.SecretId
	resp.TmpSecretKey = genCosStsResp.SecretKey
	resp.StartTime = genCosStsResp.StartTime
	resp.ExpiredTime = genCosStsResp.ExpiredTime
	return resp, nil
}

func (s *UserService) GetUser(ctx context.Context, req *core_api.GetUserReq) (resp *core_api.GetUserResp, err error) {
	getUserResp, err := s.CloudMindContent.GetUser(ctx, &content.GetUserReq{UserId: req.UserId})
	if err != nil {
		return resp, err
	}
	return &core_api.GetUserResp{
		UserId: req.UserId,
		Name:   getUserResp.Name,
		Url:    getUserResp.Url,
	}, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *core_api.UpdateUserReq) (resp *core_api.UpdateUserResp, err error) {
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	if _, err = s.CloudMindContent.UpdateUser(ctx, &content.UpdateUserReq{
		UserId:      userData.UserId,
		Name:        req.Name,
		Sex:         req.Sex,
		FullName:    req.FullName,
		IdCard:      req.IdCard,
		Description: req.Description,
	}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *UserService) GetUserDetail(ctx context.Context, req *core_api.GetUserDetailReq) (resp *core_api.GetUserDetailResp, err error) {
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	getUserResp, err := s.CloudMindContent.GetUser(ctx, &content.GetUserReq{
		UserId: userData.GetUserId(),
	})
	if err != nil {
		return resp, err
	}

	return &core_api.GetUserDetailResp{
		Name:        getUserResp.Name,
		Sex:         getUserResp.Sex,
		FullName:    getUserResp.FullName,
		IdCard:      getUserResp.IdCard,
		Description: getUserResp.Description,
		Url:         getUserResp.Url,
	}, nil
}

func (s *UserService) SearchUser(ctx context.Context, req *core_api.SearchUserReq) (resp *core_api.SearchUserResp, err error) {
	resp = new(core_api.SearchUserResp)
	users, err := s.CloudMindContent.SearchUser(ctx, &content.SearchUserReq{
		Keyword:           req.Keyword,
		PaginationOptions: convertor.MakePaginationOptions(req.Limit, req.Offset, req.LastToken, req.Backward),
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
