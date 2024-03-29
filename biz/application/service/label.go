package service

import (
	"context"
	"github.com/CloudStriver/cloudmind-core-api/biz/adaptor"
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/consts"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/convertor"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/platform_comment"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/platform/comment"
	"github.com/google/wire"
	"github.com/samber/lo"
)

type ILabelService interface {
	CreateLabel(ctx context.Context, req *core_api.CreateLabelReq) (resp *core_api.CreateLabelResp, err error)
	DeleteLabel(ctx context.Context, req *core_api.DeleteLabelReq) (resp *core_api.DeleteLabelResp, err error)
	GetLabel(ctx context.Context, req *core_api.GetLabelReq) (resp *core_api.GetLabelResp, err error)
	GetLabelsInBatch(ctx context.Context, req *core_api.GetLabelsInBatchReq) (resp *core_api.GetLabelsInBatchResp, err error)
	UpdateLabel(ctx context.Context, req *core_api.UpdateLabelReq) (resp *core_api.UpdateLabelResp, err error)
	GetLabels(ctx context.Context, req *core_api.GetLabelsReq) (resp *core_api.GetLabelsResp, err error)
}

var LabelServiceSet = wire.NewSet(
	wire.Struct(new(LabelService), "*"),
	wire.Bind(new(ILabelService), new(*LabelService)),
)

type LabelService struct {
	Config          *config.Config
	PlatformComment platform_comment.IPlatFormComment
}

func (s *LabelService) CreateLabel(ctx context.Context, req *core_api.CreateLabelReq) (resp *core_api.CreateLabelResp, err error) {
	resp = new(core_api.CreateLabelResp)
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	var res *comment.CreateLabelResp
	if res, err = s.PlatformComment.CreateLabel(ctx, &comment.CreateLabelReq{Label: convertor.CoreLabelToLabel(req.Label)}); err != nil {
		return resp, err
	}
	resp.LabelId = res.Id
	return resp, nil
}

func (s *LabelService) DeleteLabel(ctx context.Context, req *core_api.DeleteLabelReq) (resp *core_api.DeleteLabelResp, err error) {
	resp = new(core_api.DeleteLabelResp)
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	if _, err = s.PlatformComment.DeleteLabel(ctx, &comment.DeleteLabelReq{Id: req.LabelId}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *LabelService) GetLabel(ctx context.Context, req *core_api.GetLabelReq) (resp *core_api.GetLabelResp, err error) {
	resp = new(core_api.GetLabelResp)

	var res *comment.GetLabelResp
	if res, err = s.PlatformComment.GetLabel(ctx, &comment.GetLabelReq{Id: req.LabelId}); err != nil {
		return resp, err
	}
	resp.Label = res.Label
	return resp, nil
}

func (s *LabelService) GetLabelsInBatch(ctx context.Context, req *core_api.GetLabelsInBatchReq) (resp *core_api.GetLabelsInBatchResp, err error) {
	resp = new(core_api.GetLabelsInBatchResp)

	var res *comment.GetLabelsInBatchResp
	if res, err = s.PlatformComment.GetLabelsInBatch(ctx, &comment.GetLabelsInBatchReq{LabelIds: req.LabelIds}); err != nil {
		return resp, err
	}
	resp.Labels = res.Labels
	return resp, nil
}

func (s *LabelService) UpdateLabel(ctx context.Context, req *core_api.UpdateLabelReq) (resp *core_api.UpdateLabelResp, err error) {
	resp = new(core_api.UpdateLabelResp)
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	if _, err = s.PlatformComment.UpdateLabel(ctx, &comment.UpdateLabelReq{Label: convertor.CoreLabelToLabel(req.Label)}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *LabelService) GetLabels(ctx context.Context, req *core_api.GetLabelsReq) (resp *core_api.GetLabelsResp, err error) {
	resp = new(core_api.GetLabelsResp)
	var res *comment.GetLabelsResp
	p := convertor.MakePaginationOptions(req.Limit, req.Offset, req.LastToken, req.Backward)
	if res, err = s.PlatformComment.GetLabels(ctx, &comment.GetLabelsReq{FilterOptions: &comment.LabelFilterOptions{Key: req.Key, Zone: req.Zone, SubZone: req.SubZone}, Pagination: p}); err != nil {
		return resp, err
	}
	resp.Labels = lo.Map(res.Labels, func(item *comment.Label, _ int) *core_api.Label {
		return convertor.LabelToCoreLabel(item)
	})
	resp.Total = res.Total
	resp.Token = res.Token
	return resp, nil
}
