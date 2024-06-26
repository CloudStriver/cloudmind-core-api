package service

import (
	"context"
	"github.com/CloudStriver/cloudmind-core-api/biz/adaptor"
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/cloudmind-core-api/biz/domain/service"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/consts"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/convertor"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/kq"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_content"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_sts"
	platformservice "github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/platform"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/utils"
	"github.com/CloudStriver/cloudmind-mq/app/util/message"
	"github.com/CloudStriver/go-pkg/utils/pconvertor"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/sts"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/platform"
	"github.com/bytedance/sonic"
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
	MakeFilePrivate(ctx context.Context, req *core_api.MakeFilePrivateReq) (resp *core_api.MakeFilePrivateResp, err error)
	CheckFile(ctx context.Context, req *core_api.CheckFileReq) (resp *core_api.CheckFileResp, err error)
}

var FileServiceSet = wire.NewSet(
	wire.Struct(new(FileService), "*"),
	wire.Bind(new(IFileService), new(*FileService)),
)

type FileService struct {
	Config                *config.Config
	PlatformSts           cloudmind_sts.ICloudMindSts
	RelationDomainService service.IRelationDomainService
	CloudMindContent      cloudmind_content.ICloudMindContent
	FileDomainService     service.IFileDomainService
	Platform              platformservice.IPlatForm
	DeleteFileRelationKq  *kq.DeleteFileRelationKq
	CreateItemKq          *kq.CreateItemKq
	DeleteItemKq          *kq.DeleteItemKq
}

func (s *FileService) FilterContent(ctx context.Context, IsSure bool, contents []*string) ([]string, error) {
	cts := lo.Map[*string, string](contents, func(item *string, index int) string {
		return *item
	})
	if IsSure {
		replaceContentResp, err := s.PlatformSts.ReplaceContent(ctx, &sts.ReplaceContentReq{
			Contents: cts,
		})
		if err != nil {
			return nil, err
		}
		for i, val := range replaceContentResp.Content {
			*contents[i] = val
		}
		return nil, nil
	} else {
		// 内容检测
		findAllContentResp, err := s.PlatformSts.FindAllContent(ctx, &sts.FindAllContentReq{
			Contents: cts,
		})
		if err != nil {
			return nil, err
		}

		return findAllContentResp.Keywords, nil
	}
}

func (s *FileService) CheckFile(ctx context.Context, req *core_api.CheckFileReq) (resp *core_api.CheckFileResp, err error) {
	resp = new(core_api.CheckFileResp)
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	var res *content.GetPublicFileResp
	if res, err = s.CloudMindContent.GetPublicFile(ctx, &content.GetPublicFileReq{FileId: req.FileId, IsGetSize: false}); err != nil {
		return resp, err
	}

	switch {
	case res.AuditStatus != int64(content.AuditStatus_AuditStatus_wait):
		return resp, consts.ErrNoAuditStatus
	}

	resp.Keywords, err = s.FilterContent(ctx, false, []*string{&res.Name, &res.Description})
	if err != nil {
		return resp, err
	}

	var auditStatus int64
	if resp.Keywords == nil || len(resp.Keywords) == 0 {
		auditStatus = int64(content.AuditStatus_AuditStatus_pass)
	} else {
		auditStatus = int64(content.AuditStatus_AuditStatus_notPass)
	}

	_, err = s.CloudMindContent.UpdatePublicFile(ctx, &content.UpdatePublicFileReq{
		FileId:      req.FileId,
		AuditStatus: auditStatus,
	})
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (s *FileService) MakeFilePrivate(ctx context.Context, req *core_api.MakeFilePrivateReq) (resp *core_api.MakeFilePrivateResp, err error) {
	resp = new(core_api.MakeFilePrivateResp)
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	var getRes *content.GetPublicFilesByIdsResp
	if getRes, err = s.CloudMindContent.GetPublicFilesByIds(ctx, &content.GetPublicFilesByIdsReq{FileIds: req.FileIds}); err != nil {
		return resp, err
	}
	for _, file := range getRes.Files {
		if file.UserId != userData.UserId {
			return resp, consts.ErrNoAccessFile
		}
	}

	files := lo.Map(getRes.Files, func(item *content.PublicFile, _ int) *content.FileParameter {
		return &content.FileParameter{
			FileId:    item.FileId,
			Path:      item.Path,
			SpaceSize: item.SpaceSize,
		}
	})

	if _, err = s.CloudMindContent.MakeFilePrivate(ctx, &content.MakeFilePrivateReq{Files: files}); err != nil {
		return resp, err
	}

	for _, file := range files {
		data, _ := sonic.Marshal(&message.DeleteItemMessage{
			ItemId: file.FileId,
		})
		if err = s.DeleteItemKq.Push(pconvertor.Bytes2String(data)); err != nil {
			return resp, err
		}
	}

	return resp, nil
}

func (s *FileService) AskDownloadFile(ctx context.Context, req *core_api.AskDownloadFileReq) (resp *core_api.AskDownloadFileResp, err error) {
	resp = new(core_api.AskDownloadFileResp)
	user, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || user.GetUserId() == "" {
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
			if getFileResp.UserId != user.UserId {
				return consts.ErrNoAccessFile
			}
			genSignedUrlResp, err1 := s.PlatformSts.GenSignedUrl(ctx, &sts.GenSignedUrlReq{
				Path: getFileResp.Md5 + utils.GetFileNameSuffix(getFileResp.Name),
				Ttl:  getFileResp.SpaceSize,
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
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || userData.GetUserId() == "" {
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
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	var res *content.GetFileResp
	if res, err = s.CloudMindContent.GetFile(ctx, &content.GetFileReq{FileId: req.FileId, IsGetSize: true}); err != nil {
		return resp, err
	}
	switch {
	case res.UserId != userData.UserId:
		return resp, consts.ErrNoAccessFile
	case res.IsDel == consts.HardDel:
		return resp, consts.ErrFileNotExist
	}
	resp = &core_api.GetPrivateFileResp{
		Name:      res.Name,
		Type:      res.Type,
		Path:      res.Path,
		SpaceSize: res.SpaceSize,
		IsDel:     res.IsDel,
		CreateAt:  res.CreateAt,
		UpdateAt:  res.UpdateAt,
		DeleteAt:  res.DeleteAt,
	}
	return resp, nil
}

func (s *FileService) GetPublicFile(ctx context.Context, req *core_api.GetPublicFileReq) (resp *core_api.GetPublicFileResp, err error) {
	resp = new(core_api.GetPublicFileResp)
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil {
		return resp, err
	}

	var res *content.GetPublicFileResp
	if res, err = s.CloudMindContent.GetPublicFile(ctx, &content.GetPublicFileReq{FileId: req.FileId, IsGetSize: true}); err != nil {
		return resp, err
	}
	resp = &core_api.GetPublicFileResp{
		UserId:       res.UserId,
		Name:         res.Name,
		Type:         res.Type,
		Path:         res.Path,
		SpaceSize:    res.SpaceSize,
		Zone:         res.Zone,
		Description:  res.Description,
		CreateAt:     res.CreateAt,
		Labels:       []*core_api.Label{},
		Author:       &core_api.FileUser{},
		FileCount:    &core_api.FileCount{},
		FileRelation: &core_api.FileRelation{},
	}
	if err = mr.Finish(func() error {
		s.FileDomainService.LoadLikeCount(ctx, resp.FileCount, req.FileId) // 点赞量
		return nil
	}, func() error {
		s.FileDomainService.LoadAuthor(ctx, resp.Author, res.UserId) // 作者
		return nil
	}, func() error {
		s.FileDomainService.LoadViewCount(ctx, resp.FileCount, req.FileId) // 浏览量
		return nil
	}, func() error {
		s.FileDomainService.LoadLiked(ctx, resp.FileRelation, req.FileId, userData.GetUserId()) // 是否点赞
		return nil
	}, func() error {
		s.FileDomainService.LoadCollected(ctx, resp.FileRelation, req.FileId, userData.GetUserId()) // 是否收藏
		return nil
	}, func() error {
		s.FileDomainService.LoadCollectCount(ctx, resp.FileCount, req.FileId) // 收藏量
		return nil
	}, func() error {
		s.FileDomainService.LoadLabels(ctx, &resp.Labels, res.Labels) // 标签集
		return nil
	}, func() error {
		if err2 := s.RelationDomainService.CreateRelation(ctx, &core_api.Relation{
			FromType:     core_api.TargetType_UserType,
			FromId:       userData.UserId,
			ToType:       core_api.TargetType_FileType,
			ToId:         req.FileId,
			RelationType: core_api.RelationType_ViewRelationType,
		}); err2 != nil {
			return err2
		}
		return nil
	}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *FileService) GetPrivateFiles(ctx context.Context, req *core_api.GetPrivateFilesReq) (resp *core_api.GetPrivateFilesResp, err error) {
	resp = new(core_api.GetPrivateFilesResp)
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	var (
		res          *content.GetFileListResp
		searchOption *content.SearchOption
		sort         = lo.ToPtr(content.SortOptions_SortOptions_createAtDesc)
	)

	if req.SortType != nil {
		sort = lo.ToPtr(content.SortOptions(*req.SortType))
	}

	if req.SearchKeyword != nil {
		searchOption = &content.SearchOption{
			SearchKeyword:  req.SearchKeyword,
			SearchSortType: content.SearchSortType(*req.SearchType),
			SearchTimeType: content.SearchTimeType(*req.SearchTimerType),
		}
	}

	p := convertor.MakePaginationOptions(req.Limit, req.Offset, req.LastToken, req.Backward)
	filter := &content.FileFilterOptions{
		OnlyUserId:   lo.ToPtr(userData.UserId),
		OnlyFatherId: req.OnlyFatherId,
		OnlyIsDel:    lo.ToPtr(consts.NotDel),
		OnlyType:     req.OnlyType,
		OnlyCategory: req.OnlyCategory,
	}

	if res, err = s.CloudMindContent.GetFileList(ctx, &content.GetFileListReq{SearchOption: searchOption, FilterOptions: filter, PaginationOptions: p, SortOptions: sort}); err != nil {
		return resp, err
	}
	resp.Files = lo.Map[*content.File, *core_api.PrivateFile](res.Files, func(item *content.File, _ int) *core_api.PrivateFile {
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
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil {
		return resp, consts.ErrNotAuthentication
	}

	var (
		res          *content.GetPublicFileListResp
		searchOption *content.SearchOption
		sort         = lo.ToPtr(content.SortOptions_SortOptions_createAtDesc)
	)

	p := convertor.MakePaginationOptions(req.Limit, req.Offset, req.LastToken, req.Backward)
	filter := &content.PublicFileFilterOptions{
		OnlyUserId:      req.OnlyUserId,
		OnlyZone:        req.OnlyZone,
		OnlyType:        req.OnlyType,
		OnlyLabelId:     req.OnlyLabelId,
		OnlyAuditStatus: req.OnlyAuditStatus,
	}

	if req.SortType != nil {
		sort = lo.ToPtr(content.SortOptions(*req.SortType))
	}

	if req.SearchKeyword != nil {
		searchOption = &content.SearchOption{
			SearchKeyword:  req.SearchKeyword,
			SearchSortType: content.SearchSortType(*req.SearchType),
			SearchTimeType: content.SearchTimeType(*req.SearchTimerType),
		}
	}

	if res, err = s.CloudMindContent.GetPublicFileList(ctx, &content.GetPublicFileListReq{SearchOption: searchOption, FilterOptions: filter, PaginationOptions: p, SortOptions: sort}); err != nil {
		return resp, err
	}
	resp.Files = lo.Map(res.Files, func(item *content.PublicFile, _ int) *core_api.PublicFile {
		file := convertor.FileToCorePublicFile(item)
		_ = mr.Finish(func() error {
			s.FileDomainService.LoadLikeCount(ctx, file.FileCount, item.FileId) // 点赞量
			return nil
		}, func() error {
			s.FileDomainService.LoadAuthor(ctx, file.Author, item.UserId) // 作者
			return nil
		}, func() error {
			s.FileDomainService.LoadViewCount(ctx, file.FileCount, item.FileId) // 浏览量
			return nil
		}, func() error {
			s.FileDomainService.LoadLiked(ctx, file.FileRelation, item.FileId, userData.GetUserId()) // 是否点赞
			return nil
		}, func() error {
			s.FileDomainService.LoadCollected(ctx, file.FileRelation, item.FileId, userData.GetUserId()) // 是否收藏
			return nil
		}, func() error {
			s.FileDomainService.LoadCollectCount(ctx, file.FileCount, item.FileId) // 收藏量
			return nil
		})
		return file
	})
	resp.Token = res.Token
	resp.Total = res.Total
	return resp, nil
}

func (s *FileService) GetRecycleBinFiles(ctx context.Context, req *core_api.GetRecycleBinFilesReq) (resp *core_api.GetRecycleBinFilesResp, err error) {
	resp = new(core_api.GetRecycleBinFilesResp)
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	var (
		res          *content.GetRecycleBinFilesResp
		searchOption *content.SearchOption
	)

	if req.SearchKeyword != nil {
		searchOption = &content.SearchOption{
			SearchKeyword:  req.SearchKeyword,
			SearchSortType: content.SearchSortType(*req.SearchType),
			SearchTimeType: content.SearchTimeType(*req.SearchTimerType),
		}
	}

	p := convertor.MakePaginationOptions(req.Limit, req.Offset, req.LastToken, req.Backward)
	if res, err = s.CloudMindContent.GetRecycleBinFiles(ctx, &content.GetRecycleBinFilesReq{
		SearchOption: searchOption,
		FilterOptions: &content.FileFilterOptions{
			OnlyUserId: lo.ToPtr(userData.UserId),
			OnlyIsDel:  lo.ToPtr(consts.SoftDel),
		},
		PaginationOptions: p,
	}); err != nil {
		return resp, err
	}
	resp.Files = lo.Map[*content.File, *core_api.PrivateFile](res.Files, func(item *content.File, _ int) *core_api.PrivateFile {
		return convertor.FileToCorePrivateFile(item)
	})
	resp.Token = res.Token
	resp.Total = res.Total
	resp.FatherIdPath = res.FatherIdPath
	resp.FatherNamePath = res.FatherNamePath
	return resp, nil
}

func (s *FileService) GetFileBySharingCode(ctx context.Context, req *core_api.GetFileBySharingCodeReq) (resp *core_api.GetFileBySharingCodeResp, err error) {
	resp = new(core_api.GetFileBySharingCodeResp)
	var res *content.GetFileBySharingCodeResp
	var shareFile *content.ParsingShareCodeResp
	if shareFile, err = s.CloudMindContent.ParsingShareCode(ctx, &content.ParsingShareCodeReq{Code: req.Code}); err != nil {
		return resp, err
	}
	if shareFile.Key != req.Key {
		return resp, consts.ErrShareFileKey
	}
	if shareFile.Status != int64(content.Validity_Validity_expired) {
		return resp, consts.ErrShareCodeNotExist
	}
	p := convertor.MakePaginationOptions(req.Limit, req.Offset, req.LastToken, req.Backward)
	sort := lo.ToPtr(content.SortOptions_SortOptions_createAtDesc)
	if req.SortType != nil {
		sort = lo.ToPtr(content.SortOptions(*req.SortType))
	}
	if res, err = s.CloudMindContent.GetFileBySharingCode(ctx, &content.GetFileBySharingCodeReq{FileIds: shareFile.FileList, OnlyFatherId: req.OnlyFatherId, PaginationOptions: p, SortOptions: sort}); err != nil {
		return resp, err
	}
	resp.Files = lo.Map[*content.File, *core_api.PrivateFile](res.Files, func(item *content.File, _ int) *core_api.PrivateFile {
		return convertor.FileToCorePrivateFile(item)
	})
	resp.Token = res.Token
	resp.Total = res.Total
	return resp, nil
}

func (s *FileService) CreateFile(ctx context.Context, req *core_api.CreateFileReq) (resp *core_api.CreateFileResp, err error) {
	resp = new(core_api.CreateFileResp)
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	var (
		res    *content.CreateFileResp
		father *content.GetFileResp
		path   string
	)

	if userData.UserId == req.FatherId {
		path = userData.UserId
	} else {
		if father, err = s.CloudMindContent.GetFile(ctx, &content.GetFileReq{FileId: req.FatherId, IsGetSize: false}); err != nil {
			return resp, err
		}
		if father.UserId != userData.UserId {
			return resp, consts.ErrNoAccessFile
		}
		if father.SpaceSize != consts.FolderSize {
			return resp, consts.ErrFileIsNotDir
		}
		path = father.Path
	}

	if res, err = s.CloudMindContent.CreateFile(ctx, &content.CreateFileReq{
		UserId:    userData.UserId,
		Name:      req.Name,
		Type:      req.Type,
		FatherId:  req.FatherId,
		SpaceSize: req.SpaceSize,
		Md5:       req.Md5,
		Category:  req.Category,
		Path:      path,
	}); err != nil {
		return resp, err
	}
	resp.FileId = res.FileId
	resp.Name = res.Name
	return resp, nil
}

func (s *FileService) UpdateFile(ctx context.Context, req *core_api.UpdateFileReq) (resp *core_api.UpdateFileResp, err error) {
	resp = new(core_api.UpdateFileResp)
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	var res *content.UpdateFileResp
	if res, err = s.CloudMindContent.UpdateFile(ctx, &content.UpdateFileReq{
		FileId: req.FileId,
		Name:   req.Name,
	}); err != nil {
		return resp, err
	}
	resp.Name = res.Name
	return resp, nil
}

func (s *FileService) MoveFile(ctx context.Context, req *core_api.MoveFileReq) (resp *core_api.MoveFileResp, err error) {
	resp = new(core_api.MoveFileResp)
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	var res *content.GetFilesByIdsResp
	if res, err = s.CloudMindContent.GetFilesByIds(ctx, &content.GetFilesByIdsReq{FileIds: []string{req.FileId, req.FatherId}}); err != nil {
		return resp, err
	}
	if req.FatherId == userData.UserId {
		res.Files[1] = &content.File{UserId: userData.UserId, Path: userData.UserId, SpaceSize: consts.FolderSize, IsDel: consts.NotDel}
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
		Name:      res.Files[0].Name,
		IsDel:     res.Files[0].IsDel,
	}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *FileService) DeleteFile(ctx context.Context, req *core_api.DeleteFileReq) (resp *core_api.DeleteFileResp, err error) {
	resp = new(core_api.DeleteFileResp)
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || userData.GetUserId() == "" {
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

	if _, err1 := s.CloudMindContent.DeleteFile(ctx, &content.DeleteFileReq{
		DeleteType: int64(req.DeleteType),
		Files:      files,
	}); err1 != nil {
		return resp, err
	}

	return resp, nil
}

func (s *FileService) EmptyRecycleBin(ctx context.Context, _ *core_api.EmptyRecycleBinReq) (resp *core_api.EmptyRecycleBinResp, err error) {
	resp = new(core_api.EmptyRecycleBinResp)
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	if _, err = s.CloudMindContent.EmptyRecycleBin(ctx, &content.EmptyRecycleBinReq{UserId: userData.UserId}); err != nil {
		return resp, err
	}

	return resp, nil
}

func (s *FileService) CompletelyRemoveFile(ctx context.Context, req *core_api.CompletelyRemoveFileReq) (resp *core_api.CompletelyRemoveFileResp, err error) {
	resp = new(core_api.CompletelyRemoveFileResp)
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || userData.GetUserId() == "" {
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
		files[i] = &content.FileParameter{
			FileId:    file.FileId,
			Path:      file.Path,
			SpaceSize: file.SpaceSize,
		}
	}

	if _, err1 := s.CloudMindContent.CompletelyRemoveFile(ctx, &content.CompletelyRemoveFileReq{Files: files, UserId: userData.UserId}); err1 != nil {
		return resp, err1
	}

	return resp, nil
}

func (s *FileService) GetShareList(ctx context.Context, req *core_api.GetShareListReq) (resp *core_api.GetShareListResp, err error) {
	resp = new(core_api.GetShareListResp)
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	var res *content.GetShareListResp
	p := convertor.MakePaginationOptions(req.Limit, req.Offset, req.LastToken, req.Backward)
	if res, err = s.CloudMindContent.GetShareList(ctx, &content.GetShareListReq{FilterOptions: &content.ShareFileFilterOptions{OnlyUserId: lo.ToPtr(userData.UserId)}, PaginationOptions: p}); err != nil {
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
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || userData.GetUserId() == "" {
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

	if res, err = s.CloudMindContent.CreateShareCode(ctx, &content.CreateShareCodeReq{
		UserId:        userData.UserId,
		Name:          req.Name,
		EffectiveTime: req.EffectiveTime,
		FileList:      req.FileList,
	}); err != nil {
		return resp, err
	}
	resp.Code = res.Code
	resp.Key = res.Key
	return resp, nil
}

func (s *FileService) DeleteShareCode(ctx context.Context, req *core_api.DeleteShareCodeReq) (resp *core_api.DeleteShareCodeResp, err error) {
	resp = new(core_api.DeleteShareCodeResp)
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	if _, err = s.CloudMindContent.DeleteShareCode(ctx, &content.DeleteShareCodeReq{Code: req.Code}); err != nil {
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
	case res.Key != req.Key:
		return resp, consts.ErrShareFileKey
	case res.Status == int64(content.Validity_Validity_expired):
		return resp, nil
	}
	res.BrowseNumber++
	if _, err = s.CloudMindContent.UpdateShareCode(ctx, &content.UpdateShareCodeReq{
		Code:         req.Code,
		BrowseNumber: res.BrowseNumber,
	}); err != nil {
		return resp, err
	}
	resp = &core_api.ParsingShareCodeResp{
		UserId:        res.UserId,
		Name:          res.Name,
		EffectiveTime: res.EffectiveTime,
		BrowseNumber:  res.BrowseNumber,
		CreateAt:      res.CreateAt,
		FileList:      res.FileList,
	}
	return resp, nil
}

func (s *FileService) SaveFileToPrivateSpace(ctx context.Context, req *core_api.SaveFileToPrivateSpaceReq) (resp *core_api.SaveFileToPrivateSpaceResp, err error) {
	resp = new(core_api.SaveFileToPrivateSpaceResp)
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	if req.FileId == req.FatherId { // 如果目标文件和要保存的文件是，则返回错误
		return resp, consts.ErrIllegalOperation
	}

	if req.ShareField != nil {
		var (
			files *content.GetFilesByIdsResp
			check *content.CheckShareFileResp
			parse *content.ParsingShareCodeResp
			save  *content.SaveShareFileToPrivateSpaceResp
		)

		if parse, err = s.CloudMindContent.ParsingShareCode(ctx, &content.ParsingShareCodeReq{Code: req.ShareField.Code}); err != nil {
			return resp, err
		}

		if parse.Key != req.ShareField.Key {
			return resp, consts.ErrShareFileKey
		}

		check, err = s.CloudMindContent.CheckShareFile(ctx, &content.CheckShareFileReq{FileIds: parse.FileList, FileId: req.FileId})
		switch {
		case err != nil:
			return resp, err
		case !check.Ok:
			return resp, consts.ErrNoAccessFile
		}

		if files, err = s.CloudMindContent.GetFilesByIds(ctx, &content.GetFilesByIdsReq{FileIds: []string{req.FileId, req.FatherId}}); err != nil {
			return resp, err
		}
		if req.FatherId == userData.UserId {
			files.Files[1] = &content.File{Path: userData.UserId}
		}

		switch {
		case files.Files[0] == nil || files.Files[1] == nil:
			return resp, consts.ErrIllegalOperation
		case files.Files[1].UserId != userData.UserId:
			return resp, consts.ErrIllegalOperation
		}

		if save, err = s.CloudMindContent.SaveShareFileToPrivateSpace(ctx, &content.SaveShareFileToPrivateSpaceReq{
			FileId:    files.Files[0].FileId,
			UserId:    userData.UserId,
			NewPath:   files.Files[1].Path,
			FatherId:  req.FatherId,
			Name:      files.Files[0].Name,
			Type:      files.Files[0].Type,
			FileMd5:   files.Files[0].Md5,
			SpaceSize: files.Files[0].SpaceSize,
		}); err != nil {
			return resp, err
		}
		resp.FileId = save.FileId
		resp.Name = save.Name
	} else {
		var (
			err1, err2 error
			publicFile *content.GetPublicFileResp
			file       *content.GetFileResp
			save       *content.SavePublicFileToPrivateSpaceResp
		)

		if err = mr.Finish(func() error {
			if publicFile, err1 = s.CloudMindContent.GetPublicFile(ctx, &content.GetPublicFileReq{FileId: req.FileId}); err1 != nil {
				return err1
			}
			return nil
		}, func() error {
			if req.FatherId == userData.UserId {
				file = &content.GetFileResp{Path: userData.UserId}
			} else {
				if file, err2 = s.CloudMindContent.GetFile(ctx, &content.GetFileReq{FileId: req.FatherId}); err2 != nil {
					return err2
				}
			}
			return nil
		}); err != nil {
			return resp, err
		}

		switch {
		case file.UserId != userData.UserId:
			return resp, consts.ErrIllegalOperation
		}

		if save, err = s.CloudMindContent.SavePublicFileToPrivateSpace(ctx, &content.SavePublicFileToPrivateSpaceReq{
			FileId:    req.FileId,
			UserId:    userData.UserId,
			NewPath:   file.Path,
			FatherId:  req.FatherId,
			Name:      publicFile.Name,
			Type:      publicFile.Type,
			FileMd5:   publicFile.Md5,
			SpaceSize: publicFile.SpaceSize,
		}); err != nil {
			return resp, err
		}
		resp.FileId = save.FileId
		resp.Name = save.Name
	}

	return resp, nil
}

func (s *FileService) AddFileToPublicSpace(ctx context.Context, req *core_api.AddFileToPublicSpaceReq) (resp *core_api.AddFileToPublicSpaceResp, err error) {
	resp = new(core_api.AddFileToPublicSpaceResp)
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	var (
		addRes *content.AddFileToPublicSpaceResp
		getRes *content.GetFileResp
	)

	if getRes, err = s.CloudMindContent.GetFile(ctx, &content.GetFileReq{FileId: req.FileId, IsGetSize: false}); err != nil {
		return resp, err
	}

	switch {
	case getRes.UserId != userData.UserId:
		return resp, consts.ErrNoAccessFile
	case getRes.IsDel != consts.NotDel:
		return resp, consts.ErrFileNotExist
	}

	if addRes, err = s.CloudMindContent.AddFileToPublicSpace(ctx, &content.AddFileToPublicSpaceReq{
		FileId:      req.FileId,
		Zone:        req.Zone,
		Description: req.Description,
		Labels:      req.Labels,
	}); err != nil {
		return resp, err
	}

	err = mr.Finish(func() error {
		subject, _ := s.Platform.GetCommentSubject(ctx, &platform.GetCommentSubjectReq{SubjectId: addRes.FileId})
		if subject != nil {
			return nil
		}
		if _, err1 := s.Platform.CreateCommentSubject(ctx, &platform.CreateCommentSubjectReq{
			SubjectId: addRes.FileId,
			UserId:    userData.UserId,
		}); err1 != nil {
			return err1
		}
		return nil
	}, func() error {
		if _, err1 := s.CloudMindContent.CreateHot(ctx, &content.CreateHotReq{HotId: addRes.FileId}); err1 != nil {
			return err1
		}
		return nil
	}, func() error {
		if err1 := s.RelationDomainService.CreateRelation(ctx, &core_api.Relation{
			FromType:     core_api.TargetType_UserType,
			FromId:       userData.UserId,
			ToType:       core_api.TargetType_FileType,
			ToId:         req.FileId,
			RelationType: core_api.RelationType_PublishRelationType,
		}); err1 != nil {
			return err1
		}
		return nil
	}, func() error {
		data, _ := sonic.Marshal(&message.CreateItemMessage{
			ItemId:   addRes.FileId,
			IsHidden: false,
			Labels:   req.Labels,
			Category: core_api.Category_name[int32(core_api.Category_PostCategory)],
		})
		if err1 := s.CreateItemKq.Push(pconvertor.Bytes2String(data)); err1 != nil {
			return err1
		}
		return nil
	})
	return resp, err
}

func (s *FileService) RecoverRecycleBinFile(ctx context.Context, req *core_api.RecoverRecycleBinFileReq) (resp *core_api.RecoverRecycleBinFileResp, err error) {
	resp = new(core_api.RecoverRecycleBinFileResp)
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || userData.GetUserId() == "" {
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
