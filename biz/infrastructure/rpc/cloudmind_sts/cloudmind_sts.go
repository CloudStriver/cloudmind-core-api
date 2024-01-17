package cloudmind_sts

import (
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/go-pkg/utils/kitex/client"
	//"github.com/CloudStriver/go-pkg/utils/kitex/client"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/sts/stsservice"
	"github.com/google/wire"
)

type ICloudMindSts interface {
	stsservice.Client
}

type CloudMindSts struct {
	stsservice.Client
}

var CloudMindStsSet = wire.NewSet(
	NewCloudMindSts,
	wire.Struct(new(CloudMindSts), "*"),
	wire.Bind(new(ICloudMindSts), new(*CloudMindSts)),
)

func NewCloudMindSts(config *config.Config) stsservice.Client {
	return client.NewClient(config.Name, "cloudmind-sts", stsservice.NewClient)
}
