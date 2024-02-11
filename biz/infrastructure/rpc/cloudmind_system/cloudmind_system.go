package cloudmind_system

import (
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/go-pkg/utils/kitex/client"
	//"github.com/CloudStriver/go-pkg/utils/kitex/client"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/system/systemservice"
	"github.com/google/wire"
)

type ICloudMindSystem interface {
	systemservice.Client
}

type CloudMindSystem struct {
	systemservice.Client
}

var CloudMindSystemSet = wire.NewSet(
	NewCloudMindSystem,
	wire.Struct(new(CloudMindSystem), "*"),
	wire.Bind(new(ICloudMindSystem), new(*CloudMindSystem)),
)

func NewCloudMindSystem(config *config.Config) systemservice.Client {
	return client.NewClient(config.Name, "cloudmind-system", systemservice.NewClient)
}
