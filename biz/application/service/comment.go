package service

import (
	"context"
	"github.com/CloudStriver/cloudmind-core-api/biz/adaptor"
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/consts"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/convertor"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/platform_comment"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/platform/comment"
	"github.com/google/wire"
	"github.com/samber/lo"
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
	resp = new(core_api.GetCommentResp)
	var res *comment.GetCommentResp
	if res, err = s.PlatformComment.GetComment(ctx, &comment.GetCommentReq{CommentId: req.CommentId}); err != nil {
		return resp, err
	}
	resp.Comment = convertor.CommentInfoToCoreCommentInfo(res.Comment)
	return resp, nil
}

func (s *CommentService) GetComments(ctx context.Context, req *core_api.GetCommentsReq) (resp *core_api.GetCommentsResp, err error) {
	resp = new(core_api.GetCommentsResp)
	var res *comment.GetCommentListResp
	p := convertor.MakePaginationOptions(req.Limit, req.Offset, req.LastToken, req.Backward)
	if res, err = s.PlatformComment.GetCommentList(ctx, &comment.GetCommentListReq{FilterOptions: &comment.CommentFilterOptions{OnlyUserId: req.OnlyUserId, OnlyAtUserId: req.OnlyAtUserId, OnlyCommentId: req.OnlyCommentId, OnlySubjectId: req.OnlySubjectId, OnlyRootId: req.OnlyRootId, OnlyFatherId: req.OnlyFatherId, OnlyState: req.OnlyState, OnlyAttrs: req.OnlyAttrs}, Pagination: p}); err != nil {
		return resp, err
	}
	resp.Comments = lo.Map(res.Comments, func(item *comment.CommentInfo, _ int) *core_api.CommentInfo {
		return convertor.CommentInfoToCoreCommentInfo(item)
	})
	resp.Token = res.Token
	resp.Total = res.Total
	return resp, nil
}

func (s *CommentService) CreateComment(ctx context.Context, req *core_api.CreateCommentReq) (resp *core_api.CreateCommentResp, err error) {
	resp = new(core_api.CreateCommentResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	var res *comment.CreateCommentResp
	if res, err = s.PlatformComment.CreateComment(ctx, &comment.CreateCommentReq{Comment: convertor.CoreCommentToComment(req.Comment)}); err != nil {
		return resp, err
	}
	resp.CommentId = res.Id
	return resp, nil
}

func (s *CommentService) UpdateComment(ctx context.Context, req *core_api.UpdateCommentReq) (resp *core_api.UpdateCommentResp, err error) {
	resp = new(core_api.UpdateCommentResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	if _, err = s.PlatformComment.UpdateComment(ctx, &comment.UpdateCommentReq{Comment: convertor.CoreCommentToComment(req.Comment)}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *CommentService) DeleteComment(ctx context.Context, req *core_api.DeleteCommentReq) (resp *core_api.DeleteCommentResp, err error) {
	resp = new(core_api.DeleteCommentResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	if _, err = s.PlatformComment.DeleteCommentWithUserId(ctx, &comment.DeleteCommentWithUserIdReq{Id: req.CommentId, UserId: userData.UserId}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *CommentService) SetCommentState(ctx context.Context, req *core_api.SetCommentStateReq) (resp *core_api.SetCommentStateResp, err error) {
	resp = new(core_api.SetCommentStateResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	if _, err = s.PlatformComment.SetCommentState(ctx, &comment.SetCommentStateReq{Id: req.CommentId, UserId: userData.UserId, State: int64(req.State)}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *CommentService) SetCommentAttrs(ctx context.Context, req *core_api.SetCommentAttrsReq) (resp *core_api.SetCommentAttrsResp, err error) {
	resp = new(core_api.SetCommentAttrsResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	if _, err = s.PlatformComment.SetCommentAttrs(ctx, &comment.SetCommentAttrsReq{Id: req.CommentId, UserId: userData.UserId, Attrs: int64(req.Attrs), SubjectId: req.SubjectId, SortTime: req.SortTime}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *CommentService) GetCommentSubject(ctx context.Context, req *core_api.GetCommentSubjectReq) (resp *core_api.GetCommentSubjectResp, err error) {
	resp = new(core_api.GetCommentSubjectResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	var res *comment.GetCommentSubjectResp
	if res, err = s.PlatformComment.GetCommentSubject(ctx, &comment.GetCommentSubjectReq{Id: req.SubjectId}); err != nil {
		return resp, err
	}
	resp.Subject = convertor.SubjectDetailsToCoreSubjectDetails(res.Subject)
	return resp, nil
}

func (s *CommentService) CreateCommentSubject(ctx context.Context, req *core_api.CreateCommentSubjectReq) (resp *core_api.CreateCommentSubjectResp, err error) {
	resp = new(core_api.CreateCommentSubjectResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	var res *comment.CreateCommentSubjectResp
	subject := convertor.CoreSubjectToSubject(req.Subject)
	if res, err = s.PlatformComment.CreateCommentSubject(ctx, &comment.CreateCommentSubjectReq{Subject: subject}); err != nil {
		return resp, err
	}
	resp.SubjectId = res.Id
	return resp, nil
}

func (s *CommentService) UpdateCommentSubject(ctx context.Context, req *core_api.UpdateCommentSubjectReq) (resp *core_api.UpdateCommentSubjectResp, err error) {
	resp = new(core_api.UpdateCommentSubjectResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	subject := convertor.CoreSubjectToSubject(req.Subject)
	if _, err = s.PlatformComment.UpdateCommentSubject(ctx, &comment.UpdateCommentSubjectReq{Subject: subject}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *CommentService) DeleteCommentSubject(ctx context.Context, req *core_api.DeleteCommentSubjectReq) (resp *core_api.DeleteCommentSubjectResp, err error) {
	resp = new(core_api.DeleteCommentSubjectResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	if _, err = s.PlatformComment.DeleteCommentSubject(ctx, &comment.DeleteCommentSubjectReq{Id: req.Id, UserId: userData.UserId}); err != nil {
		return resp, err
	}
	return resp, nil
}
