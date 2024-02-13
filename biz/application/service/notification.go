package service

import (
	"context"
	"github.com/CloudStriver/cloudmind-core-api/biz/adaptor"
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/consts"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/convertor"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/kq"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_system"
	"github.com/CloudStriver/cloudmind-mq/app/util/message"
	"github.com/CloudStriver/go-pkg/utils/pconvertor"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/basic"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/system"
	"github.com/bytedance/sonic"
	"github.com/google/wire"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type INotificationService interface {
	GetNotifications(ctx context.Context, req *core_api.GetNotificationsReq) (resp *core_api.GetNotificationsResp, err error)
	GetNotificationCount(ctx context.Context, req *core_api.GetNotificationCountReq) (resp *core_api.GetNotificationCountResp, err error)
	UpdateNotifications(ctx context.Context, req *core_api.UpdateNotificationsReq) (resp *core_api.UpdateNotificationsResp, err error)
}

var NotificationServiceSet = wire.NewSet(
	wire.Struct(new(NotificationService), "*"),
	wire.Bind(new(INotificationService), new(*NotificationService)),
)

type NotificationService struct {
	Config                *config.Config
	CloudMindSystem       cloudmind_system.ICloudMindSystem
	UpdateNotificationsKq *kq.UpdateNotificationsKq
	Redis                 *redis.Redis
}

func (s *NotificationService) GetNotifications(ctx context.Context, req *core_api.GetNotificationsReq) (resp *core_api.GetNotificationsResp, err error) {
	resp = new(core_api.GetNotificationsResp)
	user := adaptor.ExtractUserMeta(ctx)
	if user.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	getNotificationsResp, err := s.CloudMindSystem.GetNotifications(ctx, &system.GetNotificationsReq{
		UserId:   user.UserId,
		OnlyType: req.OnlyType,
		PaginationOptions: &basic.PaginationOptions{
			Limit:     req.Limit,
			LastToken: req.LastToken,
			Backward:  req.Backward,
			Offset:    req.Offset,
		},
	})
	if err != nil {
		return resp, err
	}

	resp.Notifications = lo.Map[*system.Notification, *core_api.Notification](getNotificationsResp.Notifications, func(item *system.Notification, index int) *core_api.Notification {
		return convertor.NotificationToCoreNotification(item)
	})
	resp.Token = getNotificationsResp.Token
	return resp, nil
}

func (s *NotificationService) GetNotificationCount(ctx context.Context, req *core_api.GetNotificationCountReq) (resp *core_api.GetNotificationCountResp, err error) {
	resp = new(core_api.GetNotificationCountResp)
	user := adaptor.ExtractUserMeta(ctx)
	if user.GetUserId() != "" {
		getNotificationCountResp, err := s.CloudMindSystem.GetNotificationCount(ctx, &system.GetNotificationCountReq{
			UserId:   user.UserId,
			OnlyType: req.OnlyType,
		})
		if err != nil {
			return resp, err
		}
		resp.Total = getNotificationCountResp.Total
	}

	return resp, nil
}

func (s *NotificationService) UpdateNotifications(ctx context.Context, req *core_api.UpdateNotificationsReq) (resp *core_api.UpdateNotificationsResp, err error) {
	resp = new(core_api.UpdateNotificationsResp)
	user := adaptor.ExtractUserMeta(ctx)
	if user.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	data, _ := sonic.Marshal(&message.UpdateNotificationsMessage{
		UserId:   user.UserId,
		OnlyType: req.OnlyType,
	})
	if err = s.UpdateNotificationsKq.Push(pconvertor.Bytes2String(data)); err != nil {
		return resp, err
	}

	return resp, nil
}
