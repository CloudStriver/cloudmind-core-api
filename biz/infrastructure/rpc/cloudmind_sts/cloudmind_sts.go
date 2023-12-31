package cloudmind_sts

import (
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/sts/stsservice"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/discovery"
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

func NewCloudMindSts(etcd discovery.Resolver) stsservice.Client {
	return stsservice.MustNewClient("cloudmind-sts", client.WithResolver(etcd))
}
