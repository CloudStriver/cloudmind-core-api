package platform

import (
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/go-pkg/utils/kitex/client"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/platform/platformservice"
	"github.com/google/wire"
)

type IPlatForm interface {
	platformservice.Client
}

type PlatForm struct {
	platformservice.Client
}

var PlatFormSet = wire.NewSet(
	NewPlatFormComment,
	wire.Struct(new(PlatForm), "*"),
	wire.Bind(new(IPlatForm), new(*PlatForm)),
)

func NewPlatFormComment(config *config.Config) platformservice.Client {
	return client.NewClient(config.Name, "platform", platformservice.NewClient)
}
