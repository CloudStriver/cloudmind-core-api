// Code generated by hertz generator.

package core_api

import (
	"context"
	"github.com/CloudStriver/cloudmind-core-api/biz/adaptor"
	"github.com/CloudStriver/cloudmind-core-api/provider"

	core_api "github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// CreateLabel .
// @router /label/createLabel [POST]
func CreateLabel(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.CreateLabelReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.CreateLabelResp)
	p := provider.Get()
	resp, err = p.LabelService.CreateLabel(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// DeleteLabel .
// @router /label/deleteLabel [POST]
func DeleteLabel(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.DeleteLabelReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.DeleteLabelResp)
	p := provider.Get()
	resp, err = p.LabelService.DeleteLabel(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// UpdateLabel .
// @router /label/updateLabel [POST]
func UpdateLabel(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.UpdateLabelReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.UpdateLabelResp)
	p := provider.Get()
	resp, err = p.LabelService.UpdateLabel(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// GetLabels .
// @router /label/getLabels [GET]
func GetLabels(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.GetLabelsReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.GetLabelsResp)
	p := provider.Get()
	resp, err = p.LabelService.GetLabels(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// GetLabel .
// @router /label/getLabel [GET]
func GetLabel(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.GetLabelReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.GetLabelResp)
	p := provider.Get()
	resp, err = p.LabelService.GetLabel(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// GetLabelsInBatch .
// @router /label/getLabelsInBatch [GET]
func GetLabelsInBatch(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.GetLabelsInBatchReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.GetLabelsInBatchResp)
	p := provider.Get()
	resp, err = p.LabelService.GetLabelsInBatch(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}
