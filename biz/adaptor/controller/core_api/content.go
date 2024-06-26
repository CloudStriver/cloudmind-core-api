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
	resp, err = p.FileService.GetFileBySharingCode(ctx, &req)
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
	resp, err = p.FileService.UpdateFile(ctx, &req)
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
	resp, err = p.FileService.MoveFile(ctx, &req)
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
	resp, err = p.FileService.SaveFileToPrivateSpace(ctx, &req)
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
	resp, err = p.FileService.AddFileToPublicSpace(ctx, &req)
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
	resp, err = p.FileService.DeleteFile(ctx, &req)
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
	resp, err = p.FileService.RecoverRecycleBinFile(ctx, &req)
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
	resp, err = p.FileService.CreateShareCode(ctx, &req)
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
	resp, err = p.FileService.GetShareList(ctx, &req)
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
	resp, err = p.FileService.DeleteShareCode(ctx, &req)
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
	resp, err = p.FileService.ParsingShareCode(ctx, &req)
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

// GetPosts .
// @router /content/getPosts [GET]
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
	resp, err = p.FileService.CreateFile(ctx, &req)
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
	resp, err = p.FileService.GetPrivateFiles(ctx, &req)
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
	resp, err = p.FileService.GetPublicFiles(ctx, &req)
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
	resp, err = p.FileService.GetRecycleBinFiles(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// GetPublicFile .
// @router /content/getPublicFile [POST]
func GetPublicFile(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.GetPublicFileReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.GetPublicFileResp)
	p := provider.Get()
	resp, err = p.FileService.GetPublicFile(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// GetPrivateFile .
// @router /content/getPrivateFile [POST]
func GetPrivateFile(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.GetPrivateFileReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.GetPrivateFileResp)
	p := provider.Get()
	resp, err = p.FileService.GetPrivateFile(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// AskUploadFile .
// @router /content/askUploadFile [POST]
func AskUploadFile(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.AskUploadFileReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.AskUploadFileResp)
	p := provider.Get()
	resp, err = p.FileService.AskUploadFile(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// AskDownloadFile .
// @router /content/askDownloadFile [POST]
func AskDownloadFile(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.AskDownloadFileReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.AskDownloadFileResp)
	p := provider.Get()
	resp, err = p.FileService.AskDownloadFile(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// CompletelyRemoveFile .
// @router /content/completelyRemoveFile [POST]
func CompletelyRemoveFile(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.CompletelyRemoveFileReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.CompletelyRemoveFileResp)
	p := provider.Get()
	resp, err = p.FileService.CompletelyRemoveFile(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// GetRecommendByUser .
// @router /content/getRecommendByUser [GET]
func GetRecommendByUser(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.GetRecommendByUserReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.GetRecommendByUserResp)
	p := provider.Get()
	resp, err = p.RecommendService.GetRecommendByUser(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// GetRecommendByItem .
// @router /content/getRecommendByItem [GET]
func GetRecommendByItem(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.GetRecommendByItemReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.GetRecommendByItemResp)
	p := provider.Get()
	resp, err = p.RecommendService.GetRecommendByItem(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// CreateFeedBack .
// @router /content/createFeedBack [POST]
func CreateFeedBack(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.CreateFeedBackReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.CreateFeedBackResp)
	p := provider.Get()
	resp, err = p.RecommendService.CreateFeedBack(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// GetPopularRecommend .
// @router /content/getPopularRecommend [GET]
func GetPopularRecommend(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.GetPopularRecommendReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.GetPopularRecommendResp)
	p := provider.Get()
	resp, err = p.RecommendService.GetPopularRecommend(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// GetLatestRecommend .
// @router /content/getLatestRecommend [GET]
func GetLatestRecommend(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.GetLatestRecommendReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.GetLatestRecommendResp)
	p := provider.Get()
	resp, err = p.RecommendService.GetLatestRecommend(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// CreateProduct .
// @router /content/createProduct [POST]
func CreateProduct(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.CreateProductReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.CreateProductResp)
	p := provider.Get()
	resp, err = p.ProductService.CreateProduct(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// GetProduct .
// @router /content/getProduct [GET]
func GetProduct(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.GetProductReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.GetProductResp)
	p := provider.Get()
	resp, err = p.ProductService.GetProduct(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// GetProducts .
// @router /content/getProducts [GET]
func GetProducts(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.GetProductsReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.GetProductsResp)
	p := provider.Get()
	resp, err = p.ProductService.GetProducts(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// UpdateProduct .
// @router /content/updateProduct [POST]
func UpdateProduct(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.UpdateProductReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.UpdateProductResp)
	p := provider.Get()
	resp, err = p.ProductService.UpdateProduct(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// DeleteProduct .
// @router /content/deleteProduct [POST]
func DeleteProduct(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.DeleteProductReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.DeleteProductResp)
	p := provider.Get()
	resp, err = p.ProductService.DeleteProduct(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// EmptyRecycleBin .
// @router /content/emptyRecycleBin [POST]
func EmptyRecycleBin(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.EmptyRecycleBinReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.EmptyRecycleBinResp)
	p := provider.Get()
	resp, err = p.FileService.EmptyRecycleBin(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// MakeFilePrivate .
// @router content/makeFilePrivate [POST]
func MakeFilePrivate(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.MakeFilePrivateReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.MakeFilePrivateResp)
	p := provider.Get()
	resp, err = p.FileService.MakeFilePrivate(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}

// CheckFile .
// @router /content/checkFile [GET]
func CheckFile(ctx context.Context, c *app.RequestContext) {
	var err error
	var req core_api.CheckFileReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(core_api.CheckFileResp)
	p := provider.Get()
	resp, err = p.FileService.CheckFile(ctx, &req)
	adaptor.PostProcess(ctx, c, &req, resp, err)
}
