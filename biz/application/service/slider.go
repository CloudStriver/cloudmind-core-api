package service

import (
	"context"
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/consts"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/convertor"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_system"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/platform_relation"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/basic"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/system"
	"github.com/google/wire"
	"github.com/samber/lo"
)

type ISliderService interface {
	CreateSlider(ctx context.Context, req *core_api.CreateSliderReq) (resp *core_api.CreateSliderResp, err error)
	UpdateSlider(ctx context.Context, req *core_api.UpdateSliderReq) (resp *core_api.UpdateSliderResp, err error)
	DeleteSlider(ctx context.Context, req *core_api.DeleteSliderReq) (resp *core_api.DeleteSliderResp, err error)
	GetSliders(ctx context.Context, c *core_api.GetSlidersReq) (*core_api.GetSlidersResp, error)
}

var SliderServiceSet = wire.NewSet(
	wire.Struct(new(SliderService), "*"),
	wire.Bind(new(ISliderService), new(*SliderService)),
)

type SliderService struct {
	Config           *config.Config
	CloudMindSystem  cloudmind_system.ICloudMindSystem
	PLatFromRelation platform_relation.IPlatFormRelation
}

func (s *SliderService) CreateSlider(ctx context.Context, req *core_api.CreateSliderReq) (resp *core_api.CreateSliderResp, err error) {
	if _, err := s.CloudMindSystem.CreateSlider(ctx, &system.CreateSliderReq{
		ImageUrl: req.ImageUrl,
		LinkUrl:  req.LinkUrl,
		IsPublic: req.IsPublic,
	}); err != nil {
		return resp, err
	}

	return resp, nil
}

func (s *SliderService) UpdateSlider(ctx context.Context, req *core_api.UpdateSliderReq) (resp *core_api.UpdateSliderResp, err error) {
	if _, err = s.CloudMindSystem.UpdateSlider(ctx, &system.UpdateSliderReq{
		SliderId: req.SliderId,
		ImageUrl: req.ImageUrl,
		LinkUrl:  req.LinkUrl,
		IsPublic: req.IsPublic,
	}); err != nil {
		return resp, err
	}

	return resp, nil
}

func (s *SliderService) DeleteSlider(ctx context.Context, req *core_api.DeleteSliderReq) (resp *core_api.DeleteSliderResp, err error) {
	if _, err = s.CloudMindSystem.DeleteSlider(ctx, &system.DeleteSliderReq{
		SliderId: req.SliderId,
	}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *SliderService) GetSliders(ctx context.Context, req *core_api.GetSlidersReq) (resp *core_api.GetSlidersResp, err error) {
	resp = new(core_api.GetSlidersResp)
	getSlidersResp, _ := s.CloudMindSystem.GetSliders(ctx, &system.GetSlidersReq{
		OnlyIsPublic:      lo.ToPtr(int64(consts.PublicSlider)),
		PaginationOptions: &basic.PaginationOptions{},
	})

	resp.Sliders = lo.Map[*system.Slider, *core_api.Slider](getSlidersResp.Sliders, func(item *system.Slider, _ int) *core_api.Slider {
		return convertor.SystemSliderToSlider(item)
	})

	resp.Total = getSlidersResp.Total
	return resp, nil
}
