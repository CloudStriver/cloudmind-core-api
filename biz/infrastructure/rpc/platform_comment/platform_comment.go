package platform_comment

import (
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/go-pkg/utils/kitex/client"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/platform/comment/commentservice"
	"github.com/google/wire"
)

type IPlatFormComment interface {
	commentservice.Client
}

type PlatFormComment struct {
	commentservice.Client
}

var PlatFormCommentSet = wire.NewSet(
	NewPlatFormComment,
	wire.Struct(new(PlatFormComment), "*"),
	wire.Bind(new(IPlatFormComment), new(*PlatFormComment)),
)

func NewPlatFormComment(config *config.Config) commentservice.Client {
	return client.NewClient(config.Name, "platform-comment", commentservice.NewClient)
}
