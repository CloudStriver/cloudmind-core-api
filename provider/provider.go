package provider

import (
	"github.com/CloudStriver/cloudmind-core-api/biz/application/service"
	domainservice "github.com/CloudStriver/cloudmind-core-api/biz/domain/service"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/kq"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_content"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_sts"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_system"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_trade"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/platform"
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
	NotificationService service.INotificationService
	CommentService      service.ICommentService
	LabelService        service.ILabelService
	RecommendService    service.IRecommendService
	ProductService      service.IProductService
	SliderService       service.ISliderService
	HotRankService      service.HotRankService
}

func Get() *Provider {
	return provider
}

var RPCSet = wire.NewSet(
	cloudmind_content.CloudMindContentSet,
	cloudmind_sts.CloudMindStsSet,
	cloudmind_system.CloudMindSystemSet,
	platform.PlatFormSet,
	cloudmind_trade.CloudMindTradeSet,
)

var ApplicationSet = wire.NewSet(
	service.FileServiceSet,
	service.RelationServiceSet,
	service.AuthServiceSet,
	service.PostServiceSet,
	service.UserServiceSet,
	service.NotificationServiceSet,
	service.RecommendServiceSet,
	service.LabelServiceSet,
	service.CommentServiceSet,
	service.ProductServiceSet,
	service.SliderServiceSet,
	service.HotRankServiceSet,
)

var DomainSet = wire.NewSet(
	domainservice.PostDomainServiceSet,
	domainservice.FileDomainServiceSet,
	domainservice.ProductDomainServiceSet,
	domainservice.UserDomainServiceSet,
	domainservice.CommentDomainServiceSet,
	domainservice.RelationDomainServiceSet,
)

var InfrastructureSet = wire.NewSet(
	config.NewConfig,
	redis.NewRedis,
	kq.NewCreateNotificationsKq,
	kq.NewCreateItemKq,
	kq.NewCreateFeedBackKq,
	kq.NewDeleteItemKq,
	kq.NewUpdateItemKq,
	kq.NewDeleteNotificationsKq,
	kq.NewDeleteFileRelationKq,
	RPCSet,
)

var AllProvider = wire.NewSet(
	ApplicationSet,
	DomainSet,
	InfrastructureSet,
)
