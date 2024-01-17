package provider

import (
	"github.com/CloudStriver/cloudmind-core-api/biz/application/service"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_content"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_sts"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/platform_relation"
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
	Config          *config.Config
	ContentService  service.IContentService
	AuthService     service.IAuthService
	RelationService service.IRelationService
}

func Get() *Provider {
	return provider
}

var RPCSet = wire.NewSet(
	cloudmind_content.CloudMindContentSet,
	cloudmind_sts.CloudMindStsSet,
	platform_relation.PlatFormRelationSet,
)

var ApplicationSet = wire.NewSet(
	service.ContentServiceSet,
	service.RelationServiceSet,
	service.AuthServiceSet,
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
	RPCSet,
)

var AllProvider = wire.NewSet(
	ApplicationSet,
	DomainSet,
	InfrastructureSet,
)
