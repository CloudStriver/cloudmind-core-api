package convertor

import (
	dto_basic "github.com/CloudStriver/cloudmind-core-api/biz/application/dto/basic"
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/basic"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/platform/relation"
)

func UserToUserDetailInfo(req *core_api.UserDetail) *content.User {
	return &content.User{
		Name:        req.Name,
		Sex:         req.Sex,
		FullName:    req.FullName,
		IdCard:      req.IdCard,
		Description: req.Description,
		Url:         req.Url,
	}
}

//
//func FileToCoreFile(req *content.FileInfo) *core_api.FileInfo {
//	return &core_api.FileInfo{
//		FileId:      req.FileId,
//		UserId:      req.UserId,
//		Name:        req.Name,
//		Type:        req.Type,
//		Path:        req.Path,
//		FatherId:    req.FatherId,
//		SpaceSize:   req.SpaceSize,
//		Md5:         req.Md5,
//		IsDel:       req.IsDel,
//		Zone:        req.Zone,
//		SubZone:     req.SubZone,
//		Description: req.Description,
//		Url:         req.Url,
//		CreateAt:    req.CreateAt,
//		UpdateAt:    req.UpdateAt,
//		Labels:      req.Labels,
//	}
//}
//
//func CoreFileToFile(req *core_api.File) *content.File {
//	return &content.File{
//		FileId:      req.FileId,
//		UserId:      req.UserId,
//		Name:        req.Name,
//		Type:        req.Type,
//		Path:        req.Path,
//		FatherId:    req.FatherId,
//		SpaceSize:   req.SpaceSize,
//		Md5:         req.Md5,
//		IsDel:       req.IsDel,
//		Zone:        req.Zone,
//		SubZone:     req.SubZone,
//		Description: req.Description,
//		Url:         req.Url,
//		Labels:      req.Labels,
//	}
//}

func ZoneToCoreZone(req *content.Zone) *core_api.Zone {
	return &core_api.Zone{
		Id:    req.Id,
		Value: req.Value,
	}
}

func CoreZoneToZone(req *core_api.Zone) *content.Zone {
	return &content.Zone{
		Id:    req.Id,
		Value: req.Value,
	}
}

func UserDetailToUser(req *content.User) *core_api.User {
	return &core_api.User{
		UserId: req.UserId,
		Name:   req.Name,
		Url:    req.Url,
	}
}

func UserDetailToCoreUserDetail(req *content.User) *core_api.UserDetail {
	return &core_api.UserDetail{
		Name:        req.Name,
		Sex:         req.Sex,
		FullName:    req.FullName,
		IdCard:      req.IdCard,
		Description: req.Description,
		Url:         req.Url,
	}
}

func PaginationOptionsToPaginationOptions(req *dto_basic.PaginationOptions) *basic.PaginationOptions {
	if req == nil {
		return &basic.PaginationOptions{}
	}
	return &basic.PaginationOptions{
		Limit:     req.Limit,
		LastToken: req.LastToken,
		Backward:  req.Backward,
		Offset:    req.Offset,
	}
}

func CoreApiRelationInfoToRelationInfo(req *core_api.RelationInfo) *relation.Relation {
	return &relation.Relation{
		FromType:     req.FromType,
		FromId:       req.FromId,
		ToType:       req.ToType,
		ToId:         req.ToId,
		RelationType: req.RelationType,
	}
}

func FilterOptionsToFilterOptions(opts *core_api.FileFilterOptions) *content.FileFilterOptions {
	if opts == nil {
		return nil
	} else {
		return &content.FileFilterOptions{
			OnlyUserId:       opts.OnlyUserId,
			OnlyFileId:       opts.OnlyFileId,
			OnlyFatherId:     opts.OnlyFatherId,
			OnlyIsDel:        opts.OnlyIsDel,
			OnlyDocumentType: opts.OnlyDocumentType,
			OnlyType:         opts.OnlyType,
		}
	}
}

func ShareOptionsToShareOptions(opts *core_api.ShareFileFilterOptions) *content.ShareFileFilterOptions {
	if opts == nil {
		return nil
	} else {
		return &content.ShareFileFilterOptions{
			OnlyCode:   opts.OnlyCode,
			OnlyUserId: opts.OnlyUserId,
		}
	}
}

func ShareCodeToCoreShareCode(opts *content.ShareCode) *core_api.ShareCode {
	return &core_api.ShareCode{
		Code:         opts.Code,
		Name:         opts.Name,
		Status:       opts.Status,
		BrowseNumber: opts.BrowseNumber,
		CreateAt:     opts.CreateAt,
	}
}

//func ShareFileToCoreShareFile(opts *content.ShareFile) *core_api.ShareFile {
//	return &core_api.ShareFile{
//		Code:          opts.Code,
//		UserId:        opts.UserId,
//		Name:          opts.Name,
//		Status:        opts.Status,
//		EffectiveTime: opts.EffectiveTime,
//		BrowseNumber:  opts.BrowseNumber,
//		CreateAt:      opts.CreateAt,
//		FileList:      opts.FileList,
//	}
//}

//func CoreShareFileToShareFile(opts *core_api.ShareFile) *content.ShareFile {
//	return &content.ShareFile{
//		Code:          opts.Code,
//		UserId:        opts.UserId,
//		Name:          opts.Name,
//		Status:        opts.Status,
//		EffectiveTime: opts.EffectiveTime,
//		BrowseNumber:  opts.BrowseNumber,
//		CreateAt:      opts.CreateAt,
//		FileList:      opts.FileList,
//	}
//}

func SearchOptionsToFileSearchOptions(opts *core_api.SearchOptions) *content.SearchOptions {
	if opts == nil {
		return nil
	} else {
		switch o := opts.Type.(type) {
		case *core_api.SearchOptions_AllFieldsKey:
			return &content.SearchOptions{Type: &content.SearchOptions_AllFieldsKey{AllFieldsKey: o.AllFieldsKey}}
		case *core_api.SearchOptions_MultiFieldsKey:
			return &content.SearchOptions{Type: &content.SearchOptions_MultiFieldsKey{MultiFieldsKey: &content.SearchField{
				Name:        o.MultiFieldsKey.Name,
				Id:          o.MultiFieldsKey.Id,
				Tag:         o.MultiFieldsKey.Tag,
				Text:        o.MultiFieldsKey.Text,
				Title:       o.MultiFieldsKey.Title,
				Description: o.MultiFieldsKey.Description,
				ProductName: o.MultiFieldsKey.ProductName,
			}}}
		}
	}
	return nil
}

func CoreUserInfoToUser(req *core_api.UserInfo) *content.User {
	return &content.User{
		UserId: req.UserId,
		Name:   req.Name,
		Sex:    req.Sex,
	}
}

func CoreUserDetailToUser(req *core_api.UserDetail) *content.User {
	return &content.User{
		Name:        req.Name,
		Sex:         req.Sex,
		FullName:    req.FullName,
		IdCard:      req.IdCard,
		Description: req.Description,
		Url:         req.Url,
	}
}
