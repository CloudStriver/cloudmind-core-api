package service

import (
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_content"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_sts"
	"github.com/google/wire"
)

type IStsService interface {
}

type StsService struct {
	PlatformSts      cloudmind_sts.ICloudMindSts
	CloudMindContent cloudmind_content.ICloudMindContent
}

var StsServiceSet = wire.NewSet(
	wire.Struct(new(StsService), "*"),
	wire.Bind(new(IStsService), new(*StsService)),
)
