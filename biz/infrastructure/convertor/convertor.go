package convertor

import (
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/basic"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/system"
	"github.com/bytedance/sonic"
)

func FileToCorePublicFile(req *content.FileInfo) *core_api.PublicFile {
	return &core_api.PublicFile{
		FileId:      req.FileId,
		UserId:      req.UserId,
		Name:        req.Name,
		Type:        req.Type,
		Path:        req.Path,
		FatherId:    req.FatherId,
		SpaceSize:   req.SpaceSize,
		IsDel:       req.IsDel,
		Zone:        req.Zone,
		SubZone:     req.SubZone,
		Description: req.Description,
		CreateAt:    req.CreateAt,
		UpdateAt:    req.UpdateAt,
		//Labels:       req.Labels,
		Author:       &core_api.User{},
		FileCount:    &core_api.PostCount{},
		FileRelation: &core_api.PostRelation{},
	}
}

func FileToCorePrivateFile(req *content.FileInfo) *core_api.PrivateFile {
	return &core_api.PrivateFile{
		FileId:    req.FileId,
		UserId:    req.UserId,
		Name:      req.Name,
		Type:      req.Type,
		Path:      req.Path,
		FatherId:  req.FatherId,
		SpaceSize: req.SpaceSize,
		IsDel:     req.IsDel,
		CreateAt:  req.CreateAt,
		UpdateAt:  req.UpdateAt,
	}
}

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

func MakePaginationOptions(limit, offset *int64, lastToken *string, backward *bool) *basic.PaginationOptions {
	return &basic.PaginationOptions{
		Limit:     limit,
		LastToken: lastToken,
		Backward:  backward,
		Offset:    offset,
	}
}

func ShareCodeToCoreShareCode(opts *content.ShareCode) *core_api.ShareCode {
	return &core_api.ShareCode{
		Code:         opts.Code,
		Name:         opts.Name,
		Status:       opts.Status,
		BrowseNumber: opts.BrowseNumber,
		CreateAt:     opts.CreateAt,
		Key:          opts.Key,
	}
}

func ShareFileToCoreShareFile(opts *content.ShareFile) *core_api.ShareFile {
	return &core_api.ShareFile{
		Code:          opts.Code,
		UserId:        opts.UserId,
		Name:          opts.Name,
		EffectiveTime: opts.EffectiveTime,
		BrowseNumber:  opts.BrowseNumber,
		CreateAt:      opts.CreateAt,
		FileList:      opts.FileList,
	}
}

func NotificationToCoreNotification(in *system.Notification) *core_api.Notification {
	msg := &Msg{}
	_ = sonic.UnmarshalString(in.Text, msg)
	return &core_api.Notification{
		FromName:   msg.FromName,
		FromId:     in.SourceUserId,
		ToName:     msg.ToName,
		ToId:       in.TargetUserId,
		ToType:     in.TargetType,
		Type:       in.Type,
		CreateTime: in.CreateTime,
	}
}

type Msg struct {
	FromName string
	ToName   string
}
