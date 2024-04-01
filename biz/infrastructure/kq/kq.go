package kq

import (
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/zeromicro/go-queue/kq"
)

type CreateNotificationsKq struct {
	*kq.Pusher
}

type CreateItemKq struct {
	*kq.Pusher
}

type CreateFeedBackKq struct {
	*kq.Pusher
}

type UpdateItemKq struct {
	*kq.Pusher
}

type DeleteItemKq struct {
	*kq.Pusher
}
type DeleteNotificationsKq struct {
	*kq.Pusher
}

type DeleteFileRelationKq struct {
	*kq.Pusher
}

func NewCreateNotificationsKq(c *config.Config) *CreateNotificationsKq {
	pusher := kq.NewPusher(c.CreateNotificationsKq.Brokers, c.CreateNotificationsKq.Topic)
	return &CreateNotificationsKq{
		Pusher: pusher,
	}
}

func NewDeleteNotificationsKq(c *config.Config) *DeleteNotificationsKq {
	pusher := kq.NewPusher(c.DeleteNotificationsKq.Brokers, c.DeleteNotificationsKq.Topic)
	return &DeleteNotificationsKq{
		Pusher: pusher,
	}
}

func NewCreateItemKq(c *config.Config) *CreateItemKq {
	pusher := kq.NewPusher(c.CreateItemKq.Brokers, c.CreateItemKq.Topic)
	return &CreateItemKq{
		Pusher: pusher,
	}
}

func NewCreateFeedBackKq(c *config.Config) *CreateFeedBackKq {
	pusher := kq.NewPusher(c.CreateFeedBackKq.Brokers, c.CreateFeedBackKq.Topic)
	return &CreateFeedBackKq{
		Pusher: pusher,
	}
}

func NewUpdateItemKq(c *config.Config) *UpdateItemKq {
	pusher := kq.NewPusher(c.UpdateItemKq.Brokers, c.UpdateItemKq.Topic)
	return &UpdateItemKq{
		Pusher: pusher,
	}
}

func NewDeleteItemKq(c *config.Config) *DeleteItemKq {
	pusher := kq.NewPusher(c.DeleteItemKq.Brokers, c.DeleteItemKq.Topic)
	return &DeleteItemKq{
		Pusher: pusher,
	}
}

func NewDeleteFileRelationKq(c *config.Config) *DeleteFileRelationKq {
	pusher := kq.NewPusher(c.DeleteFileRelationKq.Brokers, c.DeleteFileRelationKq.Topic)
	return &DeleteFileRelationKq{
		Pusher: pusher,
	}
}
