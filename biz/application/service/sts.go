package service

import (
	"context"
	"fmt"
	"github.com/CloudStriver/cloudmind-core-api/biz/adaptor"
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/consts"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_content"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_sts"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/sts"

	"github.com/google/wire"
)

type IStsService interface {
	ApplySignedUrl(ctx context.Context, req *core_api.ApplySignedUrlReq) (*core_api.ApplySignedUrlResp, error)
}

type StsService struct {
	PlatformSts      cloudmind_sts.ICloudMindSts
	CloudMindContent cloudmind_content.ICloudMindContent
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
	//if req.File {
	//	 TODO：扣取用户流量
	//}
	// TODO：判断是否已经上传过
	//getFileIsExist, err := s.CloudMindContent.GetFileIsExist(ctx, &content.GetFileIsExistReq{
	//	Md5: req.Md5,
	//})
	//if err != nil {
	//	return resp, err
	//}
	//
	//if getFileIsExist.Ok {
	//	resp.Url = "exist"
	//	return resp, nil
	//}

	genCosStsResp, err := s.PlatformSts.GenCosSts(ctx, &sts.GenCosStsReq{Path: "users/*"})
	if err != nil {
		return nil, err
	}
	resp.SessionToken = genCosStsResp.SessionToken
	resp.TmpSecretId = genCosStsResp.SecretId
	resp.TmpSecretKey = genCosStsResp.SecretKey
	resp.StartTime = genCosStsResp.StartTime
	resp.ExpiredTime = genCosStsResp.ExpiredTime
	return resp, nil
}
