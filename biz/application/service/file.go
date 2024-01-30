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
	"github.com/zeromicro/go-zero/core/mr"
)

type IFileService interface {
	GetPrivateFile(ctx context.Context, req *core_api.GetFileReq) (resp *core_api.GetFileResp, err error)
	GetPublicFile(ctx context.Context, req *core_api.GetFileReq) (resp *core_api.GetFileResp, err error)
	GetPrivateFiles(ctx context.Context, req *core_api.GetPrivateFilesReq) (resp *core_api.GetPrivateFilesResp, err error)
	GetPublicFiles(ctx context.Context, req *core_api.GetPublicFilesReq) (resp *core_api.GetPublicFilesResp, err error)
	GetRecycleBinFiles(ctx context.Context, req *core_api.GetRecycleBinFilesReq) (resp *core_api.GetRecycleBinFilesResp, err error)
	GetFileBySharingCode(ctx context.Context, req *core_api.GetFileBySharingCodeReq) (resp *core_api.GetFileBySharingCodeResp, err error)
	CreateFile(ctx context.Context, req *core_api.CreateFileReq) (resp *core_api.CreateFileResp, err error)
	UpdateFile(ctx context.Context, req *core_api.UpdateFileReq) (resp *core_api.UpdateFileResp, err error)
	MoveFile(ctx context.Context, req *core_api.MoveFileReq) (resp *core_api.MoveFileResp, err error)
	DeleteFile(ctx context.Context, req *core_api.DeleteFileReq) (resp *core_api.DeleteFileResp, err error)
	RecoverRecycleBinFile(ctx context.Context, req *core_api.RecoverRecycleBinFileReq) (resp *core_api.RecoverRecycleBinFileResp, err error)
	GetShareList(ctx context.Context, req *core_api.GetShareListReq) (resp *core_api.GetShareListResp, err error)
	CreateShareCode(ctx context.Context, req *core_api.CreateShareCodeReq) (resp *core_api.CreateShareCodeResp, err error)
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
	FileDomainService service.IFileDomainService
}

func (s *FileService) GetPrivateFile(ctx context.Context, req *core_api.GetFileReq) (resp *core_api.GetFileResp, err error) {
	resp = new(core_api.GetFileResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	var res *content.GetFileResp
	req.FilterOptions.OnlyUserId = lo.ToPtr(userData.UserId)
	filter := convertor.FilterOptionsToFilterOptions(req.FilterOptions)
	if res, err = s.CloudMindContent.GetFile(ctx, &content.GetFileReq{FilterOptions: filter, IsGetSize: req.IsGetSize}); err != nil {
		return resp, err
	}
	resp.File = convertor.FileToCoreFile(res.File)
	return resp, nil
}

func (s *FileService) GetPublicFile(ctx context.Context, req *core_api.GetFileReq) (resp *core_api.GetFileResp, err error) {
	resp = new(core_api.GetFileResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	var res *content.GetFileResp
	req.FilterOptions.OnlyDocumentType = lo.ToPtr(int64(core_api.DocumentType_DocumentType_public))
	filter := convertor.FilterOptionsToFilterOptions(req.FilterOptions)
	if res, err = s.CloudMindContent.GetFile(ctx, &content.GetFileReq{FilterOptions: filter, IsGetSize: req.IsGetSize}); err != nil {
		return resp, err
	}
	resp.File = convertor.FileToCoreFile(res.File)
	_ = mr.Finish(func() error {
		s.FileDomainService.LoadLikeCount(ctx, resp.File)
		return nil
	}, func() error {
		s.FileDomainService.LoadAuthor(ctx, resp.File, res.File.UserId)
		return nil
	}, func() error {
		s.FileDomainService.LoadCollectCount(ctx, resp.File)
		return nil
	}, func() error {
		s.FileDomainService.LoadLiked(ctx, resp.File, res.File.UserId)
		return nil
	})
	return resp, nil
}

func (s *FileService) GetPrivateFiles(ctx context.Context, req *core_api.GetPrivateFilesReq) (resp *core_api.GetPrivateFilesResp, err error) {
	resp = new(core_api.GetPrivateFilesResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	var res *content.GetFileListResp
	p := convertor.PaginationOptionsToPaginationOptions(req.PaginationOptions)
	req.FilterOptions.OnlyIsDel = lo.ToPtr(consts.NotDel)
	req.FilterOptions.OnlyUserId = lo.ToPtr(userData.UserId)
	filter := convertor.FilterOptionsToFilterOptions(req.FilterOptions)
	if req.SearchOptions != nil {
		searchOptions := convertor.SearchOptionsToFileSearchOptions(req.SearchOptions)
		if res, err = s.CloudMindContent.GetFileList(ctx, &content.GetFileListReq{SearchOptions: searchOptions, FilterOptions: filter, PaginationOptions: p}); err != nil {
			return resp, err
		}
	} else {
		if res, err = s.CloudMindContent.GetFileList(ctx, &content.GetFileListReq{FilterOptions: filter, PaginationOptions: p}); err != nil {
			return resp, err
		}
	}
	resp.Files = lo.Map[*content.FileInfo, *core_api.FileInfo](res.Files, func(item *content.FileInfo, _ int) *core_api.FileInfo {
		return convertor.FileToCoreFile(item)
	})
	resp.Token = res.Token
	resp.Total = res.Total
	return resp, nil
}

func (s *FileService) GetPublicFiles(ctx context.Context, req *core_api.GetPublicFilesReq) (resp *core_api.GetPublicFilesResp, err error) {
	resp = new(core_api.GetPublicFilesResp)
	var res *content.GetFileListResp
	p := convertor.PaginationOptionsToPaginationOptions(req.PaginationOptions)
	req.FilterOptions.OnlyDocumentType = lo.ToPtr(int64(core_api.DocumentType_DocumentType_public))
	req.FilterOptions.OnlyIsDel = lo.ToPtr(consts.NotDel)
	filter := convertor.FilterOptionsToFilterOptions(req.FilterOptions)
	if req.SearchOptions != nil {
		searchOptions := convertor.SearchOptionsToFileSearchOptions(req.SearchOptions)
		if res, err = s.CloudMindContent.GetFileList(ctx, &content.GetFileListReq{SearchOptions: searchOptions, FilterOptions: filter, PaginationOptions: p}); err != nil {
			return resp, err
		}
	} else {
		if res, err = s.CloudMindContent.GetFileList(ctx, &content.GetFileListReq{FilterOptions: filter, PaginationOptions: p}); err != nil {
			return resp, err
		}
	}
	resp.Files = lo.Map[*content.FileInfo, *core_api.FileInfo](res.Files, func(item *content.FileInfo, _ int) *core_api.FileInfo {
		return convertor.FileToCoreFile(item)
	})
	resp.Token = res.Token
	resp.Total = res.Total
	return resp, nil
}

func (s *FileService) GetRecycleBinFiles(ctx context.Context, req *core_api.GetRecycleBinFilesReq) (resp *core_api.GetRecycleBinFilesResp, err error) {
	resp = new(core_api.GetRecycleBinFilesResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	var res *content.GetFileListResp
	p := convertor.PaginationOptionsToPaginationOptions(req.PaginationOptions)
	req.FilterOptions.OnlyUserId = lo.ToPtr(userData.UserId)
	req.FilterOptions.OnlyIsDel = lo.ToPtr(consts.SoftDel)
	filter := convertor.FilterOptionsToFilterOptions(req.FilterOptions)
	if res, err = s.CloudMindContent.GetFileList(ctx, &content.GetFileListReq{FilterOptions: filter, PaginationOptions: p}); err != nil {
		return resp, err
	}
	resp.Files = lo.Map[*content.FileInfo, *core_api.FileInfo](res.Files, func(item *content.FileInfo, _ int) *core_api.FileInfo {
		return convertor.FileToCoreFile(item)
	})
	resp.Token = res.Token
	resp.Total = res.Total
	return resp, nil
}

func (s *FileService) GetFileBySharingCode(ctx context.Context, req *core_api.GetFileBySharingCodeReq) (resp *core_api.GetFileBySharingCodeResp, err error) {
	resp = new(core_api.GetFileBySharingCodeResp)
	var res *content.GetFileBySharingCodeResp
	filter := convertor.FilterOptionsToFilterOptions(req.FilterOptions)
	p := convertor.PaginationOptionsToPaginationOptions(req.PaginationOptions)
	if res, err = s.CloudMindContent.GetFileBySharingCode(ctx, &content.GetFileBySharingCodeReq{SharingCode: req.SharingCode, FilterOptions: filter, PaginationOptions: p}); err != nil {
		return resp, err
	}

	resp.Files = lo.Map[*content.FileInfo, *core_api.FileInfo](res.Files, func(item *content.FileInfo, _ int) *core_api.FileInfo {
		return convertor.FileToCoreFile(item)
	})
	resp.Token = res.Token
	resp.Total = res.Total
	return resp, nil
}

func (s *FileService) CreateFile(ctx context.Context, req *core_api.CreateFileReq) (resp *core_api.CreateFileResp, err error) {
	resp = new(core_api.CreateFileResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	var res *content.CreateFileResp
	req.File.UserId = userData.UserId
	file := convertor.CoreFileToFile(req.File)
	file.IsDel = int64(core_api.IsDel_Is_no)
	if res, err = s.CloudMindContent.CreateFile(ctx, &content.CreateFileReq{File: file}); err != nil {
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

	var res *content.GetShareListResp
	shareOptions := convertor.ShareOptionsToShareOptions(req.ShareFileFilterOptions)
	p := convertor.PaginationOptionsToPaginationOptions(req.PaginationOptions)
	if res, err = s.CloudMindContent.GetShareList(ctx, &content.GetShareListReq{ShareFileFilterOptions: shareOptions, PaginationOptions: p}); err != nil {
		return resp, err
	}

	resp.ShareCodes = lo.Map[*content.ShareCode, *core_api.ShareCode](res.ShareCodes, func(item *content.ShareCode, _ int) *core_api.ShareCode {
		return convertor.ShareCodeToCoreShareCode(item)
	})
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

	var res *content.CreateShareCodeResp
	req.ShareFile.UserId = userData.UserId
	sharefile := convertor.CoreShareFileToShareFile(req.ShareFile)
	if res, err = s.CloudMindContent.CreateShareCode(ctx, &content.CreateShareCodeReq{ShareFile: sharefile}); err != nil {
		return resp, err
	}
	resp.Code = res.Code
	resp.Key = res.Key
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

	var res *content.ParsingShareCodeResp
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

	var res *content.SaveFileToPrivateSpaceResp
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

func (s *FileService) RecoverRecycleBinFile(ctx context.Context, req *core_api.RecoverRecycleBinFileReq) (resp *core_api.RecoverRecycleBinFileResp, err error) {
	resp = new(core_api.RecoverRecycleBinFileResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	if _, err = s.CloudMindContent.RecoverRecycleBinFile(ctx, &content.RecoverRecycleBinFileReq{FileId: req.FileId, UserId: userData.UserId}); err != nil {
		return resp, err
	}
	return resp, nil
}
