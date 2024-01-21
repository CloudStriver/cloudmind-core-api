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
	resp, err = p.ContentService.UpdateUser(ctx, &req)
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
	resp, err = p.ContentService.SearchUser(ctx, &req)
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
	//p := provider.Get()
	//resp, err = p.ContentService.GetFromRelations(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
	c.JSON(consts.StatusOK, resp)
}

// CreateUser .
// @router /content/user/create [POST]
func CreateUser(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.CreateUserReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.CreateUserResp)

	c.JSON(consts.StatusOK, resp)
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

	c.JSON(consts.StatusOK, resp)
}

// DeleteUser .
// @router /content/user/delete [POST]
func DeleteUser(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.DeleteUserReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.DeleteUserResp)

	c.JSON(consts.StatusOK, resp)
}

// GetFileIsExist .
// @router /content/file/exist [GET]
func GetFileIsExist(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.GetFileIsExistReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.GetFileIsExistResp)

	c.JSON(consts.StatusOK, resp)
}

// GetFile .
// @router /content/file/get [GET]
func GetFile(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.GetFileReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.GetFileResp)

	c.JSON(consts.StatusOK, resp)
}

// GetFileList .
// @router /content/file/list [GET]
func GetFileList(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.GetFileListReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.GetFileListResp)

	c.JSON(consts.StatusOK, resp)
}

// GetFolderSize .
// @router /content/folder/size [GET]
func GetFolderSize(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.GetFolderSizeReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.GetFolderSizeResp)

	c.JSON(consts.StatusOK, resp)
}

// GetFileCount .
// @router /content/file/count [GET]
func GetFileCount(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.GetFileCountReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.GetFileCountResp)

	c.JSON(consts.StatusOK, resp)
}

// GetFileBySharingCode .
// @router /content/sharecode/file [GET]
func GetFileBySharingCode(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.GetFileIsExistReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.GetFileIsExistResp)

	c.JSON(consts.StatusOK, resp)
}

// CreateFolder .
// @router /content/folder/create [POST]
func CreateFolder(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.CreateFolderReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.GetFileIsExistResp)

	c.JSON(consts.StatusOK, resp)
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

	c.JSON(consts.StatusOK, resp)
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

	c.JSON(consts.StatusOK, resp)
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

	c.JSON(consts.StatusOK, resp)
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

	c.JSON(consts.StatusOK, resp)
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

	c.JSON(consts.StatusOK, resp)
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

	c.JSON(consts.StatusOK, resp)
}

// CreateLabel .
// @router /content/label/create [POST]
func CreateLabel(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.CreateLabelReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.CreateLabelResp)

	c.JSON(consts.StatusOK, resp)
}

// UpdateLabel .
// @router /content/label/update [POST]
func UpdateLabel(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.UpdateLabelReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.UpdateLabelResp)

	c.JSON(consts.StatusOK, resp)
}

// GetLabel .
// @router /content/label/get [GET]
func GetLabel(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.GetLabelReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.GetLabelResp)

	c.JSON(consts.StatusOK, resp)
}

// DeleteLabel .
// @router /content/label/delete [POST]
func DeleteLabel(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.DeleteLabelReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.DeleteLabelResp)

	c.JSON(consts.StatusOK, resp)
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

	c.JSON(consts.StatusOK, resp)
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

	c.JSON(consts.StatusOK, resp)
}

// UpdateShareCode .
// @router /content/sharecode/update [POST]
func UpdateShareCode(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.UpdateShareCodeReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.UpdateShareCodeResp)

	c.JSON(consts.StatusOK, resp)
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

	c.JSON(consts.StatusOK, resp)
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

	c.JSON(consts.StatusOK, resp)
}

// DeleteShareFile .
// @router /content/sharefile/delete [POST]
func DeleteShareFile(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.DeleteShareFileReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.DeleteShareFileResp)

	c.JSON(consts.StatusOK, resp)
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

	c.JSON(consts.StatusOK, resp)
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

	c.JSON(consts.StatusOK, resp)
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

	c.JSON(consts.StatusOK, resp)
}

// GetPost .
// @router /content/post/get [GET]
func GetPost(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.GetPostReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.GetPostResp)

	c.JSON(consts.StatusOK, resp)
}

// GetPosts .
// @router /content/post/posts [GET]
func GetPosts(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.GetPostsReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.GetPostsResp)

	c.JSON(consts.StatusOK, resp)
}
