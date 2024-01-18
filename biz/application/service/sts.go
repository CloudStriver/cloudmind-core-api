package service

import (
	"context"
	"fmt"
	"github.com/CloudStriver/cloudmind-core-api/biz/adaptor"
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/consts"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_sts"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/sts"

	"github.com/google/wire"

	"net/http"
)

type IStsService interface {
	ApplySignedUrl(ctx context.Context, req *core_api.ApplySignedUrlReq) (*core_api.ApplySignedUrlResp, error)
}

type StsService struct {
	PlatformSts cloudmind_sts.ICloudMindSts
}

var StsServiceSet = wire.NewSet(
	wire.Struct(new(StsService), "*"),
	wire.Bind(new(IStsService), new(*StsService)),
)

func (s *StsService) ApplySignedUrl(ctx context.Context, req *core_api.ApplySignedUrlReq) (*core_api.ApplySignedUrlResp, error) {
	user := adaptor.ExtractUserMeta(ctx)
	if user.GetUserId() == "" {
		return nil, consts.ErrNotAuthentication
	}
	resp := new(core_api.ApplySignedUrlResp)
	userId := user.GetUserId()
	fmt.Println(userId)
	if req.File {
		// TODO：扣取用户流量
	}

	// TODO：判断是否已经上传过

	data, err := s.PlatformSts.GenCosSts(ctx, &sts.GenCosStsReq{Path: "users/*"})
	if err != nil {
		return nil, err
	}
	resp.SessionToken = data.SessionToken
	data2, err := s.PlatformSts.GenSignedUrl(ctx, &sts.GenSignedUrlReq{
		SecretId:  data.SecretId,
		SecretKey: data.SecretKey,
		Method:    http.MethodPut,
		Path:      "users/" + req.Md5 + req.Suffix,
	})
	resp.Url = data2.SignedUrl
	return resp, nil
}
