// Code generated by hertz generator.

package core_api

import (
	"context"
	"github.com/CloudStriver/cloudmind-core-api/biz/adaptor"
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/cloudmind-core-api/provider"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// CreateRelation .
// @router /relation/create [POST]
func CreateRelation(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.CreateRelationReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.CreateRelationResp)
	p := provider.Get()
	resp, err = p.RelationService.CreateRelation(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// GetRelation .
// @router /relation [GET]
func GetRelation(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.GetRelationReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.GetRelationResp)
	p := provider.Get()
	resp, err = p.RelationService.GetRelation(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// GetFromRelations .
// @router /relations/from [GET]
func GetFromRelations(ctx context.Context, c *app.RequestContext) {
	// this my demo
	var err error
	var req core_api.GetFromRelationsReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.GetFromRelationsResp)
	p := provider.Get()
	resp, err = p.RelationService.GetFromRelations(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// GetToRelations .
// @router /relations/to [GET]
func GetToRelations(ctx context.Context, c *app.RequestContext) {
	// this my demo
	var err error
	var req core_api.GetToRelationsReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.GetToRelationsResp)
	p := provider.Get()
	resp, err = p.RelationService.GetToRelations(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// DeleteRelation .
// @router /relation/delete [POST]
func DeleteRelation(ctx context.Context, c *app.RequestContext) {
	// this my demo
	var err error
	var req core_api.DeleteRelationReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.DeleteRelationResp)
	p := provider.Get()
	resp, err = p.RelationService.DeleteRelation(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}
