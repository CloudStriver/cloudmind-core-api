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

// CreateComment .
// @router /comment/createComment [POST]
func CreateComment(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.CreateCommentReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.CreateCommentResp)
	p := provider.Get()
	resp, err = p.CommentService.CreateComment(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// GetComment .
// @router /comment/getComment [GET]
func GetComment(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.GetCommentReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.GetCommentResp)
	p := provider.Get()
	resp, err = p.CommentService.GetComment(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// GetComments .
// @router /comment/getComments [GET]
func GetComments(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.GetCommentsReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.GetCommentsResp)
	p := provider.Get()
	resp, err = p.CommentService.GetComments(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// DeleteComment .
// @router /comment/deleteComment [POST]
func DeleteComment(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.DeleteCommentReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.DeleteCommentResp)
	p := provider.Get()
	resp, err = p.CommentService.DeleteComment(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// UpdateComment .
// @router /comment/updateComment [POST]
func UpdateComment(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.UpdateCommentReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.UpdateCommentResp)
	p := provider.Get()
	resp, err = p.CommentService.UpdateComment(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// SetCommentAttrs .
// @router /comment/setCommentAttrs [POST]
func SetCommentAttrs(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.SetCommentAttrsReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.SetCommentAttrsResp)
	p := provider.Get()
	resp, err = p.CommentService.SetCommentAttrs(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// GetCommentSubject .
// @router /comment/getCommentSubject [GET]
func GetCommentSubject(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.GetCommentSubjectReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.GetCommentSubjectResp)
	p := provider.Get()
	resp, err = p.CommentService.GetCommentSubject(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// UpdateCommentSubject .
// @router /comment/updateCommentSubject [POST]
func UpdateCommentSubject(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.UpdateCommentSubjectReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.UpdateCommentSubjectResp)
	p := provider.Get()
	resp, err = p.CommentService.UpdateCommentSubject(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// DeleteCommentSubject .
// @router /comment/deleteCommentSubject [POST]
func DeleteCommentSubject(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.DeleteCommentSubjectReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.DeleteCommentSubjectResp)
	p := provider.Get()
	resp, err = p.CommentService.DeleteCommentSubject(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// CreateRelation .
// @router /relation/createRelation [POST]
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

// GetFromRelations .
// @router /relation/getFromRelations [GET]
func GetFromRelations(ctx context.Context, c *app.RequestContext) {
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
// @router /relation/getToRelations [GET]
func GetToRelations(ctx context.Context, c *app.RequestContext) {
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

// GetRelation .
// @router /relation/getRelation [GET]
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

// DeleteRelation .
// @router /relation/deleteRelation [POST]
func DeleteRelation(ctx context.Context, c *app.RequestContext) {
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

// GetRelationPaths .
// @router /relation/getRelationPaths [GET]
func GetRelationPaths(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.GetRelationPathsReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.GetRelationPathsResp)
	p := provider.Get()
	resp, err = p.RelationService.GetRelationPaths(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}
