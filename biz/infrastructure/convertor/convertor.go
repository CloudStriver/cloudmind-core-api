package convertor

import (
	dto_basic "github.com/CloudStriver/cloudmind-core-api/biz/application/dto/basic"
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/basic"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/platform/relation"
)

func UserToUserDetailInfo(req *core_api.UserDetail) *content.UserDetailInfo {
	return &content.UserDetailInfo{
		Name:        req.Name,
		Sex:         req.Sex,
		FullName:    req.FullName,
		IdCard:      req.IdCard,
		Description: req.Description,
		Url:         req.Url,
	}
}

func FileToCoreFile(req *content.FileInfo) *core_api.FileInfo {
	return &core_api.FileInfo{
		FileId:      req.FileId,
		UserId:      req.UserId,
		Name:        req.Name,
		Type:        core_api.Type(req.Type),
		Path:        req.Path,
		FatherId:    req.FatherId,
		SpaceSize:   req.SpaceSize,
		Md5:         req.Md5,
		IsDel:       req.IsDel,
		Tag:         req.Tag,
		Description: req.Description,
		UpdateAt:    req.UpdateAt,
	}
}

func CoreFileToFile(req *core_api.File) *content.File {
	return &content.File{
		FileId:      req.FileId,
		UserId:      req.UserId,
		Name:        req.Name,
		Type:        content.Type(req.Type),
		Path:        req.Path,
		FatherId:    req.FatherId,
		SpaceSize:   req.SpaceSize,
		Md5:         req.Md5,
		IsDel:       req.IsDel,
		Tag:         req.Tag,
		Description: req.Description,
		UpdateAt:    req.UpdateAt,
	}
}

func LabelToCoreLabel(req *content.Label) *core_api.Label {
	return &core_api.Label{
		Id:    req.Id,
		Value: req.Value,
	}
}

func CoreLabelToLabel(req *core_api.Label) *content.Label {
	return &content.Label{
		Id:    req.Id,
		Value: req.Value,
	}
}

func PostInfoToCorePostInfo(req *content.PostInfo) *core_api.PostInfo {
	return &core_api.PostInfo{
		PostId: req.PostId,
		UserId: req.UserId,
		Title:  req.Title,
		Text:   req.Text,
		Tags:   req.Tags,
		Status: req.Status,
		Url:    req.Url,
	}
}

func CorePostInfoToPostInfo(req *core_api.PostInfo) *content.PostInfo {
	return &content.PostInfo{
		PostId: req.PostId,
		UserId: req.UserId,
		Title:  req.Title,
		Text:   req.Text,
		Tags:   req.Tags,
		Status: req.Status,
		Url:    req.Url,
	}
}

func PostToCorePost(req *content.Post) *core_api.Post {
	return &core_api.Post{
		PostId: req.PostId,
		UserId: req.UserId,
		Title:  req.Title,
		Text:   req.Text,
		Tags:   req.Tags,
		Status: req.Status,
		Url:    req.Url,
	}
}

func CorePostToPost(req *core_api.PostInfo) *content.PostInfo {
	return &content.PostInfo{
		PostId: req.PostId,
		UserId: req.UserId,
		Title:  req.Title,
		Text:   req.Text,
		Tags:   req.Tags,
		Status: req.Status,
		Url:    req.Url,
	}
}

func UserDetailToUser(req *content.User) *core_api.User {
	return &core_api.User{
		UserId: req.UserId,
		Name:   req.Name,
		Url:    req.Url,
	}
}

func PaginationOptionsToPaginationOptions(req *dto_basic.PaginationOptions) *basic.PaginationOptions {
	return &basic.PaginationOptions{
		Limit:     req.Limit,
		LastToken: req.LastToken,
		Backward:  req.Backward,
		Offset:    req.Offset,
	}
}

func CoreApiRelationInfoToRelationInfo(req *core_api.RelationInfo) *relation.RelationInfo {
	return &relation.RelationInfo{
		FromType:     req.FromType,
		FromId:       req.FromId,
		ToType:       req.ToType,
		ToId:         req.ToId,
		RelationType: req.RelationType,
	}
}

func FilterOptionsToFilterOptions(opts *core_api.FileFilterOptions) (filter *content.FileFilterOptions) {
	if opts == nil {
		filter = &content.FileFilterOptions{}
	} else {
		filter = &content.FileFilterOptions{
			OnlyUserId:      opts.OnlyUserId,
			OnlyFileId:      opts.OnlyFileId,
			OnlyFatherId:    opts.OnlyFatherId,
			OnlyFileType:    opts.OnlyFileType,
			OnlyTags:        opts.OnlyTags,
			IsDel:           opts.IsDel,
			DocumentType:    opts.DocumentType,
			OnlyMd5:         opts.OnlyMd5,
			OnlySetRelation: opts.OnlySetRelation,
		}
	}
	return filter
}

func ShareOptionsToShareOptions(opts *core_api.ShareFileFilterOptions) (filter *content.ShareFileFilterOptions) {
	if opts == nil {
		filter = &content.ShareFileFilterOptions{}
	} else {
		filter = &content.ShareFileFilterOptions{
			OnlyCode:   opts.OnlyCode,
			OnlyUserId: opts.OnlyUserId,
		}
	}
	return filter
}

func ShareCodeToShareCode(opts *content.ShareCode) *core_api.ShareCode {
	return &core_api.ShareCode{
		Code:         opts.Code,
		Name:         opts.Name,
		Status:       opts.Status,
		BrowseNumber: opts.BrowseNumber,
		CreateAt:     opts.CreateAt,
	}
}

func ShareFileToCoreShareFile(opts *content.ShareFile) *core_api.ShareFile {
	return &core_api.ShareFile{
		Code:          opts.Code,
		UserId:        opts.UserId,
		Name:          opts.Name,
		Status:        opts.Status,
		EffectiveTime: opts.EffectiveTime,
		BrowseNumber:  opts.BrowseNumber,
		CreateAt:      opts.CreateAt,
		FileList:      opts.FileList,
	}
}

func CoreShareFileToShareFile(opts *core_api.ShareFile) *content.ShareFile {
	return &content.ShareFile{
		Code:          opts.Code,
		UserId:        opts.UserId,
		Name:          opts.Name,
		Status:        opts.Status,
		EffectiveTime: opts.EffectiveTime,
		BrowseNumber:  opts.BrowseNumber,
		CreateAt:      opts.CreateAt,
		FileList:      opts.FileList,
	}
}

func SearchOptionsToFileSearchOptions(opts *core_api.SearchOptions) (filter *content.SearchOptions) {
	if opts == nil {
		filter = &content.SearchOptions{}
	} else {
		switch o := opts.Type.(type) {
		case *core_api.SearchOptions_AllFieldsKey:
			filter = &content.SearchOptions{Type: &content.SearchOptions_AllFieldsKey{AllFieldsKey: o.AllFieldsKey}}
		case *core_api.SearchOptions_MultiFieldsKey:
			filter = &content.SearchOptions{Type: &content.SearchOptions_MultiFieldsKey{MultiFieldsKey: &content.SearchField{
				Name:  o.MultiFieldsKey.Name,
				Id:    o.MultiFieldsKey.Id,
				Tag:   o.MultiFieldsKey.Tag,
				Text:  o.MultiFieldsKey.Text,
				Title: o.MultiFieldsKey.Title,
			}}}
		}
	}
	return filter
}
