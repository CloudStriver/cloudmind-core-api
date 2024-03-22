package convertor

import (
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/basic"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/system"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/platform/comment"
	"github.com/bytedance/sonic"
)

func FileToCorePublicFile(req *content.FileInfo) *core_api.PublicFile {
	return &core_api.PublicFile{
		FileId:       req.FileId,
		UserId:       req.UserId,
		Name:         req.Name,
		Type:         req.Type,
		Path:         req.Path,
		FatherId:     req.FatherId,
		SpaceSize:    req.SpaceSize,
		IsDel:        req.IsDel,
		Zone:         req.Zone,
		SubZone:      req.SubZone,
		Description:  req.Description,
		CreateAt:     req.CreateAt,
		UpdateAt:     req.UpdateAt,
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

func CommentInfoToCoreCommentInfo(req *comment.CommentInfo) *core_api.CommentInfo {
	return &core_api.CommentInfo{
		Id:              req.Id,
		SubjectId:       req.SubjectId,
		RootId:          req.RootId,
		FatherId:        req.FatherId,
		Count:           req.Count,
		State:           req.State,
		Attrs:           req.Attrs,
		UserId:          req.UserId,
		AtUserId:        req.AtUserId,
		Content:         req.Content,
		Meta:            req.Meta,
		CreateTime:      req.CreateTime,
		Like:            0,
		Author:          &core_api.User{},
		CommentRelation: &core_api.CommentRelation{},
	}
}

func CoreCommentToComment(req *core_api.Comment) *comment.Comment {
	return &comment.Comment{
		Id:        req.Id,
		SubjectId: req.SubjectId,
		RootId:    req.RootId,
		FatherId:  req.FatherId,
		Count:     req.Count,
		State:     int64(req.State),
		Attrs:     int64(req.Attrs),
		Labels:    req.Labels,
		UserId:    req.UserId,
		AtUserId:  req.AtUserId,
		Content:   req.Content,
		Meta:      req.Meta,
	}
}

func CoreLabelToLabel(req *core_api.Label) *comment.Label {
	return &comment.Label{
		LabelId: req.LabelId,
		Value:   req.Value,
	}
}

func LabelToCoreLabel(req *comment.Label) *core_api.Label {
	return &core_api.Label{
		LabelId: req.LabelId,
		Value:   req.Value,
	}
}

func SubjectDetailsToCoreSubjectDetails(req *comment.SubjectDetails) *core_api.SubjectDetails {
	return &core_api.SubjectDetails{
		Id:           req.Id,
		UserId:       req.UserId,
		TopCommentId: req.TopCommentId,
		RootCount:    req.RootCount,
		AllCount:     req.AllCount,
		State:        req.State,
		Attrs:        req.Attrs,
	}
}

func CoreSubjectToSubject(req *core_api.Subject) *comment.Subject {
	return &comment.Subject{
		Id:           req.Id,
		UserId:       req.UserId,
		TopCommentId: req.TopCommentId,
		RootCount:    req.RootCount,
		AllCount:     req.AllCount,
		State:        int64(req.State),
		Attrs:        int64(req.Attrs),
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
		NotificationId: in.NotificationId,
		FromName:       msg.FromName,
		FromId:         in.SourceUserId,
		ToName:         msg.ToName,
		ToId:           in.TargetUserId,
		ToType:         in.TargetType,
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
