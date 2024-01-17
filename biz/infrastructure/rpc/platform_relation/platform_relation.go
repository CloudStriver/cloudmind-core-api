package platform_relation

import (
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/go-pkg/utils/kitex/client"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/platform/relation/relationservice"
	"github.com/google/wire"
)

type IPlatFormRelation interface {
	relationservice.Client
}

type PlatFormRelation struct {
	relationservice.Client
}

var PlatFormRelationSet = wire.NewSet(
	NewPlatFormRelation,
	wire.Struct(new(PlatFormRelation), "*"),
	wire.Bind(new(IPlatFormRelation), new(*PlatFormRelation)),
)

func NewPlatFormRelation(config *config.Config) relationservice.Client {
	return client.NewClient(config.Name, "platform-relation", relationservice.NewClient)
}
