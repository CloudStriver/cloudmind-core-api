package service

import (
	"context"
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

type IContentService interface {
	GetFileIsExist(ctx context.Context, req *core_api.GetFileIsExistReq) (resp *core_api.GetFileIsExistResp, err error)
	GetFile(ctx context.Context, req *core_api.GetFileReq) (resp *core_api.GetFileResp, err error)
	GetFileList(ctx context.Context, req *core_api.GetFileListReq) (resp *core_api.GetFileListResp, err error)
	GetFileBySharingCode(ctx context.Context, req *core_api.GetFileBySharingCodeReq) (resp *core_api.GetFileBySharingCodeResp, err error)
	GetFolderSize(ctx context.Context, req *core_api.GetFolderSizeReq) (resp *core_api.GetFolderSizeResp, err error)
	CreateFolder(ctx context.Context, req *core_api.CreateFolderReq) (resp *core_api.CreateFolderResp, err error)
	UpdateFile(ctx context.Context, req *core_api.UpdateFileReq) (resp *core_api.UpdateFileResp, err error)
	MoveFile(ctx context.Context, req *core_api.MoveFileReq) (resp *core_api.MoveFileResp, err error)
	DeleteFile(ctx context.Context, req *core_api.DeleteFileReq) (resp *core_api.DeleteFileResp, err error)
	GetShareList(ctx context.Context, req *core_api.GetShareListReq) (resp *core_api.GetShareListResp, err error)
	CreateShareCode(ctx context.Context, req *core_api.CreateShareCodeReq) (resp *core_api.CreateShareCodeResp, err error)
	UpdateShareCode(ctx context.Context, req *core_api.UpdateShareCodeReq) (resp *core_api.UpdateShareCodeResp, err error)
	DeleteShareCode(ctx context.Context, req *core_api.DeleteShareCodeReq) (resp *core_api.DeleteShareCodeResp, err error)
	ParsingShareCode(ctx context.Context, req *core_api.ParsingShareCodeReq) (resp *core_api.ParsingShareCodeResp, err error)
	SaveFileToPrivateSpace(ctx context.Context, req *core_api.SaveFileToPrivateSpaceReq) (resp *core_api.SaveFileToPrivateSpaceResp, err error)
	AddFileToPublicSpace(ctx context.Context, req *core_api.AddFileToPublicSpaceReq) (resp *core_api.AddFileToPublicSpaceResp, err error)
	GetLabel(ctx context.Context, req *core_api.GetLabelReq) (resp *core_api.GetLabelResp, err error)
	CreateLabel(ctx context.Context, req *core_api.CreateLabelReq) (resp *core_api.CreateLabelResp, err error)
	UpdateLabel(ctx context.Context, req *core_api.UpdateLabelReq) (resp *core_api.UpdateLabelResp, err error)
	DeleteLabel(ctx context.Context, req *core_api.DeleteLabelReq) (resp *core_api.DeleteLabelResp, err error)
	CreatePost(ctx context.Context, req *core_api.CreatePostReq) (resp *core_api.CreatePostResp, err error)
	GetPost(ctx context.Context, req *core_api.GetPostReq) (resp *core_api.GetPostResp, err error)
	GetPosts(ctx context.Context, req *core_api.GetPostsReq) (resp *core_api.GetPostsResp, err error)
	UpdatePost(ctx context.Context, req *core_api.UpdatePostReq) (resp *core_api.UpdatePostResp, err error)
	DeletePost(ctx context.Context, req *core_api.DeletePostReq) (resp *core_api.DeletePostResp, err error)
	GetUser(ctx context.Context, req *core_api.GetUserReq) (resp *core_api.GetUserResp, err error)
	CreateUser(ctx context.Context, req *core_api.CreateUserReq) (resp *core_api.CreateUserResp, err error)
	UpdateUser(ctx context.Context, req *core_api.UpdateUserReq) (resp *core_api.UpdateUserResp, err error)
	GetUserDetail(ctx context.Context, req *core_api.GetUserDetailReq) (resp *core_api.GetUserDetailResp, err error)
	SearchUser(ctx context.Context, req *core_api.SearchUserReq) (resp *core_api.SearchUserResp, err error)
	DeleteUser(ctx context.Context, req *core_api.DeleteUserReq) (resp *core_api.DeleteUserResp, err error)
}

var ContentServiceSet = wire.NewSet(
	wire.Struct(new(ContentService), "*"),
	wire.Bind(new(IContentService), new(*ContentService)),
)

type ContentService struct {
	Config           *config.Config
	CloudMindContent cloudmind_content.ICloudMindContent
}

func (s *ContentService) GetFileIsExist(ctx context.Context, req *core_api.GetFileIsExistReq) (resp *core_api.GetFileIsExistResp, err error) {
	resp = new(core_api.GetFileIsExistResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	res := new(content.GetFileIsExistResp)
	if res, err = s.CloudMindContent.GetFileIsExist(ctx, &content.GetFileIsExistReq{Md5: req.Md5}); err != nil {
		return resp, err
	}
	resp.Ok = res.Ok
	return resp, nil
}

func (s *ContentService) GetFile(ctx context.Context, req *core_api.GetFileReq) (resp *core_api.GetFileResp, err error) {
	resp = new(core_api.GetFileResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	filter := convertor.FilterOptionsToFilterOptions(req.FilterOptions)
	filter.OnlyUserId = lo.ToPtr(userData.UserId)
	res := new(content.GetFileResp)
	if res, err = s.CloudMindContent.GetFile(ctx, &content.GetFileReq{FilterOptions: filter, IsGetSize: req.IsGetSize}); err != nil {
		return resp, err
	}
	resp.File = convertor.FileToCoreFile(res.File)
	return resp, nil
}

func (s *ContentService) GetFileList(ctx context.Context, req *core_api.GetFileListReq) (resp *core_api.GetFileListResp, err error) {
	resp = new(core_api.GetFileListResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	res := new(content.GetFileListResp)
	p := convertor.PaginationOptionsToPaginationOptions(req.PaginationOptions)
	if req.SearchOptions == nil {
		searchOptions := convertor.SearchOptionsToFileSearchOptions(req.SearchOptions)
		if res, err = s.CloudMindContent.GetFileList(ctx, &content.GetFileListReq{SearchOptions: searchOptions, PaginationOptions: p}); err != nil {
			return resp, err
		}
	} else {
		filter := convertor.FilterOptionsToFilterOptions(req.FilterOptions)
		filter.OnlyUserId = lo.ToPtr(userData.UserId)
		if res, err = s.CloudMindContent.GetFileList(ctx, &content.GetFileListReq{FilterOptions: filter, PaginationOptions: p}); err != nil {
			return resp, err
		}
	}
	resp.Files = make([]*core_api.FileInfo, 0, len(res.Files))
	for _, file := range res.Files {
		resp.Files = append(resp.Files, convertor.FileToCoreFile(file))
	}
	resp.Token = res.Token
	resp.Total = res.Total
	return resp, nil
}

func (s *ContentService) GetFileBySharingCode(ctx context.Context, req *core_api.GetFileBySharingCodeReq) (resp *core_api.GetFileBySharingCodeResp, err error) {
	resp = new(core_api.GetFileBySharingCodeResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	res := new(content.GetFileBySharingCodeResp)
	filter := convertor.FilterOptionsToFilterOptions(req.FilterOptions)
	p := convertor.PaginationOptionsToPaginationOptions(req.PaginationOptions)
	if res, err = s.CloudMindContent.GetFileBySharingCode(ctx, &content.GetFileBySharingCodeReq{SharingCode: req.SharingCode, FilterOptions: filter, PaginationOptions: p}); err != nil {
		return resp, err
	}
	resp.Files = make([]*core_api.FileInfo, 0, len(res.Files))
	for _, file := range res.Files {
		resp.Files = append(resp.Files, convertor.FileToCoreFile(file))
	}
	resp.Token = res.Token
	resp.Total = res.Total
	return resp, nil
}

func (s *ContentService) GetFolderSize(ctx context.Context, req *core_api.GetFolderSizeReq) (resp *core_api.GetFolderSizeResp, err error) {
	resp = new(core_api.GetFolderSizeResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	res := new(content.GetFolderSizeResp)
	filter := convertor.FilterOptionsToFilterOptions(req.FilterOptions)
	if res, err = s.CloudMindContent.GetFolderSize(ctx, &content.GetFolderSizeReq{FilterOptions: filter}); err != nil {
		return resp, err
	}
	resp.SpaceSize = res.SpaceSize
	return resp, nil
}

func (s *ContentService) CreateFolder(ctx context.Context, req *core_api.CreateFolderReq) (resp *core_api.CreateFolderResp, err error) {
	resp = new(core_api.CreateFolderResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	res := new(content.CreateFolderResp)
	req.File.UserId = userData.UserId
	file := convertor.CoreFileToFile(req.File)
	if res, err = s.CloudMindContent.CreateFolder(ctx, &content.CreateFolderReq{File: file}); err != nil {
		return resp, err
	}
	resp.FileId = res.FileId
	return resp, nil
}

func (s *ContentService) UpdateFile(ctx context.Context, req *core_api.UpdateFileReq) (resp *core_api.UpdateFileResp, err error) {
	resp = new(core_api.UpdateFileResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	req.File.UserId = userData.UserId
	file := convertor.CoreFileToFile(req.File)
	if _, err = s.CloudMindContent.UpdateFile(ctx, &content.UpdateFileReq{File: file}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *ContentService) MoveFile(ctx context.Context, req *core_api.MoveFileReq) (resp *core_api.MoveFileResp, err error) {
	resp = new(core_api.MoveFileResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	if _, err = s.CloudMindContent.MoveFile(ctx, &content.MoveFileReq{UserId: userData.UserId, FileId: req.FileId, FatherId: req.FatherId}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *ContentService) DeleteFile(ctx context.Context, req *core_api.DeleteFileReq) (resp *core_api.DeleteFileResp, err error) {
	resp = new(core_api.DeleteFileResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	if _, err = s.CloudMindContent.DeleteFile(ctx, &content.DeleteFileReq{UserId: userData.UserId, FileId: req.FileId, DeleteType: content.IsDel(req.DeleteType), ClearCommunity: req.ClearCommunity}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *ContentService) GetShareList(ctx context.Context, req *core_api.GetShareListReq) (resp *core_api.GetShareListResp, err error) {
	resp = new(core_api.GetShareListResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	res := new(content.GetShareListResp)
	shareOptions := convertor.ShareOptionsToShareOptions(req.ShareFileFilterOptions)
	p := convertor.PaginationOptionsToPaginationOptions(req.PaginationOptions)
	if res, err = s.CloudMindContent.GetShareList(ctx, &content.GetShareListReq{ShareFileFilterOptions: shareOptions, PaginationOptions: p}); err != nil {
		return resp, err
	}
	resp.ShareCodes = make([]*core_api.ShareCode, 0, len(res.ShareCodes))
	for _, shareCode := range res.ShareCodes {
		resp.ShareCodes = append(resp.ShareCodes, convertor.ShareCodeToShareCode(shareCode))
	}
	resp.Token = res.Token
	resp.Total = res.Total
	return resp, nil
}

func (s *ContentService) CreateShareCode(ctx context.Context, req *core_api.CreateShareCodeReq) (resp *core_api.CreateShareCodeResp, err error) {
	resp = new(core_api.CreateShareCodeResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	res := new(content.CreateShareCodeResp)
	req.ShareFile.UserId = userData.UserId
	sharefile := convertor.CoreShareFileToShareFile(req.ShareFile)
	if res, err = s.CloudMindContent.CreateShareCode(ctx, &content.CreateShareCodeReq{ShareFile: sharefile}); err != nil {
		return resp, err
	}
	resp.Code = res.Code
	return resp, nil
}

func (s *ContentService) UpdateShareCode(ctx context.Context, req *core_api.UpdateShareCodeReq) (resp *core_api.UpdateShareCodeResp, err error) {
	resp = new(core_api.UpdateShareCodeResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	req.ShareFile.UserId = userData.UserId
	sharefile := convertor.CoreShareFileToShareFile(req.ShareFile)
	if _, err = s.CloudMindContent.UpdateShareCode(ctx, &content.UpdateShareCodeReq{ShareFile: sharefile}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *ContentService) DeleteShareCode(ctx context.Context, req *core_api.DeleteShareCodeReq) (resp *core_api.DeleteShareCodeResp, err error) {
	resp = new(core_api.DeleteShareCodeResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	req.ShareFileFilterOptions.OnlyUserId = lo.ToPtr(userData.UserId)
	shareOptions := convertor.ShareOptionsToShareOptions(req.ShareFileFilterOptions)
	if _, err = s.CloudMindContent.DeleteShareCode(ctx, &content.DeleteShareCodeReq{ShareFileFilterOptions: shareOptions}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *ContentService) ParsingShareCode(ctx context.Context, req *core_api.ParsingShareCodeReq) (resp *core_api.ParsingShareCodeResp, err error) {
	resp = new(core_api.ParsingShareCodeResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	res := new(content.ParsingShareCodeResp)
	if res, err = s.CloudMindContent.ParsingShareCode(ctx, &content.ParsingShareCodeReq{Code: req.Code}); err != nil {
		return resp, err
	}
	resp.ShareFile = convertor.ShareFileToCoreShareFile(res.ShareFile)
	return resp, nil
}

func (s *ContentService) SaveFileToPrivateSpace(ctx context.Context, req *core_api.SaveFileToPrivateSpaceReq) (resp *core_api.SaveFileToPrivateSpaceResp, err error) {
	resp = new(core_api.SaveFileToPrivateSpaceResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	res := new(content.SaveFileToPrivateSpaceResp)
	if res, err = s.CloudMindContent.SaveFileToPrivateSpace(ctx, &content.SaveFileToPrivateSpaceReq{UserId: userData.UserId, FileId: req.FileId, FatherId: req.FatherId, DocumentType: content.DocumentType(req.DocumentType)}); err != nil {
		return resp, err
	}
	resp.FileId = res.FileId
	return resp, nil
}

func (s *ContentService) AddFileToPublicSpace(ctx context.Context, req *core_api.AddFileToPublicSpaceReq) (resp *core_api.AddFileToPublicSpaceResp, err error) {
	resp = new(core_api.AddFileToPublicSpaceResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	req.File.UserId = userData.UserId
	file := convertor.CoreFileToFile(req.File)
	if _, err = s.CloudMindContent.AddFileToPublicSpace(ctx, &content.AddFileToPublicSpaceReq{File: file}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *ContentService) GetLabel(ctx context.Context, req *core_api.GetLabelReq) (resp *core_api.GetLabelResp, err error) {
	resp = new(core_api.GetLabelResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	res := new(content.GetLabelResp)
	if res, err = s.CloudMindContent.GetLabel(ctx, &content.GetLabelReq{Id: req.Id}); err != nil {
		return resp, err
	}
	resp.Label = convertor.LabelToCoreLabel(res.Label)
	return resp, nil
}

func (s *ContentService) CreateLabel(ctx context.Context, req *core_api.CreateLabelReq) (resp *core_api.CreateLabelResp, err error) {
	resp = new(core_api.CreateLabelResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	res := new(content.CreateLabelResp)
	label := convertor.CoreLabelToLabel(req.Label)
	if res, err = s.CloudMindContent.CreateLabel(ctx, &content.CreateLabelReq{Label: label}); err != nil {
		return resp, err
	}
	resp.Id = res.Id
	return resp, nil
}

func (s *ContentService) UpdateLabel(ctx context.Context, req *core_api.UpdateLabelReq) (resp *core_api.UpdateLabelResp, err error) {
	resp = new(core_api.UpdateLabelResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	label := convertor.CoreLabelToLabel(req.Label)
	if _, err = s.CloudMindContent.UpdateLabel(ctx, &content.UpdateLabelReq{Label: label}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *ContentService) DeleteLabel(ctx context.Context, req *core_api.DeleteLabelReq) (resp *core_api.DeleteLabelResp, err error) {
	resp = new(core_api.DeleteLabelResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	if _, err = s.CloudMindContent.DeleteLabel(ctx, &content.DeleteLabelReq{Id: req.Id}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *ContentService) CreatePost(ctx context.Context, req *core_api.CreatePostReq) (resp *core_api.CreatePostResp, err error) {
	resp = new(core_api.CreatePostResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	req.PostInfo.UserId = userData.UserId
	post := convertor.CorePostInfoToPostInfo(req.PostInfo)
	if _, err = s.CloudMindContent.CreatePost(ctx, &content.CreatePostReq{PostInfo: post}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *ContentService) GetPost(ctx context.Context, req *core_api.GetPostReq) (resp *core_api.GetPostResp, err error) {
	resp = new(core_api.GetPostResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	if _, err = s.CloudMindContent.GetPost(ctx, &content.GetPostReq{PostId: req.PostId}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *ContentService) GetPosts(ctx context.Context, req *core_api.GetPostsReq) (resp *core_api.GetPostsResp, err error) {
	resp = new(core_api.MoveFileResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	if _, err = s.CloudMindContent.MoveFile(ctx, &content.MoveFileReq{UserId: req.UserId, FileId: req.FileId, FatherId: req.FatherId}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *ContentService) UpdatePost(ctx context.Context, req *core_api.UpdatePostReq) (resp *core_api.UpdatePostResp, err error) {
	resp = new(core_api.DeleteFileResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	if _, err = s.CloudMindContent.DeleteFile(ctx, &content.DeleteFileReq{UserId: userData.UserId, FileId: req.FileId, DeleteType: content.IsDel(req.DeleteType), ClearCommunity: req.ClearCommunity}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *ContentService) DeletePost(ctx context.Context, req *core_api.DeletePostReq) (resp *core_api.DeletePostResp, err error) {
	resp = new(core_api.MoveFileResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	if _, err = s.CloudMindContent.MoveFile(ctx, &content.MoveFileReq{UserId: req.UserId, FileId: req.FileId, FatherId: req.FatherId}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *ContentService) GetUser(ctx context.Context, req *core_api.GetUserReq) (resp *core_api.GetUserResp, err error) {
	resp = new(core_api.MoveFileResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	if _, err = s.CloudMindContent.MoveFile(ctx, &content.MoveFileReq{UserId: req.UserId, FileId: req.FileId, FatherId: req.FatherId}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *ContentService) CreateUser(ctx context.Context, req *core_api.CreateUserReq) (resp *core_api.CreateUserResp, err error) {
	resp = new(core_api.MoveFileResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	if _, err = s.CloudMindContent.MoveFile(ctx, &content.MoveFileReq{UserId: req.UserId, FileId: req.FileId, FatherId: req.FatherId}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *ContentService) UpdateUser(ctx context.Context, req *core_api.UpdateUserReq) (resp *core_api.UpdateUserResp, err error) {
	resp = new(core_api.UpdateUserResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	user := convertor.UserToUserDetailInfo(req.UserDetail)
	user.UserId = userData.UserId

	if _, err = s.CloudMindContent.UpdateUser(ctx, &content.UpdateUserReq{UserDetailInfo: user}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *ContentService) GetUserDetail(ctx context.Context, req *core_api.GetUserDetailReq) (resp *core_api.GetUserDetailResp, err error) {
	resp = new(core_api.MoveFileResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	if _, err = s.CloudMindContent.MoveFile(ctx, &content.MoveFileReq{UserId: req.UserId, FileId: req.FileId, FatherId: req.FatherId}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *ContentService) SearchUser(ctx context.Context, req *core_api.SearchUserReq) (resp *core_api.SearchUserResp, err error) {
	resp = new(core_api.SearchUserResp)
	users, err := s.CloudMindContent.SearchUser(ctx, &content.SearchUserReq{
		Keyword:           req.Keyword,
		PaginationOptions: convertor.PaginationOptionsToPaginationOptions(req.PaginationOptions),
	})
	if err != nil {
		return resp, err
	}
	resp.Users = make([]*core_api.User, 0, len(users.Users))
	for _, user := range users.Users {
		resp.Users = append(resp.Users, convertor.UserDetailToUser(user))
	}
	resp.Total = users.Total
	resp.LastToken = users.LastToken
	return resp, nil
}

func (s *ContentService) DeleteUser(ctx context.Context, req *core_api.DeleteUserReq) (resp *core_api.DeleteUserResp, err error) {
	resp = new(core_api.MoveFileResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	if _, err = s.CloudMindContent.MoveFile(ctx, &content.MoveFileReq{UserId: req.UserId, FileId: req.FileId, FatherId: req.FatherId}); err != nil {
		return resp, err
	}
	return resp, nil
}
