package service

import (
	"context"
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
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/sts"
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

func (s *UserService) FiltetContet(ctx context.Context, IsSure bool, contents []*string) ([]string, error) {
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
		for i, content := range replaceContentResp.Content {
			*contents[i] = content
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

func (s *UserService) GetUser(ctx context.Context, req *core_api.GetUserReq) (resp *core_api.GetUserResp, err error) {
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil {
		return resp, consts.ErrNotAuthentication
	}
	getUserResp, err := s.CloudMindContent.GetUser(ctx, &content.GetUserReq{UserId: req.UserId})
	if err != nil {
		return resp, err
	}
	resp = &core_api.GetUserResp{
		UserId:        req.UserId,
		Name:          getUserResp.Name,
		Url:           getUserResp.Url,
		Description:   getUserResp.Description,
		Sex:           getUserResp.Sex,
		CreateTime:    getUserResp.CreateTime,
		BackgroundUrl: getUserResp.BackgroundUrl,
	}
	if err = mr.Finish(func() error {
		s.UserDomainService.LoadLabel(ctx, getUserResp.Labels)
		resp.Labels = getUserResp.Labels
		return nil
	}, func() error {
		if userData.GetUserId() != req.UserId {
			s.UserDomainService.LoadFollowed(ctx, &resp.Followed, userData.GetUserId(), req.UserId)
		}
		return nil
	}, func() error {
		s.UserDomainService.LoadFollowCount(ctx, &resp.FollowCount, req.UserId)
		return nil
	}, func() error {
		s.UserDomainService.LoadFollowedCount(ctx, &resp.FollowedCount, req.UserId)
		return nil
	}); err != nil {
		return resp, err
	}

	return resp, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *core_api.UpdateUserReq) (resp *core_api.UpdateUserResp, err error) {
	resp = new(core_api.UpdateUserResp)
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	if req.Name != "" || req.FullName != "" || req.Description != "" {
		resp.Keywords, err = s.FiltetContet(ctx, req.IsSure, []*string{&req.Name, &req.FullName, &req.Description})
		if err != nil {
			return resp, err
		}
		if resp.Keywords != nil {
			return resp, nil
		}
	}

	if _, err = s.CloudMindContent.UpdateUser(ctx, &content.UpdateUserReq{
		UserId:        userData.UserId,
		Name:          req.Name,
		Sex:           req.Sex,
		FullName:      req.FullName,
		IdCard:        req.IdCard,
		Description:   req.Description,
		Labels:        req.LabelIds,
		Url:           req.Url,
		BackgroundUrl: req.BackgroundUrl,
	}); err != nil {
		return resp, err
	}
	if len(req.LabelIds) != 0 {
		data, _ := sonic.Marshal(&message.UpdateItemMessage{
			ItemId: userData.UserId,
			Labels: req.LabelIds,
		})
		if err = s.UpdateItemKq.Push(pconvertor.Bytes2String(data)); err != nil {
			return resp, err
		}
	}
	return resp, nil
}

func (s *UserService) GetUserDetail(ctx context.Context, _ *core_api.GetUserDetailReq) (resp *core_api.GetUserDetailResp, err error) {
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || userData.GetUserId() == "" {
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
		Memory:      getBalanceResp.Memory,
		Point:       getBalanceResp.Point,
		Labels:      getUserResp.Labels,
		CreateTime:  getUserResp.CreateTime,
	}, nil
}

func (s *UserService) SearchUser(ctx context.Context, req *core_api.SearchUserReq) (resp *core_api.SearchUserResp, err error) {
	resp = new(core_api.SearchUserResp)
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil {
		return resp, consts.ErrNotAuthentication
	}

	users, err := s.CloudMindContent.GetUsers(ctx, &content.GetUsersReq{
		SearchOption: &content.SearchOption{
			SearchKeyword:  req.Keyword,
			SearchSortType: content.SearchSortType(req.SearchType),
			SearchTimeType: content.SearchTimeType(req.SearchTimerType),
		},
		PaginationOptions: convertor.MakePaginationOptions(req.Limit, req.Offset, req.LastToken, req.Backward),
	})
	if err != nil {
		return resp, err
	}
	resp.Users = make([]*core_api.User, len(users.Users))

	if err = mr.Finish(lo.Map[*content.User](users.Users, func(user *content.User, i int) func() error {
		return func() error {
			resp.Users[i] = &core_api.User{
				UserId:      user.UserId,
				Name:        user.Name,
				Url:         user.Url,
				Labels:      user.Labels,
				Description: user.Description,
			}
			if err = mr.Finish(func() error {
				if userData.GetUserId() != user.UserId && userData.GetUserId() != "" {
					s.UserDomainService.LoadFollowed(ctx, &resp.Users[i].Followed, userData.UserId, user.UserId)
				}
				return nil
			}, func() error {
				s.UserDomainService.LoadFollowedCount(ctx, &resp.Users[i].FollowedCount, user.UserId)
				return nil
			}); err != nil {
				return err
			}
			return nil
		}
	})...); err != nil {
		return resp, err
	}
	resp.LastToken = users.LastToken
	return resp, nil
}
