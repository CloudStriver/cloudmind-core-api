package service

import (
	"context"
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/platform_comment"
	"github.com/google/wire"
)

type ILabelService interface {
	CreateLabel(ctx context.Context, req *core_api.CreateLabelReq) (resp *core_api.CreateLabelResp, err error)
	DeleteLabel(ctx context.Context, req *core_api.DeleteLabelReq) (resp *core_api.DeleteLabelResp, err error)
	GetLabel(ctx context.Context, req *core_api.GetLabelReq) (resp *core_api.GetLabelResp, err error)
	GetLabelsInBatch(ctx context.Context, req *core_api.GetLabelsInBatchReq) (resp *core_api.GetLabelsInBatchResp, err error)
	UpdateLabel(ctx context.Context, req *core_api.UpdateLabelReq) (resp *core_api.UpdateLabelResp, err error)
	GetLabels(ctx context.Context, req *core_api.GetLabelsReq) (resp *core_api.GetLabelsResp, err error)
	CreateObject(ctx context.Context, req *core_api.CreateObjectReq) (resp *core_api.CreateObjectResp, err error)
	CreateObjects(ctx context.Context, req *core_api.CreateObjectsReq) (resp *core_api.CreateObjectsResp, err error)
	DeleteObject(ctx context.Context, req *core_api.DeleteObjectReq) (resp *core_api.DeleteObjectResp, err error)
	GetObjects(ctx context.Context, req *core_api.GetObjectsReq) (resp *core_api.GetObjectsResp, err error)
	UpdateObject(ctx context.Context, req *core_api.UpdateObjectReq) (resp *core_api.UpdateObjectResp, err error)
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
	return
}

func (s *LabelService) DeleteLabel(ctx context.Context, req *core_api.DeleteLabelReq) (resp *core_api.DeleteLabelResp, err error) {
	return
}

func (s *LabelService) GetLabel(ctx context.Context, req *core_api.GetLabelReq) (resp *core_api.GetLabelResp, err error) {
	return
}

func (s *LabelService) GetLabelsInBatch(ctx context.Context, req *core_api.GetLabelsInBatchReq) (resp *core_api.GetLabelsInBatchResp, err error) {
	return
}

func (s *LabelService) UpdateLabel(ctx context.Context, req *core_api.UpdateLabelReq) (resp *core_api.UpdateLabelResp, err error) {
	return
}

func (s *LabelService) GetLabels(ctx context.Context, req *core_api.GetLabelsReq) (resp *core_api.GetLabelsResp, err error) {
	return
}

func (s *LabelService) CreateObject(ctx context.Context, req *core_api.CreateObjectReq) (resp *core_api.CreateObjectResp, err error) {
	return
}

func (s *LabelService) CreateObjects(ctx context.Context, req *core_api.CreateObjectsReq) (resp *core_api.CreateObjectsResp, err error) {
	return
}

func (s *LabelService) DeleteObject(ctx context.Context, req *core_api.DeleteObjectReq) (resp *core_api.DeleteObjectResp, err error) {
	return
}

func (s *LabelService) GetObjects(ctx context.Context, req *core_api.GetObjectsReq) (resp *core_api.GetObjectsResp, err error) {
	return
}

func (s *LabelService) UpdateObject(ctx context.Context, req *core_api.UpdateObjectReq) (resp *core_api.UpdateObjectResp, err error) {
	return
}
