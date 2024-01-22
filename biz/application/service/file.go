package service

import (
	"context"
	"github.com/CloudStriver/cloudmind-core-api/biz/adaptor"
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/cloudmind-core-api/biz/domain/service"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/consts"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/convertor"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_content"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
	"github.com/google/wire"
	"github.com/samber/lo"
)

type IFileService interface {
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
}

var FileServiceSet = wire.NewSet(
	wire.Struct(new(FileService), "*"),
	wire.Bind(new(IFileService), new(*FileService)),
)

type FileService struct {
	Config            *config.Config
	CloudMindContent  cloudmind_content.ICloudMindContent
	PostDomainService service.IPostDomainService
}

func (s *FileService) GetFileIsExist(ctx context.Context, req *core_api.GetFileIsExistReq) (resp *core_api.GetFileIsExistResp, err error) {
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

func (s *FileService) GetFile(ctx context.Context, req *core_api.GetFileReq) (resp *core_api.GetFileResp, err error) {
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

func (s *FileService) GetFileList(ctx context.Context, req *core_api.GetFileListReq) (resp *core_api.GetFileListResp, err error) {
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

func (s *FileService) GetFileBySharingCode(ctx context.Context, req *core_api.GetFileBySharingCodeReq) (resp *core_api.GetFileBySharingCodeResp, err error) {
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

func (s *FileService) GetFolderSize(ctx context.Context, req *core_api.GetFolderSizeReq) (resp *core_api.GetFolderSizeResp, err error) {
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

func (s *FileService) CreateFolder(ctx context.Context, req *core_api.CreateFolderReq) (resp *core_api.CreateFolderResp, err error) {
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

func (s *FileService) UpdateFile(ctx context.Context, req *core_api.UpdateFileReq) (resp *core_api.UpdateFileResp, err error) {
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

func (s *FileService) MoveFile(ctx context.Context, req *core_api.MoveFileReq) (resp *core_api.MoveFileResp, err error) {
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

func (s *FileService) DeleteFile(ctx context.Context, req *core_api.DeleteFileReq) (resp *core_api.DeleteFileResp, err error) {
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

func (s *FileService) GetShareList(ctx context.Context, req *core_api.GetShareListReq) (resp *core_api.GetShareListResp, err error) {
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

func (s *FileService) CreateShareCode(ctx context.Context, req *core_api.CreateShareCodeReq) (resp *core_api.CreateShareCodeResp, err error) {
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

func (s *FileService) UpdateShareCode(ctx context.Context, req *core_api.UpdateShareCodeReq) (resp *core_api.UpdateShareCodeResp, err error) {
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

func (s *FileService) DeleteShareCode(ctx context.Context, req *core_api.DeleteShareCodeReq) (resp *core_api.DeleteShareCodeResp, err error) {
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

func (s *FileService) ParsingShareCode(ctx context.Context, req *core_api.ParsingShareCodeReq) (resp *core_api.ParsingShareCodeResp, err error) {
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

func (s *FileService) SaveFileToPrivateSpace(ctx context.Context, req *core_api.SaveFileToPrivateSpaceReq) (resp *core_api.SaveFileToPrivateSpaceResp, err error) {
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

func (s *FileService) AddFileToPublicSpace(ctx context.Context, req *core_api.AddFileToPublicSpaceReq) (resp *core_api.AddFileToPublicSpaceResp, err error) {
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
