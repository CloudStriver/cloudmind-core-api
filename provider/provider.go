package provider

import (
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/etcd"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_content"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_sts"
	"github.com/cloudwego/kitex/pkg/discovery"
	"github.com/google/wire"

	"github.com/CloudStriver/cloudmind-core-api/biz/application/service"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
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
	Config         *config.Config
	Etcd           discovery.Resolver
	ContentService service.IContentService
	AuthService    service.IAuthService
}

func Get() *Provider {
	return provider
}

var RPCSet = wire.NewSet(
	cloudmind_content.CloudMindContentSet,
	cloudmind_sts.CloudMindStsSet,
)

var ApplicationSet = wire.NewSet(
	service.ContentServiceSet,
	//service.CollectionServiceSet,
	service.AuthServiceSet,
	//service.CommentServiceSet,
	//service.UserServiceSet,
	//service.MomentServiceSet,
	//service.LikeServiceSet,
	//service.PostServiceSet,
	//service.SystemServiceSet,
	//service.StsServiceSet,
	//service.PlanServiceSet,
	//service.IncentiveServiceSet,
)

var DomainSet = wire.NewSet(
//domainservice.UserDomainServiceSet,
//domainservice.CommentDomainServiceSet,
//domainservice.MomentDomainServiceSet,
//domainservice.PostDomainServiceSet,
//domainservice.CatImageDomainServiceSet,
)

var InfrastructureSet = wire.NewSet(
	config.NewConfig,
	etcd.NewEtcd,
	RPCSet,
)

var AllProvider = wire.NewSet(
	ApplicationSet,
	DomainSet,
	InfrastructureSet,
)
