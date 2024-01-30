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

// UpdateUser .
// @router /content/user/update [GET]
func UpdateUser(ctx context.Context, c *app.RequestContext) {
	// this my demo
	var err error
	var req core_api.UpdateUserReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	resp := new(core_api.UpdateUserResp)
	p := provider.Get()
	resp, err = p.UserService.UpdateUser(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// SearchUser .
// @router /content/user/search [GET]
func SearchUser(ctx context.Context, c *app.RequestContext) {
	// this my demo
	var err error
	var req core_api.SearchUserReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	resp := new(core_api.SearchUserResp)
	p := provider.Get()
	resp, err = p.UserService.SearchUser(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// GetUser .
// @router /content/user/get [GET]
func GetUser(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.GetUserReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	resp := new(core_api.GetUserResp)
	p := provider.Get()
	resp, err = p.UserService.GetUser(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// GetUserDetail .
// @router /content/user/detail [GET]
func GetUserDetail(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.GetUserDetailReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.GetUserDetailResp)
	p := provider.Get()
	resp, err = p.UserService.GetUserDetail(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// GetFileBySharingCode .
// @router /content/sharecode/file [GET]
func GetFileBySharingCode(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.GetFileBySharingCodeReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.GetFileBySharingCodeResp)
	p := provider.Get()
	resp, err = p.ContentService.GetFileBySharingCode(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// UpdateFile .
// @router /content/file/update [POST]
func UpdateFile(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.UpdateFileReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.UpdateFileResp)
	p := provider.Get()
	resp, err = p.ContentService.UpdateFile(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// MoveFile .
// @router /content/file/move [POST]
func MoveFile(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.MoveFileReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.MoveFileResp)
	p := provider.Get()
	resp, err = p.ContentService.MoveFile(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// SaveFileToPrivateSpace .
// @router /content/file/save [POST]
func SaveFileToPrivateSpace(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.SaveFileToPrivateSpaceReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.SaveFileToPrivateSpaceResp)
	p := provider.Get()
	resp, err = p.ContentService.SaveFileToPrivateSpace(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// AddFileToPublicSpace .
// @router /content/file/add [POST]
func AddFileToPublicSpace(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.AddFileToPublicSpaceReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.AddFileToPublicSpaceResp)
	p := provider.Get()
	resp, err = p.ContentService.AddFileToPublicSpace(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// DeleteFile .
// @router /content/file/delete [POST]
func DeleteFile(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.DeleteFileReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.DeleteFileResp)
	p := provider.Get()
	resp, err = p.ContentService.DeleteFile(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// RecoverRecycleBinFile .
// @router /content/file/recover [POST]
func RecoverRecycleBinFile(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.RecoverRecycleBinFileReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.RecoverRecycleBinFileResp)
	p := provider.Get()
	resp, err = p.ContentService.RecoverRecycleBinFile(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// CreateShareCode .
// @router /content/sharecode/create [POST]
func CreateShareCode(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.CreateShareCodeReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.CreateShareCodeResp)
	p := provider.Get()
	resp, err = p.ContentService.CreateShareCode(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// GetShareList .
// @router /content/sharefile/list [GET]
func GetShareList(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.GetShareListReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.GetShareListResp)
	p := provider.Get()
	resp, err = p.ContentService.GetShareList(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// DeleteShareCode .
// @router /content/sharecode/delete [POST]
func DeleteShareCode(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.DeleteShareCodeReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.DeleteShareCodeResp)
	p := provider.Get()
	resp, err = p.ContentService.DeleteShareCode(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// ParsingShareCode .
// @router /content/sharecode/parsing [GET]
func ParsingShareCode(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.ParsingShareCodeReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.ParsingShareCodeResp)
	p := provider.Get()
	resp, err = p.ContentService.ParsingShareCode(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// CreatePost .
// @router /content/post/create [POST]
func CreatePost(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.CreatePostReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.CreatePostResp)
	p := provider.Get()
	resp, err = p.PostService.CreatePost(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// DeletePost .
// @router /content/post/delete [POST]
func DeletePost(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.DeletePostReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.DeletePostResp)
	p := provider.Get()
	resp, err = p.PostService.DeletePost(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// UpdatePost .
// @router /content/post/update [POST]
func UpdatePost(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.UpdatePostReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.UpdatePostResp)
	p := provider.Get()
	resp, err = p.PostService.UpdatePost(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

func GetPosts(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.GetPostsReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.GetPostsResp)
	p := provider.Get()
	resp, err = p.PostService.GetPosts(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// GetPost .
// @router /content/getPost [GET]
func GetPost(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.GetPostReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.GetPostResp)
	p := provider.Get()
	resp, err = p.PostService.GetPost(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// CreateFile .
// @router /content/createFile [POST]
func CreateFile(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.CreateFileReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.CreateFileResp)
	p := provider.Get()
	resp, err = p.ContentService.CreateFile(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// GetPrivateFiles .
// @router /content/getPrivateFiles [POST]
func GetPrivateFiles(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.GetPrivateFilesReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.GetPrivateFilesResp)
	p := provider.Get()
	resp, err = p.ContentService.GetPrivateFiles(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// GetPublicFiles .
// @router /content/getPublicFiles [POST]
func GetPublicFiles(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.GetPublicFilesReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.GetPublicFilesResp)
	p := provider.Get()
	resp, err = p.ContentService.GetPublicFiles(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// GetRecycleBinFiles .
// @router /content/getRecycleBinFiles [POST]
func GetRecycleBinFiles(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.GetRecycleBinFilesReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.GetRecycleBinFilesResp)
	p := provider.Get()
	resp, err = p.ContentService.GetRecycleBinFiles(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// CreateZone .
// @router /content/createZone [POST]
func CreateZone(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.CreateZoneReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.CreateZoneResp)
	p := provider.Get()
	resp, err = p.ZoneService.CreateZone(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// UpdateZone .
// @router /content/updateZone [POST]
func UpdateZone(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.UpdateZoneReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.UpdateZoneResp)
	p := provider.Get()
	resp, err = p.ZoneService.UpdateZone(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// GetZone .
// @router /content/getZone [GET]
func GetZone(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.GetZoneReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.GetZoneResp)
	p := provider.Get()
	resp, err = p.ZoneService.GetZone(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// DeleteZone .
// @router /content/deleteZone [POST]
func DeleteZone(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.DeleteZoneReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.DeleteZoneResp)
	p := provider.Get()
	resp, err = p.ZoneService.DeleteZone(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// GetPublicFile .
// @router /content/getPublicFile [POST]
func GetPublicFile(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.GetFileReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.GetFileResp)
	p := provider.Get()
	resp, err = p.ContentService.GetPublicFile(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// GetPrivateFile .
// @router /content/getPrivateFile [POST]
func GetPrivateFile(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.GetFileReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.GetFileResp)
	p := provider.Get()
	resp, err = p.ContentService.GetPublicFile(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}
