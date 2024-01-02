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

	p := provider.Get()
	resp, err := p.AuthService.Register(ctx, &req)
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

	p := provider.Get()
	resp, err := p.AuthService.EmailLogin(ctx, &req)
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

	p := provider.Get()
	resp, err := p.AuthService.GithubLogin(ctx, &req)
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

	p := provider.Get()
	resp, err := p.AuthService.RefreshToken(ctx, &req)
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

	p := provider.Get()
	resp, err := p.AuthService.SendEmail(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// GetCaptcha .
// @router /auth/captcha [GET]
func GetCaptcha(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.GetCaptchaReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	p := provider.Get()
	resp, err := p.AuthService.GetCaptcha(ctx, &req)
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

	p := provider.Get()
	resp, err := p.AuthService.SetPasswordByEmail(ctx, &req)
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

	p := provider.Get()
	resp, err := p.AuthService.SetPasswordByPassword(ctx, &req)
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

	p := provider.Get()
	resp, err := p.AuthService.GiteeLogin(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}
