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
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/basic"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/system"
	"github.com/google/wire"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type INotificationService interface {
	GetNotifications(ctx context.Context, req *core_api.GetNotificationsReq) (resp *core_api.GetNotificationsResp, err error)
	GetNotificationCount(ctx context.Context, req *core_api.GetNotificationCountReq) (resp *core_api.GetNotificationCountResp, err error)
	DeleteNotifications(ctx context.Context, req *core_api.DeleteNotificationsReq) (resp *core_api.DeleteNotificationsResp, err error)
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
	DeleteNotificationsKq *kq.DeleteNotificationsKq
	Redis                 *redis.Redis
}

func (s *NotificationService) GetNotifications(ctx context.Context, req *core_api.GetNotificationsReq) (resp *core_api.GetNotificationsResp, err error) {
	resp = new(core_api.GetNotificationsResp)
	user := adaptor.ExtractUserMeta(ctx)
	if user.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	getNotificationsResp, err := s.CloudMindSystem.GetNotifications(ctx, &system.GetNotificationsReq{
		OnlyUserId: lo.ToPtr(user.UserId),
		OnlyType:   req.OnlyType,
		OnlyIsRead: req.OnlyIsRead,
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
	resp.Total = getNotificationsResp.Total
	resp.Token = getNotificationsResp.Token
	return resp, nil
}

func (s *NotificationService) GetNotificationCount(ctx context.Context, req *core_api.GetNotificationCountReq) (resp *core_api.GetNotificationCountResp, err error) {
	resp = new(core_api.GetNotificationCountResp)
	user := adaptor.ExtractUserMeta(ctx)
	if user.GetUserId() != "" {
		getNotificationCountResp, err := s.CloudMindSystem.GetNotificationCount(ctx, &system.GetNotificationCountReq{
			OnlyUserId: lo.ToPtr(user.UserId),
			OnlyType:   req.OnlyType,
			OnlyIsRead: req.OnlyIsRead,
		})
		if err != nil {
			return resp, err
		}
		resp.Total = getNotificationCountResp.Total
	}

	return resp, nil
}

func (s *NotificationService) DeleteNotifications(ctx context.Context, req *core_api.DeleteNotificationsReq) (resp *core_api.DeleteNotificationsResp, err error) {
	user := adaptor.ExtractUserMeta(ctx)
	if user.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	if err = s.DeleteNotificationsKq.Add(user.UserId, &message.DeleteNotificationsMessage{
		OnlyNotificationIds: req.OnlyNotificationIds,
		OnlyUserId:          lo.ToPtr(user.UserId),
		OnlyType:            req.OnlyType,
		OnlyIsRead:          req.OnlyIsRead,
	}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *NotificationService) UpdateNotifications(ctx context.Context, req *core_api.UpdateNotificationsReq) (resp *core_api.UpdateNotificationsResp, err error) {
	user := adaptor.ExtractUserMeta(ctx)
	if user.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	if err = s.UpdateNotificationsKq.Add(user.UserId, &message.UpdateNotificationsMessage{
		OnlyNotificationIds: req.OnlyNotificationIds,
		OnlyUserId:          lo.ToPtr(user.UserId),
		OnlyType:            req.OnlyType,
		OnlyIsRead:          lo.ToPtr(consts.NotificationNotRead),
		IsRead:              consts.NotificationRead,
	}); err != nil {
		return resp, err
	}
	return resp, nil
}
