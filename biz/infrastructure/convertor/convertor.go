package convertor

import (
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/basic"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/system"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/platform"
	"github.com/bytedance/sonic"
)

func FileToCorePublicFile(req *content.File) *core_api.PublicFile {
	return &core_api.PublicFile{
		Id:        req.Id,
		UserId:    req.UserId,
		Name:      req.Name,
		Type:      req.Type,
		SpaceSize: req.SpaceSize,
		IsDel:     req.IsDel,
		//Zone:         req.Zone,
		//SubZone:      req.SubZone,
		//Description:  req.Description,
		CreateAt:     req.CreateAt,
		UpdateAt:     req.UpdateAt,
		Author:       &core_api.FileUser{},
		FileCount:    &core_api.FileCount{},
		FileRelation: &core_api.FileRelation{},
	}
}

func FileToCorePrivateFile(req *content.File) *core_api.PrivateFile {
	return &core_api.PrivateFile{
		Id:        req.Id,
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

func CommentToCoreComment(req *platform.Comment) *core_api.Comment {
	return &core_api.Comment{
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
		Author:          &core_api.SimpleUser{},
		CommentRelation: &core_api.CommentRelation{},
	}
}

//func CoreCommentToComment(req *core_api.Comment) *platform.Comment {
//	return &platform.Comment{
//		Id:        req.Id,
//		SubjectId: req.SubjectId,
//		RootId:    req.RootId,
//		FatherId:  req.FatherId,
//		//Count:     req.Count,
//		State:    int64(req.State),
//		Attrs:    int64(req.Attrs),
//		Labels:   req.Labels,
//		UserId:   req.UserId,
//		AtUserId: req.AtUserId,
//		Content:  req.Content,
//		Meta:     req.Meta,
//	}
//}

func CoreLabelToLabel(req *core_api.Label) *platform.Label {
	return &platform.Label{
		Id:       "",
		FatherId: "",
		Value:    req.Value,
	}
}

func LabelToCoreLabel(req *platform.Label) *core_api.Label {
	return &core_api.Label{
		Id:    req.Id,
		Value: req.Value,
	}
}

//func SubjectDetailsToCoreSubjectDetails(req *platform.SubjectDetails) *core_api.SubjectDetails {
//	return &core_api.SubjectDetails{
//		Id:           req.Id,
//		UserId:       req.UserId,
//		TopCommentId: req.TopCommentId,
//		RootCount:    req.RootCount,
//		AllCount:     req.AllCount,
//		State:        req.State,
//		Attrs:        req.Attrs,
//	}
//}

func CoreSubjectToSubject(req *core_api.Subject) *platform.Subject {
	return &platform.Subject{
		Id:           req.Id,
		UserId:       req.UserId,
		TopCommentId: req.TopCommentId,
		RootCount:    0,
		AllCount:     0,
		//RootCount:    req.RootCount,
		//AllCount:     req.AllCount,
		State: int64(req.State),
		Attrs: int64(req.Attrs),
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
