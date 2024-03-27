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
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_sts"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/platform_comment"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/utils"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/sts"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/platform/comment"
	"github.com/google/wire"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/mr"
	"strings"
)

type IFileService interface {
	GetPrivateFile(ctx context.Context, req *core_api.GetPrivateFileReq) (resp *core_api.GetPrivateFileResp, err error)
	GetPublicFile(ctx context.Context, req *core_api.GetPublicFileReq) (resp *core_api.GetPublicFileResp, err error)
	GetPrivateFiles(ctx context.Context, req *core_api.GetPrivateFilesReq) (resp *core_api.GetPrivateFilesResp, err error)
	GetPublicFiles(ctx context.Context, req *core_api.GetPublicFilesReq) (resp *core_api.GetPublicFilesResp, err error)
	GetRecycleBinFiles(ctx context.Context, req *core_api.GetRecycleBinFilesReq) (resp *core_api.GetRecycleBinFilesResp, err error)
	GetFileBySharingCode(ctx context.Context, req *core_api.GetFileBySharingCodeReq) (resp *core_api.GetFileBySharingCodeResp, err error)
	CreateFile(ctx context.Context, req *core_api.CreateFileReq) (resp *core_api.CreateFileResp, err error)
	UpdateFile(ctx context.Context, req *core_api.UpdateFileReq) (resp *core_api.UpdateFileResp, err error)
	MoveFile(ctx context.Context, req *core_api.MoveFileReq) (resp *core_api.MoveFileResp, err error)
	DeleteFile(ctx context.Context, req *core_api.DeleteFileReq) (resp *core_api.DeleteFileResp, err error)
	EmptyRecycleBin(ctx context.Context, req *core_api.EmptyRecycleBinReq) (resp *core_api.EmptyRecycleBinResp, err error)
	CompletelyRemoveFile(ctx context.Context, req *core_api.CompletelyRemoveFileReq) (resp *core_api.CompletelyRemoveFileResp, err error)
	RecoverRecycleBinFile(ctx context.Context, req *core_api.RecoverRecycleBinFileReq) (resp *core_api.RecoverRecycleBinFileResp, err error)
	GetShareList(ctx context.Context, req *core_api.GetShareListReq) (resp *core_api.GetShareListResp, err error)
	CreateShareCode(ctx context.Context, req *core_api.CreateShareCodeReq) (resp *core_api.CreateShareCodeResp, err error)
	DeleteShareCode(ctx context.Context, req *core_api.DeleteShareCodeReq) (resp *core_api.DeleteShareCodeResp, err error)
	ParsingShareCode(ctx context.Context, req *core_api.ParsingShareCodeReq) (resp *core_api.ParsingShareCodeResp, err error)
	SaveFileToPrivateSpace(ctx context.Context, req *core_api.SaveFileToPrivateSpaceReq) (resp *core_api.SaveFileToPrivateSpaceResp, err error)
	AddFileToPublicSpace(ctx context.Context, req *core_api.AddFileToPublicSpaceReq) (resp *core_api.AddFileToPublicSpaceResp, err error)
	AskUploadFile(ctx context.Context, req *core_api.AskUploadFileReq) (resp *core_api.AskUploadFileResp, err error)
	AskDownloadFile(ctx context.Context, req *core_api.AskDownloadFileReq) (resp *core_api.AskDownloadFileResp, err error)
}

var FileServiceSet = wire.NewSet(
	wire.Struct(new(FileService), "*"),
	wire.Bind(new(IFileService), new(*FileService)),
)

type FileService struct {
	Config            *config.Config
	PlatformSts       cloudmind_sts.ICloudMindSts
	CloudMindContent  cloudmind_content.ICloudMindContent
	FileDomainService service.IFileDomainService
	PlatformComment   platform_comment.IPlatFormComment
}

func (s *FileService) AskDownloadFile(ctx context.Context, req *core_api.AskDownloadFileReq) (resp *core_api.AskDownloadFileResp, err error) {
	resp = new(core_api.AskDownloadFileResp)
	user := adaptor.ExtractUserMeta(ctx)
	if user.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	resp.Urls = make([]string, len(req.FileIds))
	if err = mr.Finish(lo.Map(req.FileIds, func(item string, i int) func() error {
		return func() error {
			getFileResp, err1 := s.CloudMindContent.GetFile(ctx, &content.GetFileReq{
				FileId:    item,
				IsGetSize: true,
			})
			if err1 != nil {
				return err1
			}
			if getFileResp.File.UserId != user.UserId {
				return consts.ErrNoAccessFile
			}
			genSignedUrlResp, err1 := s.PlatformSts.GenSignedUrl(ctx, &sts.GenSignedUrlReq{
				Path: getFileResp.File.Md5 + utils.GetFileNameSuffix(getFileResp.File.Name),
				Ttl:  getFileResp.File.SpaceSize,
			})
			if err1 != nil {
				return err1
			}
			resp.Urls[i] = genSignedUrlResp.SignedUrl
			return nil
		}
	})...); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *FileService) AskUploadFile(ctx context.Context, req *core_api.AskUploadFileReq) (resp *core_api.AskUploadFileResp, err error) {
	resp = new(core_api.AskUploadFileResp)
	user := adaptor.ExtractUserMeta(ctx)
	if user.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	//userId := user.GetUserId()
	//fmt.Println(userId)
	//if req.File {
	//	 TODO：扣取用户流量
	//}
	// TODO：判断是否已经上传过
	//getFileIsExist, err := s.CloudMindContent.GetFileIsExist(ctx, &content.GetFileIsExistReq{
	//	Md5: req.Md5,
	//})
	//if err != nil {
	//	return resp, err
	//}
	//
	//if getFileIsExist.Ok {
	//	return resp, nil
	//}

	genCosStsResp, err := s.PlatformSts.GenCosSts(ctx, &sts.GenCosStsReq{
		Path:   "users/" + req.Name,
		IsFile: true,
		Time:   req.FileSize / (1024 * 1024),
	})
	if err != nil {
		return resp, err
	}
	resp.SessionToken = genCosStsResp.SessionToken
	resp.TmpSecretId = genCosStsResp.SecretId
	resp.TmpSecretKey = genCosStsResp.SecretKey
	resp.StartTime = genCosStsResp.StartTime
	resp.ExpiredTime = genCosStsResp.ExpiredTime
	return resp, nil
}

func (s *FileService) GetPrivateFile(ctx context.Context, req *core_api.GetPrivateFileReq) (resp *core_api.GetPrivateFileResp, err error) {
	resp = new(core_api.GetPrivateFileResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	var res *content.GetFileResp
	if res, err = s.CloudMindContent.GetFile(ctx, &content.GetFileReq{FileId: req.FileId, IsGetSize: req.IsGetSize}); err != nil {
		return resp, err
	}
	switch {
	case res.File.UserId != userData.UserId:
		return resp, consts.ErrNoAccessFile
	case res.File.IsDel == consts.HardDel:
		return resp, consts.ErrFileNotExist
	}
	resp.File = convertor.FileToCorePrivateFile(res.File)
	return resp, nil
}

func (s *FileService) GetPublicFile(ctx context.Context, req *core_api.GetPublicFileReq) (resp *core_api.GetPublicFileResp, err error) {
	resp = new(core_api.GetPublicFileResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	var res *content.GetFileResp
	if res, err = s.CloudMindContent.GetFile(ctx, &content.GetFileReq{FileId: req.FileId, IsGetSize: req.IsGetSize}); err != nil {
		return resp, err
	}
	switch {
	case res.File.Zone == "" || res.File.SubZone == "":
		return resp, consts.ErrNoAccessFile
	case res.File.IsDel == consts.HardDel:
		return resp, consts.ErrFileNotExist
	}
	resp.File = convertor.FileToCorePublicFile(res.File)
	if err = mr.Finish(func() error {
		s.FileDomainService.LoadLikeCount(ctx, resp.File) // 点赞量
		return nil
	}, func() error {
		s.FileDomainService.LoadAuthor(ctx, resp.File, res.File.UserId) // 作者
		return nil
	}, func() error {
		s.FileDomainService.LoadViewCount(ctx, resp.File) // 浏览量
		return nil
	}, func() error {
		s.FileDomainService.LoadLiked(ctx, resp.File, userData.GetUserId()) // 是否点赞
		return nil
	}, func() error {
		s.FileDomainService.LoadCollected(ctx, resp.File, userData.GetUserId()) // 是否收藏
		return nil
	}, func() error {
		s.FileDomainService.LoadCollectCount(ctx, resp.File) // 收藏量
		return nil
	}, func() error {
		s.FileDomainService.LoadLabels(ctx, resp.File, res.File.Labels) // 标签集
		return nil
	}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *FileService) GetPrivateFiles(ctx context.Context, req *core_api.GetPrivateFilesReq) (resp *core_api.GetPrivateFilesResp, err error) {
	resp = new(core_api.GetPrivateFilesResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	sort := lo.ToPtr(content.SortOptions_SortOptions_createAtDesc)
	if req.SortType != nil {
		sort = lo.ToPtr(content.SortOptions(*req.SortType))
	}

	var res *content.GetFileListResp
	var searchOptions *content.SearchOptions
	p := convertor.MakePaginationOptions(req.Limit, req.Offset, req.LastToken, req.Backward)

	filter := &content.FileFilterOptions{
		OnlyUserId:   lo.ToPtr(userData.UserId),
		OnlyFatherId: req.OnlyFatherId,
		OnlyIsDel:    lo.ToPtr(consts.NotDel),
		OnlyType:     req.OnlyType,
		OnlyCategory: req.OnlyCategory,
	}
	if req.AllFieldsKey != nil {
		searchOptions = &content.SearchOptions{Type: &content.SearchOptions_AllFieldsKey{AllFieldsKey: *req.AllFieldsKey}}
	} else if req.Name != nil || req.Id != nil {
		searchOptions = &content.SearchOptions{Type: &content.SearchOptions_MultiFieldsKey{MultiFieldsKey: &content.SearchField{Name: req.Name, Id: req.Id}}}
	}

	if res, err = s.CloudMindContent.GetFileList(ctx, &content.GetFileListReq{SearchOptions: searchOptions, FilterOptions: filter, PaginationOptions: p, SortOptions: sort}); err != nil {
		return resp, err
	}
	resp.Files = lo.Map[*content.FileInfo, *core_api.PrivateFile](res.Files, func(item *content.FileInfo, _ int) *core_api.PrivateFile {
		return convertor.FileToCorePrivateFile(item)
	})

	resp.Token = res.Token
	resp.Total = res.Total
	resp.FatherNamePath = res.FatherNamePath
	resp.FatherIdPath = res.FatherIdPath
	return resp, nil
}

func (s *FileService) GetPublicFiles(ctx context.Context, req *core_api.GetPublicFilesReq) (resp *core_api.GetPublicFilesResp, err error) {
	resp = new(core_api.GetPublicFilesResp)
	userData := adaptor.ExtractUserMeta(ctx)
	var res *content.GetFileListResp
	var searchOptions *content.SearchOptions

	p := convertor.MakePaginationOptions(req.Limit, req.Offset, req.LastToken, req.Backward)
	filter := &content.FileFilterOptions{
		OnlyFatherId:     req.OnlyFatherId,
		OnlyZone:         req.OnlyZone,
		OnlySubZone:      req.OnlySubZone,
		OnlyIsDel:        lo.ToPtr(consts.NotDel),
		OnlyType:         req.OnlyType,
		OnlyCategory:     req.OnlyCategory,
		OnlyLabelId:      req.OnlyLabelId,
		OnlyDocumentType: lo.ToPtr(int64(core_api.DocumentType_DocumentType_public)),
	}
	sort := lo.ToPtr(content.SortOptions_SortOptions_createAtDesc)
	if req.SortType != nil {
		sort = lo.ToPtr(content.SortOptions(*req.SortType))
	}

	if req.AllFieldsKey != nil {
		searchOptions = &content.SearchOptions{Type: &content.SearchOptions_AllFieldsKey{AllFieldsKey: *req.AllFieldsKey}}
	} else if req.Name != nil || req.Id != nil || req.Description != nil {
		searchOptions = &content.SearchOptions{Type: &content.SearchOptions_MultiFieldsKey{MultiFieldsKey: &content.SearchField{Name: req.Name, Id: req.Id, Description: req.Description}}}
	}

	if res, err = s.CloudMindContent.GetFileList(ctx, &content.GetFileListReq{SearchOptions: searchOptions, FilterOptions: filter, PaginationOptions: p, SortOptions: sort}); err != nil {
		return resp, err
	}
	resp.Files = lo.Map[*content.FileInfo, *core_api.PublicFile](res.Files, func(item *content.FileInfo, _ int) *core_api.PublicFile {
		file := convertor.FileToCorePublicFile(item)
		_ = mr.Finish(func() error {
			s.FileDomainService.LoadLikeCount(ctx, file) // 点赞量
			return nil
		}, func() error {
			s.FileDomainService.LoadAuthor(ctx, file, item.UserId) // 作者
			return nil
		}, func() error {
			s.FileDomainService.LoadViewCount(ctx, file) // 浏览量
			return nil
		}, func() error {
			s.FileDomainService.LoadLiked(ctx, file, userData.GetUserId()) // 是否点赞
			return nil
		}, func() error {
			s.FileDomainService.LoadCollected(ctx, file, userData.GetUserId()) // 是否收藏
			return nil
		}, func() error {
			s.FileDomainService.LoadCollectCount(ctx, file) // 收藏量
			return nil
		})
		return file
	})
	resp.Token = res.Token
	resp.Total = res.Total
	resp.FatherNamePath = res.FatherNamePath
	resp.FatherIdPath = res.FatherIdPath
	return resp, nil
}

func (s *FileService) GetRecycleBinFiles(ctx context.Context, req *core_api.GetRecycleBinFilesReq) (resp *core_api.GetRecycleBinFilesResp, err error) {
	resp = new(core_api.GetRecycleBinFilesResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	var res *content.GetRecycleBinFilesResp
	p := convertor.MakePaginationOptions(req.Limit, req.Offset, req.LastToken, req.Backward)
	if res, err = s.CloudMindContent.GetRecycleBinFiles(ctx, &content.GetRecycleBinFilesReq{FilterOptions: &content.FileFilterOptions{OnlyUserId: lo.ToPtr(userData.UserId), OnlyIsDel: lo.ToPtr(consts.SoftDel)}, PaginationOptions: p}); err != nil {
		return resp, err
	}
	resp.Files = lo.Map[*content.FileInfo, *core_api.PrivateFile](res.Files, func(item *content.FileInfo, _ int) *core_api.PrivateFile {
		return convertor.FileToCorePrivateFile(item)
	})
	resp.Token = res.Token
	resp.Total = res.Total
	return resp, nil
}

func (s *FileService) GetFileBySharingCode(ctx context.Context, req *core_api.GetFileBySharingCodeReq) (resp *core_api.GetFileBySharingCodeResp, err error) {
	resp = new(core_api.GetFileBySharingCodeResp)
	var res *content.GetFileBySharingCodeResp
	var shareFile *content.ParsingShareCodeResp
	if shareFile, err = s.CloudMindContent.ParsingShareCode(ctx, &content.ParsingShareCodeReq{Code: req.Code}); err != nil {
		return resp, err
	}
	if shareFile.ShareFile.Key != req.Key {
		return resp, consts.ErrShareFileKey
	}
	p := convertor.MakePaginationOptions(req.Limit, req.Offset, req.LastToken, req.Backward)
	sort := lo.ToPtr(content.SortOptions_SortOptions_createAtDesc)
	if req.SortType != nil {
		sort = lo.ToPtr(content.SortOptions(*req.SortType))
	}
	if res, err = s.CloudMindContent.GetFileBySharingCode(ctx, &content.GetFileBySharingCodeReq{FileIds: shareFile.ShareFile.FileList, OnlyFatherId: req.OnlyFatherId, PaginationOptions: p, SortOptions: sort}); err != nil {
		return resp, err
	}
	resp.Files = lo.Map[*content.FileInfo, *core_api.PrivateFile](res.Files, func(item *content.FileInfo, _ int) *core_api.PrivateFile {
		return convertor.FileToCorePrivateFile(item)
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
	var father *content.GetFileResp
	file := &content.File{
		UserId:    userData.UserId,
		Name:      req.Name,
		Type:      req.Type,
		Category:  req.Category,
		FatherId:  req.FatherId,
		SpaceSize: req.SpaceSize,
		Md5:       req.Md5,
		IsDel:     consts.NotDel,
	}

	if file.UserId == file.FatherId {
		file.Path = file.UserId
	} else {
		if father, err = s.CloudMindContent.GetFile(ctx, &content.GetFileReq{FileId: req.FatherId, IsGetSize: false}); err != nil {
			return resp, err
		}
		if father.File.UserId != userData.UserId {
			return resp, consts.ErrNoAccessFile
		}
		if father.File.SpaceSize != consts.FolderSize {
			return resp, consts.ErrFileIsNotDir
		}
		file.Path = father.File.Path
	}

	if res, err = s.CloudMindContent.CreateFile(ctx, &content.CreateFileReq{File: file}); err != nil {
		return resp, err
	}
	resp.FileId = res.FileId
	resp.Name = res.Name
	return resp, nil
}

func (s *FileService) UpdateFile(ctx context.Context, req *core_api.UpdateFileReq) (resp *core_api.UpdateFileResp, err error) {
	resp = new(core_api.UpdateFileResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	file := &content.File{
		FileId: req.FileId,
		UserId: userData.UserId,
		Name:   req.Name,
	}
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

	var res *content.GetFilesByIdsResp
	if res, err = s.CloudMindContent.GetFilesByIds(ctx, &content.GetFilesByIdsReq{FileIds: []string{req.FileId, req.FatherId}}); err != nil {
		return resp, err
	}
	if req.FatherId == userData.UserId {
		res.Files[1] = &content.FileInfo{UserId: userData.UserId, Path: userData.UserId, SpaceSize: -1, IsDel: consts.NotDel}
	}
	for _, file := range res.Files {
		if file == nil {
			return resp, consts.ErrFileNotExist
		}
		switch {
		case file.UserId != userData.UserId:
			return resp, consts.ErrNoAccessFile
		case file.IsDel == consts.HardDel:
			return resp, consts.ErrFileNotExist
		}
	}

	switch {
	case res.Files[1].SpaceSize != consts.FolderSize:
		return resp, consts.ErrFileIsNotDir
	case strings.HasPrefix(res.Files[1].Path, res.Files[0].Path):
		return resp, consts.ErrIllegalOperation
	}

	if _, err = s.CloudMindContent.MoveFile(ctx, &content.MoveFileReq{
		FatherId:  req.FatherId,
		NewPath:   res.Files[1].Path,
		FileId:    res.Files[0].FileId,
		OldPath:   res.Files[0].Path,
		SpaceSize: res.Files[0].SpaceSize,
	}); err != nil {
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
	var res *content.GetFilesByIdsResp
	if res, err = s.CloudMindContent.GetFilesByIds(ctx, &content.GetFilesByIdsReq{FileIds: req.FileIds}); err != nil {
		return resp, err
	}

	files := make([]*content.FileParameter, len(req.FileIds))
	for i, file := range res.Files {
		if file.UserId != userData.UserId {
			return resp, consts.ErrNoAccessFile
		}
		switch req.DeleteType {
		case core_api.IsDel_Is_soft:
			if file.IsDel != consts.NotDel {
				return resp, consts.ErrIllegalOperation
			}
		case core_api.IsDel_Is_hard:
			if file.IsDel != consts.SoftDel {
				return resp, consts.ErrIllegalOperation
			}
		}
		files[i] = &content.FileParameter{
			FileId:    file.FileId,
			Path:      file.Path,
			SpaceSize: file.SpaceSize,
		}
	}

	if _, err = s.CloudMindContent.DeleteFile(ctx, &content.DeleteFileReq{
		DeleteType:     int64(req.DeleteType),
		ClearCommunity: req.ClearCommunity,
		Files:          files,
	}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *FileService) EmptyRecycleBin(ctx context.Context, req *core_api.EmptyRecycleBinReq) (resp *core_api.EmptyRecycleBinResp, err error) {
	resp = new(core_api.EmptyRecycleBinResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	if _, err = s.CloudMindContent.EmptyRecycleBin(ctx, &content.EmptyRecycleBinReq{UserId: userData.UserId}); err != nil {
		return resp, err
	}

	return resp, nil
}

func (s *FileService) CompletelyRemoveFile(ctx context.Context, req *core_api.CompletelyRemoveFileReq) (resp *core_api.CompletelyRemoveFileResp, err error) {
	resp = new(core_api.CompletelyRemoveFileResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	var res *content.GetFileResp
	if res, err = s.CloudMindContent.GetFile(ctx, &content.GetFileReq{FileId: req.FileId, IsGetSize: false}); err != nil {
		return resp, err
	}
	if res.File.UserId != userData.UserId {
		return resp, consts.ErrNoAccessFile
	}
	if _, err = s.CloudMindContent.CompletelyRemoveFile(ctx, &content.CompletelyRemoveFileReq{FileId: req.FileId, SpaceSize: res.File.SpaceSize, Path: res.File.Path}); err != nil {
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
	p := convertor.MakePaginationOptions(req.Limit, req.Offset, req.LastToken, req.Backward)
	if res, err = s.CloudMindContent.GetShareList(ctx, &content.GetShareListReq{ShareFileFilterOptions: &content.ShareFileFilterOptions{OnlyUserId: lo.ToPtr(userData.UserId)}, PaginationOptions: p}); err != nil {
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
	var files *content.GetFilesByIdsResp
	if files, err = s.CloudMindContent.GetFilesByIds(ctx, &content.GetFilesByIdsReq{FileIds: req.FileList}); err != nil {
		return resp, err
	}

	if len(files.Files) != len(req.FileList) { // 如果文件列表长度不一致， 说明文件列表存在问题 则返回错误
		return resp, consts.ErrIllegalOperation
	}

	for _, file := range files.Files {
		if file.UserId != userData.UserId {
			return resp, consts.ErrIllegalOperation
		}
		if file.IsDel != consts.NotDel {
			return resp, consts.ErrIllegalOperation
		}
	}

	if res, err = s.CloudMindContent.CreateShareCode(ctx, &content.CreateShareCodeReq{ShareFile: &content.ShareFile{
		UserId:        userData.UserId,
		Name:          req.Name,
		EffectiveTime: req.EffectiveTime,
		BrowseNumber:  consts.InitNumber,
		FileList:      req.FileList,
	}}); err != nil {
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
	if _, err = s.CloudMindContent.DeleteShareCode(ctx, &content.DeleteShareCodeReq{Code: req.OnlyCode}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *FileService) ParsingShareCode(ctx context.Context, req *core_api.ParsingShareCodeReq) (resp *core_api.ParsingShareCodeResp, err error) {
	resp = new(core_api.ParsingShareCodeResp)
	var res *content.ParsingShareCodeResp
	if res, err = s.CloudMindContent.ParsingShareCode(ctx, &content.ParsingShareCodeReq{Code: req.Code}); err != nil {
		return resp, err
	}
	switch {
	case res.ShareFile.Key != req.Key:
		return resp, consts.ErrShareFileKey
	case res.ShareFile.Status == int64(content.Validity_Validity_expired):
		return resp, nil
	}
	res.ShareFile.BrowseNumber++
	if _, err = s.CloudMindContent.UpdateShareCode(ctx, &content.UpdateShareCodeReq{ShareFile: res.ShareFile}); err != nil {
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

	if req.DocumentType == core_api.DocumentType_DocumentType_personal {
		var ok *content.CheckShareFileResp
		var res *content.ParsingShareCodeResp
		if res, err = s.CloudMindContent.ParsingShareCode(ctx, &content.ParsingShareCodeReq{Code: *req.Code}); err != nil {
			return resp, err
		}
		if res.ShareFile.Key != *req.Key {
			return resp, consts.ErrShareFileKey
		}
		ok, err = s.CloudMindContent.CheckShareFile(ctx, &content.CheckShareFileReq{FileIds: res.ShareFile.FileList, FileId: req.FileId})
		switch {
		case err != nil:
			return resp, err
		case !ok.Ok:
			return resp, consts.ErrNoAccessFile
		}
	}

	if req.FileId == req.FatherId {
		return resp, consts.ErrIllegalOperation
	} // 如果目标文件和要保存的文件是同一个用户的，则返回错误
	var files *content.GetFilesByIdsResp
	if files, err = s.CloudMindContent.GetFilesByIds(ctx, &content.GetFilesByIdsReq{FileIds: []string{req.FileId, req.FatherId}}); err != nil {
		return resp, err
	}
	if req.FatherId == userData.UserId {
		files.Files[1] = &content.FileInfo{UserId: userData.UserId, Path: userData.UserId, SpaceSize: -1, IsDel: consts.NotDel}
	}
	for _, file := range files.Files {
		if file == nil {
			return resp, consts.ErrFileNotExist
		}
	}

	switch {
	case files.Files[1].SpaceSize != consts.FolderSize:
		return resp, consts.ErrFileIsNotDir
	case files.Files[0].UserId == files.Files[1].UserId || userData.UserId != files.Files[1].UserId || req.FileId == req.FatherId:
		return resp, consts.ErrIllegalOperation
	case req.DocumentType == core_api.DocumentType_DocumentType_public && (files.Files[0].Zone == "" || files.Files[0].SubZone == ""):
		return resp, consts.ErrNoAccessFile
	}

	var res *content.SaveFileToPrivateSpaceResp
	if res, err = s.CloudMindContent.SaveFileToPrivateSpace(ctx, &content.SaveFileToPrivateSpaceReq{
		FileId:       files.Files[0].FileId,
		UserId:       userData.UserId,
		NewPath:      files.Files[1].Path,
		FatherId:     req.FatherId,
		Name:         files.Files[0].Name,
		Type:         files.Files[0].Type,
		FileMd5:      files.Files[0].Md5,
		SpaceSize:    files.Files[0].SpaceSize,
		DocumentType: int64(req.DocumentType),
	}); err != nil {
		return resp, err
	}
	resp.FileId = res.FileId
	resp.Name = res.Name
	return resp, nil
}

func (s *FileService) AddFileToPublicSpace(ctx context.Context, req *core_api.AddFileToPublicSpaceReq) (resp *core_api.AddFileToPublicSpaceResp, err error) {
	resp = new(core_api.AddFileToPublicSpaceResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	var file *content.GetFileResp
	if file, err = s.CloudMindContent.GetFile(ctx, &content.GetFileReq{FileId: req.FileId, IsGetSize: false}); err != nil {
		return resp, err
	}

	switch {
	case file.File.UserId != userData.UserId:
		return resp, consts.ErrNoAccessFile
	case file.File.IsDel != consts.NotDel:
		return resp, consts.ErrFileNotExist
	}

	file.File.Zone = req.Zone
	file.File.SubZone = req.SubZone
	file.File.Description = req.Description
	file.File.Labels = req.Labels
	err = mr.Finish(func() error {
		_, err1 := s.CloudMindContent.AddFileToPublicSpace(ctx, &content.AddFileToPublicSpaceReq{
			FileId:      file.File.FileId,
			Path:        file.File.Path,
			SpaceSize:   file.File.SpaceSize,
			Zone:        file.File.Zone,
			SubZone:     file.File.SubZone,
			Description: file.File.Description,
			Labels:      file.File.Labels,
		})
		return err1
	}, func() error {
		subject, _ := s.PlatformComment.GetCommentSubject(ctx, &comment.GetCommentSubjectReq{Id: file.File.FileId})
		if subject.GetSubject() != nil {
			return nil
		}
		_, err2 := s.PlatformComment.CreateCommentSubject(ctx, &comment.CreateCommentSubjectReq{Subject: &comment.Subject{
			Id:        file.File.FileId,
			UserId:    file.File.UserId,
			RootCount: lo.ToPtr(consts.InitNumber),
			AllCount:  lo.ToPtr(consts.InitNumber),
			State:     int64(core_api.State_Normal),
			Attrs:     int64(core_api.Attrs_None),
		}})
		return err2
	}, func() error {
		if _, err = s.CloudMindContent.CreateHot(ctx, &content.CreateHotReq{HotId: file.File.UserId}); err != nil {
			return err
		}
		return nil
	})
	return resp, err
}

func (s *FileService) RecoverRecycleBinFile(ctx context.Context, req *core_api.RecoverRecycleBinFileReq) (resp *core_api.RecoverRecycleBinFileResp, err error) {
	resp = new(core_api.RecoverRecycleBinFileResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	var res *content.GetFilesByIdsResp
	if res, err = s.CloudMindContent.GetFilesByIds(ctx, &content.GetFilesByIdsReq{FileIds: req.FileIds}); err != nil {
		return resp, err
	}

	files := make([]*content.FileParameter, len(req.FileIds))
	for i, file := range res.Files {
		switch {
		case file.UserId != userData.UserId:
			return resp, consts.ErrNoAccessFile
		case file.IsDel != consts.SoftDel:
			return resp, consts.ErrFileNotExist
		}
		files[i] = &content.FileParameter{
			FileId:    file.FileId,
			Path:      file.Path,
			SpaceSize: file.SpaceSize,
		}
	}

	if _, err = s.CloudMindContent.RecoverRecycleBinFile(ctx, &content.RecoverRecycleBinFileReq{Files: files}); err != nil {
		return resp, err
	}
	return resp, nil
}
