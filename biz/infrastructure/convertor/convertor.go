package convertor

import (
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/basic"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/system"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/platform"
	"github.com/bytedance/sonic"
)

func FileToCorePublicFile(req *content.PublicFile) *core_api.PublicFile {
	return &core_api.PublicFile{
		FileId:       req.FileId,
		UserId:       req.UserId,
		Name:         req.Name,
		Type:         req.Type,
		Path:         req.Path,
		SpaceSize:    req.SpaceSize,
		Zone:         req.Zone,
		Description:  req.Description,
		AuditStatus:  req.AuditStatus,
		CreateAt:     req.CreateAt,
		Author:       &core_api.FileUser{},
		FileCount:    &core_api.FileCount{},
		FileRelation: &core_api.FileRelation{},
	}
}

func FileToCorePrivateFile(req *content.File) *core_api.PrivateFile {
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
		DeleteAt:  req.DeleteAt,
	}
}

func CommentToCoreComment(req *platform.Comment) *core_api.Comment {
	return &core_api.Comment{
		CommentId:       req.CommentId,
		SubjectId:       req.SubjectId,
		RootId:          req.RootId,
		FatherId:        req.FatherId,
		Count:           req.Count,
		State:           req.State,
		Attrs:           req.Attrs,
		AtUserId:        req.AtUserId,
		Content:         req.Content,
		Meta:            req.Meta,
		CreateTime:      req.CreateTime,
		Author:          &core_api.SimpleUser{},
		CommentRelation: &core_api.CommentRelation{},
		ItemType:        core_api.TargetType(req.Type),
		AtUser:          &core_api.SimpleUser{},
	}
}

func CommentToCoreCommentNode(req *platform.Comment) *core_api.CommentNode {
	return &core_api.CommentNode{
		CommentId:       req.CommentId,
		SubjectId:       req.SubjectId,
		RootId:          req.RootId,
		FatherId:        req.FatherId,
		Count:           req.Count,
		State:           req.State,
		Attrs:           req.Attrs,
		AtUserId:        req.AtUserId,
		Content:         req.Content,
		Meta:            req.Meta,
		CreateTime:      req.CreateTime,
		Author:          &core_api.SimpleUser{},
		CommentRelation: &core_api.CommentRelation{},
		AtUser:          &core_api.SimpleUser{},
	}
}

func LabelToCoreLabel(req *platform.Label) *core_api.Label {
	return &core_api.Label{
		Id:    req.LabelId,
		Value: req.Value,
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

func NotificationToCoreNotification(in *system.Notification) *core_api.Notification {
	msg := &Msg{}
	_ = sonic.UnmarshalString(in.Text, msg)
	return &core_api.Notification{
		NotificationId: in.NotificationId,
		FromName:       msg.FromName,
		FromId:         in.SourceUserId,
		ToName:         msg.ToName,
		ToId:           in.SourceContentId,
		ToUserId:       in.TargetUserId,
		Type:           in.Type,
		CreateTime:     in.CreateTime,
	}
}

type Msg struct {
	FromName string
	ToName   string
}

func SystemSliderToSlider(in *system.Slider) *core_api.Slider {
	return &core_api.Slider{
		SliderId: in.SliderId,
		ImageUrl: in.ImageUrl,
		LinkUrl:  in.LinkUrl,
		IsPublic: in.IsPublic,
	}
}
