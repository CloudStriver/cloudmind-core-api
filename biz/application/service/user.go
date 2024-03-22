package service

import (
	"context"
	"fmt"
	"github.com/CloudStriver/cloudmind-core-api/biz/adaptor"
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	domainservice "github.com/CloudStriver/cloudmind-core-api/biz/domain/service"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/consts"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/convertor"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/kq"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_content"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_sts"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_trade"
	"github.com/CloudStriver/cloudmind-mq/app/util/message"
	"github.com/CloudStriver/go-pkg/utils/pconvertor"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/trade"
	"github.com/bytedance/sonic"
	"github.com/google/wire"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/mr"
)

type IUserService interface {
	GetUser(ctx context.Context, req *core_api.GetUserReq) (resp *core_api.GetUserResp, err error)
	UpdateUser(ctx context.Context, req *core_api.UpdateUserReq) (resp *core_api.UpdateUserResp, err error)
	GetUserDetail(ctx context.Context, req *core_api.GetUserDetailReq) (resp *core_api.GetUserDetailResp, err error)
	SearchUser(ctx context.Context, req *core_api.SearchUserReq) (resp *core_api.SearchUserResp, err error)
}

var UserServiceSet = wire.NewSet(
	wire.Struct(new(UserService), "*"),
	wire.Bind(new(IUserService), new(*UserService)),
)

type UserService struct {
	Config            *config.Config
	CloudMindContent  cloudmind_content.ICloudMindContent
	CloudMindTrade    cloudmind_trade.ICloudMindTrade
	UserDomainService domainservice.IUserDomainService
	CloudMindSts      cloudmind_sts.ICloudMindSts
	UpdateItemKq      *kq.UpdateItemKq
}

func (s *UserService) GetUser(ctx context.Context, req *core_api.GetUserReq) (resp *core_api.GetUserResp, err error) {
	getUserResp, err := s.CloudMindContent.GetUser(ctx, &content.GetUserReq{UserId: req.UserId})
	if err != nil {
		return resp, err
	}
	s.UserDomainService.LoadLabel(ctx, getUserResp.Labels)

	return &core_api.GetUserResp{
		UserId: req.UserId,
		Name:   getUserResp.Name,
		Url:    getUserResp.Url,
		Tags:   getUserResp.Labels,
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
		Labels:      req.Labels,
		Url:         req.Url,
	}); err != nil {
		return resp, err
	}
	if len(req.Labels) != 0 {
		data, _ := sonic.Marshal(&message.UpdateItemMessage{
			ItemId: userData.UserId,
			Labels: req.Labels,
		})
		if err = s.UpdateItemKq.Push(pconvertor.Bytes2String(data)); err != nil {
			return resp, err
		}
	}
	return resp, nil
}

func (s *UserService) GetUserDetail(ctx context.Context, _ *core_api.GetUserDetailReq) (resp *core_api.GetUserDetailResp, err error) {
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	var (
		err1, err2     error
		getUserResp    *content.GetUserResp
		getBalanceResp *trade.GetBalanceResp
	)

	if err = mr.Finish(func() error {
		if getUserResp, err1 = s.CloudMindContent.GetUser(ctx, &content.GetUserReq{
			UserId: userData.UserId,
		}); err1 != nil {
			return err1
		}
		s.UserDomainService.LoadLabel(ctx, getUserResp.Labels)
		return nil
	}, func() error {
		if getBalanceResp, err2 = s.CloudMindTrade.GetBalance(ctx, &trade.GetBalanceReq{
			UserId: userData.UserId,
		}); err2 != nil {
			return err2
		}
		return nil
	}); err != nil {
		return resp, err
	}

	return &core_api.GetUserDetailResp{
		Name:        getUserResp.Name,
		Sex:         getUserResp.Sex,
		FullName:    getUserResp.FullName,
		IdCard:      getUserResp.IdCard,
		Description: getUserResp.Description,
		Url:         getUserResp.Url,
		Flow:        getBalanceResp.Flow,
		Momery:      getBalanceResp.Memory,
		Point:       getBalanceResp.Point,
		Labels:      getUserResp.Labels,
	}, nil
}

func (s *UserService) SearchUser(ctx context.Context, req *core_api.SearchUserReq) (resp *core_api.SearchUserResp, err error) {
	resp = new(core_api.SearchUserResp)
	var searchOptions *content.SearchOptions

	if req.AllFieldsKey != nil {
		searchOptions = &content.SearchOptions{
			Type: &content.SearchOptions_AllFieldsKey{
				AllFieldsKey: *req.AllFieldsKey,
			},
		}
	}
	if req.Name != nil || req.Description != nil {
		searchOptions = &content.SearchOptions{
			Type: &content.SearchOptions_MultiFieldsKey{
				MultiFieldsKey: &content.SearchField{
					Name:        req.Name,
					Description: req.Description,
				},
			},
		}
	}

	users, err := s.CloudMindContent.GetUsers(ctx, &content.GetUsersReq{
		SearchOptions:     searchOptions,
		PaginationOptions: convertor.MakePaginationOptions(req.Limit, req.Offset, req.LastToken, req.Backward),
	})
	if err != nil {
		fmt.Println(err)
		return resp, err
	}
	resp.Users = lo.Map[*content.User, *core_api.User](users.Users, func(user *content.User, _ int) *core_api.User {
		return convertor.UserDetailToUser(user)
	})
	resp.LastToken = users.LastToken
	return resp, nil
}
