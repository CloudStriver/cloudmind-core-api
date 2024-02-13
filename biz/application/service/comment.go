package service

import (
	"context"
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/platform_comment"
	"github.com/google/wire"
)

type ICommentService interface {
	GetComment(ctx context.Context, req *core_api.GetCommentReq) (resp *core_api.GetCommentResp, err error)
	GetComments(ctx context.Context, req *core_api.GetCommentsReq) (resp *core_api.GetCommentsResp, err error)
	CreateComment(ctx context.Context, req *core_api.CreateCommentReq) (resp *core_api.CreateCommentResp, err error)
	UpdateComment(ctx context.Context, req *core_api.UpdateCommentReq) (resp *core_api.UpdateCommentResp, err error)
	DeleteComment(ctx context.Context, req *core_api.DeleteCommentReq) (resp *core_api.DeleteCommentResp, err error)
	SetCommentState(ctx context.Context, req *core_api.SetCommentStateReq) (resp *core_api.SetCommentStateResp, err error)
	SetCommentAttrs(ctx context.Context, req *core_api.SetCommentAttrsReq) (resp *core_api.SetCommentAttrsResp, err error)
	GetCommentSubject(ctx context.Context, req *core_api.GetCommentSubjectReq) (resp *core_api.GetCommentSubjectResp, err error)
	CreateCommentSubject(ctx context.Context, req *core_api.CreateCommentSubjectReq) (resp *core_api.CreateCommentSubjectResp, err error)
	UpdateCommentSubject(ctx context.Context, req *core_api.UpdateCommentSubjectReq) (resp *core_api.UpdateCommentSubjectResp, err error)
	DeleteCommentSubject(ctx context.Context, req *core_api.DeleteCommentSubjectReq) (resp *core_api.DeleteCommentSubjectResp, err error)
	//SetCommentSubjectState(ctx context.Context, req *core_api.SetCommentSubjectStateReq) (resp *core_api.SetCommentSubjectStateResp, err error)
	//SetCommentSubjectAttrs(ctx context.Context, req *core_api.SetCommentSubjectAttrsReq) (resp *core_api.SetCommentSubjectAttrsResp, err error)
}

var CommentServiceSet = wire.NewSet(
	wire.Struct(new(CommentService), "*"),
	wire.Bind(new(ICommentService), new(*CommentService)),
)

type CommentService struct {
	Config          *config.Config
	PlatformComment platform_comment.IPlatFormComment
}

func (s *CommentService) GetComment(ctx context.Context, req *core_api.GetCommentReq) (resp *core_api.GetCommentResp, err error) {
	return
}

func (s *CommentService) GetComments(ctx context.Context, req *core_api.GetCommentsReq) (resp *core_api.GetCommentsResp, err error) {
	return
}

func (s *CommentService) CreateComment(ctx context.Context, req *core_api.CreateCommentReq) (resp *core_api.CreateCommentResp, err error) {
	return
}

func (s *CommentService) UpdateComment(ctx context.Context, req *core_api.UpdateCommentReq) (resp *core_api.UpdateCommentResp, err error) {
	return
}

func (s *CommentService) DeleteComment(ctx context.Context, req *core_api.DeleteCommentReq) (resp *core_api.DeleteCommentResp, err error) {
	return
}

func (s *CommentService) SetCommentState(ctx context.Context, req *core_api.SetCommentStateReq) (resp *core_api.SetCommentStateResp, err error) {
	return
}

func (s *CommentService) SetCommentAttrs(ctx context.Context, req *core_api.SetCommentAttrsReq) (resp *core_api.SetCommentAttrsResp, err error) {
	return
}

func (s *CommentService) GetCommentSubject(ctx context.Context, req *core_api.GetCommentSubjectReq) (resp *core_api.GetCommentSubjectResp, err error) {
	return
}

func (s *CommentService) CreateCommentSubject(ctx context.Context, req *core_api.CreateCommentSubjectReq) (resp *core_api.CreateCommentSubjectResp, err error) {
	return
}

func (s *CommentService) UpdateCommentSubject(ctx context.Context, req *core_api.UpdateCommentSubjectReq) (resp *core_api.UpdateCommentSubjectResp, err error) {
	return
}

func (s *CommentService) DeleteCommentSubject(ctx context.Context, req *core_api.DeleteCommentSubjectReq) (resp *core_api.DeleteCommentSubjectResp, err error) {
	return
}

//func (s *CommentService) SetCommentSubjectState(ctx context.Context, req *core_api.SetCommentSubjectStateReq) (resp *core_api.SetCommentSubjectStateResp, err error) {
//	return
//}
//
//func (s *CommentService) SetCommentSubjectAttrs(ctx context.Context, req *core_api.SetCommentSubjectAttrsReq) (resp *core_api.SetCommentSubjectAttrsResp, err error) {
//	return
//}
