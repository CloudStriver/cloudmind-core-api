package cloudmind_trade

import (
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/go-pkg/utils/kitex/client"
	//"github.com/CloudStriver/go-pkg/utils/kitex/client"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/trade/tradeservice"
	"github.com/google/wire"
)

type ICloudMindTrade interface {
	tradeservice.Client
}

type CloudMindTrade struct {
	tradeservice.Client
}

var CloudMindTradeSet = wire.NewSet(
	NewCloudMindTrade,
	wire.Struct(new(CloudMindTrade), "*"),
	wire.Bind(new(ICloudMindTrade), new(*CloudMindTrade)),
)

func NewCloudMindTrade(config *config.Config) tradeservice.Client {
	return client.NewClient(config.Name, "cloudmind-trade", tradeservice.NewClient)
}
