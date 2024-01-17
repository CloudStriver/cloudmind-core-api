package cloudmind_content

import (
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/go-pkg/utils/kitex/client"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content/contentservice"
	"github.com/google/wire"
)

type ICloudMindContent interface {
	contentservice.Client
}

type CloudMindContent struct {
	contentservice.Client
}

var CloudMindContentSet = wire.NewSet(
	NewCloudMindContent,
	wire.Struct(new(CloudMindContent), "*"),
	wire.Bind(new(ICloudMindContent), new(*CloudMindContent)),
)

func NewCloudMindContent(config *config.Config) contentservice.Client {
	return client.NewClient(config.Name, "cloudmind-content", contentservice.NewClient)
}
