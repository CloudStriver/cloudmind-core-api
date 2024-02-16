// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package provider

import (
	service2 "github.com/CloudStriver/cloudmind-core-api/biz/application/service"
	"github.com/CloudStriver/cloudmind-core-api/biz/domain/service"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/kq"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_content"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_sts"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_system"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_trade"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/platform_comment"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/platform_relation"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/store/redis"
)

// Injectors from wire.go:

func NewProvider() (*Provider, error) {
	configConfig, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	client := cloudmind_sts.NewCloudMindSts(configConfig)
	cloudMindSts := &cloudmind_sts.CloudMindSts{
		Client: client,
	}
	contentserviceClient := cloudmind_content.NewCloudMindContent(configConfig)
	cloudMindContent := &cloudmind_content.CloudMindContent{
		Client: contentserviceClient,
	}
	relationserviceClient := platform_relation.NewPlatFormRelation(configConfig)
	platFormRelation := &platform_relation.PlatFormRelation{
		Client: relationserviceClient,
	}
	commentserviceClient := platform_comment.NewPlatFormComment(configConfig)
	platFormComment := &platform_comment.PlatFormComment{
		Client: commentserviceClient,
	}
	fileDomainService := &service.FileDomainService{
		CloudMindUser:    cloudMindContent,
		PlatformRelation: platFormRelation,
		PlatformComment:  platFormComment,
	}
	fileService := &service2.FileService{
		Config:            configConfig,
		PlatformSts:       cloudMindSts,
		CloudMindContent:  cloudMindContent,
		FileDomainService: fileDomainService,
		PlatformComment:   platFormComment,
	}
	postDomainService := &service.PostDomainService{
		CloudMindUser:    cloudMindContent,
		PlatformRelation: platFormRelation,
	}
	createItemsKq := kq.NewCreateItemsKq(configConfig)
	updateItemKq := kq.NewUpdateItemKq(configConfig)
	deleteItemKq := kq.NewDeleteItemKq(configConfig)
	postService := &service2.PostService{
		Config:            configConfig,
		CloudMindContent:  cloudMindContent,
		PostDomainService: postDomainService,
		PLatFromRelation:  platFormRelation,
		CreateItemsKq:     createItemsKq,
		UpdateItemKq:      updateItemKq,
		DeleteItemKq:      deleteItemKq,
	}
	tradeserviceClient := cloudmind_trade.NewCloudMindTrade(configConfig)
	cloudMindTrade := &cloudmind_trade.CloudMindTrade{
		Client: tradeserviceClient,
	}
	redisRedis := redis.NewRedis(configConfig)
	authService := &service2.AuthService{
		Config:           configConfig,
		CloudMindContent: cloudMindContent,
		CloudMindSts:     cloudMindSts,
		CloudMindTrade:   cloudMindTrade,
		CreateItemsKq:    createItemsKq,
		Redis:            redisRedis,
	}
	createNotificationsKq := kq.NewCreateNotificationsKq(configConfig)
	createFeedBacksKq := kq.NewCreateFeedBacksKq(configConfig)
	relationService := &service2.RelationService{
		Config:               configConfig,
		PlatFormRelation:     platFormRelation,
		CloudMindContent:     cloudMindContent,
		PostDomainService:    postDomainService,
		CreateNotificationKq: createNotificationsKq,
		CreateFeedBacksKq:    createFeedBacksKq,
	}
	userService := &service2.UserService{
		Config:           configConfig,
		CloudMindContent: cloudMindContent,
		CloudMindTrade:   cloudMindTrade,
		PlatformSts:      cloudMindSts,
	}
	zoneService := &service2.ZoneService{
		Config:           configConfig,
		CloudMindContent: cloudMindContent,
	}
	systemserviceClient := cloudmind_system.NewCloudMindSystem(configConfig)
	cloudMindSystem := &cloudmind_system.CloudMindSystem{
		Client: systemserviceClient,
	}
	updateNotificationsKq := kq.NewUpdateNotificationsKq(configConfig)
	notificationService := &service2.NotificationService{
		Config:                configConfig,
		CloudMindSystem:       cloudMindSystem,
		UpdateNotificationsKq: updateNotificationsKq,
		Redis:                 redisRedis,
	}
	commentService := &service2.CommentService{
		Config:          configConfig,
		PlatformComment: platFormComment,
	}
	labelService := &service2.LabelService{
		Config:          configConfig,
		PlatformComment: platFormComment,
	}
	userDomainService := &service.UserDomainService{
		Config:           configConfig,
		PlatFormRelation: platFormRelation,
	}
	recommendService := &service2.RecommendService{
		Config:            configConfig,
		CloudMindContent:  cloudMindContent,
		PostDomainService: postDomainService,
		CreateFeedBacks:   createFeedBacksKq,
		UserDomainService: userDomainService,
	}
	productDomainService := &service.ProductDomainService{
		CloudMindUser:    cloudMindContent,
		PlatformRelation: platFormRelation,
		CloudMindTrade:   cloudMindTrade,
	}
	productService := &service2.ProductService{
		Config:               configConfig,
		CloudMindContent:     cloudMindContent,
		ProductDomainService: productDomainService,
		CloudMindTrade:       cloudMindTrade,
		CreateItemsKq:        createItemsKq,
		UpdateItemKq:         updateItemKq,
		DeleteItemKq:         deleteItemKq,
	}
	providerProvider := &Provider{
		Config:              configConfig,
		FileService:         fileService,
		PostService:         postService,
		AuthService:         authService,
		RelationService:     relationService,
		UserService:         userService,
		ZoneService:         zoneService,
		NotificationService: notificationService,
		CommentService:      commentService,
		LabelService:        labelService,
		RecommendService:    recommendService,
		ProductService:      productService,
	}
	return providerProvider, nil
}
