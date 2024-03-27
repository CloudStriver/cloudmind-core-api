package service

import (
	"context"
	"fmt"
	"github.com/CloudStriver/cloudmind-core-api/biz/adaptor"
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/consts"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/convertor"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_content"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
	"github.com/google/wire"
	"github.com/samber/lo"
)

type IZoneService interface {
	GetZone(ctx context.Context, req *core_api.GetZoneReq) (resp *core_api.GetZoneResp, err error)
	GetZones(ctx context.Context, req *core_api.GetZonesReq) (resp *core_api.GetZonesResp, err error)
	CreateZone(ctx context.Context, req *core_api.CreateZoneReq) (resp *core_api.CreateZoneResp, err error)
	UpdateZone(ctx context.Context, req *core_api.UpdateZoneReq) (resp *core_api.UpdateZoneResp, err error)
	DeleteZone(ctx context.Context, req *core_api.DeleteZoneReq) (resp *core_api.DeleteZoneResp, err error)
}

type ZoneService struct {
	Config           *config.Config
	CloudMindContent cloudmind_content.ICloudMindContent
}

var ZoneServiceSet = wire.NewSet(
	wire.Struct(new(ZoneService), "*"),
	wire.Bind(new(IZoneService), new(*ZoneService)),
)

func (s *ZoneService) GetZone(ctx context.Context, req *core_api.GetZoneReq) (resp *core_api.GetZoneResp, err error) {
	resp = new(core_api.GetZoneResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	var res *content.GetZoneResp
	if res, err = s.CloudMindContent.GetZone(ctx, &content.GetZoneReq{Id: req.Id}); err != nil {
		return resp, err
	}
	resp.Zone = convertor.ZoneToCoreZone(res.Zone)
	return resp, nil
}

func (s *ZoneService) GetZones(ctx context.Context, req *core_api.GetZonesReq) (resp *core_api.GetZonesResp, err error) {
	resp = new(core_api.GetZonesResp)
	var res *content.GetZonesResp

	p := convertor.MakePaginationOptions(req.Limit, req.Offset, req.LastToken, req.Backward)
	if res, err = s.CloudMindContent.GetZones(ctx, &content.GetZonesReq{FatherId: req.FatherId, PaginationOptions: p}); err != nil {
		return resp, err
	}

	fmt.Printf("[%v\n]", res.Zones)
	resp.Zones = lo.Map(res.Zones, func(item *content.Zone, _ int) *core_api.Zone {
		return convertor.ZoneToCoreZone(item)
	})
	resp.Token = res.Token
	resp.Total = res.Total
	return resp, nil
}

func (s *ZoneService) CreateZone(ctx context.Context, req *core_api.CreateZoneReq) (resp *core_api.CreateZoneResp, err error) {
	resp = new(core_api.CreateZoneResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	var res *content.CreateZoneResp
	zone := &content.Zone{
		FatherId: req.FatherId,
		Value:    req.Value,
	}
	if res, err = s.CloudMindContent.CreateZone(ctx, &content.CreateZoneReq{Zone: zone}); err != nil {
		return resp, err
	}
	resp.Id = res.Id
	return resp, nil
}

func (s *ZoneService) UpdateZone(ctx context.Context, req *core_api.UpdateZoneReq) (resp *core_api.UpdateZoneResp, err error) {
	resp = new(core_api.UpdateZoneResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	zone := convertor.CoreZoneToZone(req.Zone)
	if _, err = s.CloudMindContent.UpdateZone(ctx, &content.UpdateZoneReq{Zone: zone}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *ZoneService) DeleteZone(ctx context.Context, req *core_api.DeleteZoneReq) (resp *core_api.DeleteZoneResp, err error) {
	resp = new(core_api.DeleteZoneResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	if _, err = s.CloudMindContent.DeleteZone(ctx, &content.DeleteZoneReq{Id: req.Id}); err != nil {
		return resp, err
	}
	return resp, nil
}
