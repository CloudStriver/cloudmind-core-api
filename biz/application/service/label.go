package service

import (
	"context"
	"github.com/CloudStriver/cloudmind-core-api/biz/adaptor"
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/consts"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/convertor"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_content"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
	"github.com/google/wire"
)

type ILabelService interface {
	GetLabel(ctx context.Context, req *core_api.GetLabelReq) (resp *core_api.GetLabelResp, err error)
	CreateLabel(ctx context.Context, req *core_api.CreateLabelReq) (resp *core_api.CreateLabelResp, err error)
	UpdateLabel(ctx context.Context, req *core_api.UpdateLabelReq) (resp *core_api.UpdateLabelResp, err error)
	DeleteLabel(ctx context.Context, req *core_api.DeleteLabelReq) (resp *core_api.DeleteLabelResp, err error)
}

var LabelServiceSet = wire.NewSet(
	wire.Struct(new(LabelService), "*"),
	wire.Bind(new(ILabelService), new(*LabelService)),
)

type LabelService struct {
	Config           *config.Config
	CloudMindContent cloudmind_content.ICloudMindContent
}

func (s *LabelService) GetLabel(ctx context.Context, req *core_api.GetLabelReq) (resp *core_api.GetLabelResp, err error) {
	resp = new(core_api.GetLabelResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	var res *content.GetLabelResp
	if res, err = s.CloudMindContent.GetLabel(ctx, &content.GetLabelReq{Id: req.Id}); err != nil {
		return resp, err
	}
	resp.Label = convertor.LabelToCoreLabel(res.Label)
	return resp, nil
}

func (s *LabelService) CreateLabel(ctx context.Context, req *core_api.CreateLabelReq) (resp *core_api.CreateLabelResp, err error) {
	resp = new(core_api.CreateLabelResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	var res *content.CreateLabelResp
	label := convertor.CoreLabelToLabel(req.Label)
	if res, err = s.CloudMindContent.CreateLabel(ctx, &content.CreateLabelReq{Label: label}); err != nil {
		return resp, err
	}
	resp.Id = res.Id
	return resp, nil
}

func (s *LabelService) UpdateLabel(ctx context.Context, req *core_api.UpdateLabelReq) (resp *core_api.UpdateLabelResp, err error) {
	resp = new(core_api.UpdateLabelResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	label := convertor.CoreLabelToLabel(req.Label)
	if _, err = s.CloudMindContent.UpdateLabel(ctx, &content.UpdateLabelReq{Label: label}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *LabelService) DeleteLabel(ctx context.Context, req *core_api.DeleteLabelReq) (resp *core_api.DeleteLabelResp, err error) {
	resp = new(core_api.DeleteLabelResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	if _, err = s.CloudMindContent.DeleteLabel(ctx, &content.DeleteLabelReq{Id: req.Id}); err != nil {
		return resp, err
	}
	return resp, nil
}
