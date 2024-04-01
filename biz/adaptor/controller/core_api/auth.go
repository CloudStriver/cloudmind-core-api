// Code generated by hertz generator.

package core_api

import (
	"context"
	"github.com/CloudStriver/cloudmind-core-api/biz/adaptor"
	core_api "github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/cloudmind-core-api/provider"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// Register .
// @router /auth/register [POST]
func Register(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.RegisterReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	resp := new(core_api.RegisterResp)
	p := provider.Get()
	resp, err = p.AuthService.Register(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// EmailLogin .
// @router /auth/login/email [POST]
func EmailLogin(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.EmailLoginReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.EmailLoginResp)
	p := provider.Get()
	resp, err = p.AuthService.EmailLogin(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// GithubLogin .
// @router /auth/login/github [POST]
func GithubLogin(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.GithubLoginReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.GithubLoginResp)
	p := provider.Get()
	resp, err = p.AuthService.GithubLogin(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// RefreshToken .
// @router /auth/refresh [POST]
func RefreshToken(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.RefreshTokenReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.RefreshTokenResp)
	p := provider.Get()
	resp, err = p.AuthService.RefreshToken(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// SendEmail .
// @router /auth/send [POST]
func SendEmail(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.SendEmailReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	resp := new(core_api.SendEmailResp)
	p := provider.Get()
	resp, err = p.AuthService.SendEmail(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// SetPasswordByEmail .
// @router /auth/reset/email [POST]
func SetPasswordByEmail(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.SetPasswordByEmailReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	resp := new(core_api.SetPasswordByEmailResp)
	p := provider.Get()
	resp, err = p.AuthService.SetPasswordByEmail(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// SetPasswordByPassword .
// @router /auth/reset/password [POST]
func SetPasswordByPassword(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.SetPasswordByPasswordReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	resp := new(core_api.SetPasswordByPasswordResp)
	p := provider.Get()
	resp, err = p.AuthService.SetPasswordByPassword(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// GiteeLogin .
// @router /auth/login/gitee [GET]
func GiteeLogin(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.GiteeLoginReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	resp := new(core_api.GiteeLoginResp)
	p := provider.Get()
	resp, err = p.AuthService.GiteeLogin(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// CheckEmail .
// @router /auth/checkEmail [GET]
func CheckEmail(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.CheckEmailReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.CheckEmailResp)
	p := provider.Get()
	resp, err = p.AuthService.CheckEmail(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// AskUploadAvatar .
// @router /auth/askUploadAvatar [POST]
func AskUploadAvatar(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.AskUploadAvatarReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.AskUploadAvatarResp)
	p := provider.Get()
	resp, err = p.AuthService.AskUploadAvatar(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// WeixinLogin .
// @router /auth/weixinLogin [GET]
func WeixinLogin(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.WeixinLoginReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.WeixinLoginResp)
	p := provider.Get()
	resp, err = p.AuthService.WeixinLogin(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// WeixinCallBack .
// @router /auth/weixinCallback [POST]
func WeixinCallBack(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.WeixinCallBackReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.WeixinCallBackResp)
	p := provider.Get()
	resp, err = p.AuthService.WeixinCallBack(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// WeixinIsLogin .
// @router /auth/weixinIsLogin [POST]
func WeixinIsLogin(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.WeixinIsLoginReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.WeixinIsLoginResp)
	p := provider.Get()
	resp, err = p.AuthService.WeixinIsLogin(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// QQLogin .
// @router /auth/qqLogin [GET]
func QQLogin(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.QQLoginReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.QQLoginResp)
	p := provider.Get()
	resp, err = p.AuthService.QQLogin(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}
