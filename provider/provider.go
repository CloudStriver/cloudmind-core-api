package provider

import (
	"github.com/CloudStriver/cloudmind-core-api/biz/application/service"
	domainservice "github.com/CloudStriver/cloudmind-core-api/biz/domain/service"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/kq"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_content"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_sts"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/platform_comment"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_system"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/platform_relation"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/store/redis"
	"github.com/google/wire"
)

var provider *Provider

func Init() {
	var err error
	provider, err = NewProvider()
	if err != nil {
		panic(err)
	}
}

// Provider 提供controller依赖的对象
type Provider struct {
	Config              *config.Config
	FileService         service.IFileService
	PostService         service.IPostService
	AuthService         service.IAuthService
	RelationService     service.IRelationService
	UserService         service.IUserService
	ZoneService         service.IZoneService
	NotificationService service.INotificationService
}

func Get() *Provider {
	return provider
}

var RPCSet = wire.NewSet(
	cloudmind_content.CloudMindContentSet,
	cloudmind_sts.CloudMindStsSet,
	platform_relation.PlatFormRelationSet,
	cloudmind_system.CloudMindSystemSet,
	platform_comment.PlatFormCommentSet,
)

var ApplicationSet = wire.NewSet(
	service.FileServiceSet,
	service.RelationServiceSet,
	service.AuthServiceSet,
	service.PostServiceSet,
	service.UserServiceSet,
	service.ZoneServiceSet,
	service.NotificationServiceSet,
)

var DomainSet = wire.NewSet(
	domainservice.PostDomainServiceSet,
	domainservice.FileDomainServiceSet,
)

var InfrastructureSet = wire.NewSet(
	config.NewConfig,
	redis.NewRedis,
	kq.NewCreateNotificationsKq,
	kq.NewDeleteNotificationsKq,
	kq.NewUpdateNotificationsKq,
	RPCSet,
)

var AllProvider = wire.NewSet(
	ApplicationSet,
	DomainSet,
	InfrastructureSet,
)
