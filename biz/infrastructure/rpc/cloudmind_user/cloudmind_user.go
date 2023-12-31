package cloudmind_user

import (
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/user/userservice"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/discovery"
	"github.com/google/wire"
)

type ICloudMindUser interface {
	userservice.Client
}

type CloudMindUser struct {
	userservice.Client
}

var CloudMindUserSet = wire.NewSet(
	NewCloudMindUser,
	wire.Struct(new(CloudMindUser), "*"),
	wire.Bind(new(ICloudMindUser), new(*CloudMindUser)),
)

func NewCloudMindUser(etcd discovery.Resolver) userservice.Client {
	return userservice.MustNewClient("cloudmind-user", client.WithResolver(etcd))
}
