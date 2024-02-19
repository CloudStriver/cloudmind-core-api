package service

import (
	"context"
	"github.com/CloudStriver/cloudmind-core-api/biz/adaptor"
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/cloudmind-core-api/biz/domain/service"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/consts"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/convertor"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/platform_comment"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/platform/comment"
	"github.com/google/wire"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/mr"
)

type ICommentService interface {
	GetComment(ctx context.Context, req *core_api.GetCommentReq) (resp *core_api.GetCommentResp, err error)
	GetComments(ctx context.Context, req *core_api.GetCommentsReq) (resp *core_api.GetCommentsResp, err error)
	CreateComment(ctx context.Context, req *core_api.CreateCommentReq) (resp *core_api.CreateCommentResp, err error)
	UpdateComment(ctx context.Context, req *core_api.UpdateCommentReq) (resp *core_api.UpdateCommentResp, err error)
	DeleteComment(ctx context.Context, req *core_api.DeleteCommentReq) (resp *core_api.DeleteCommentResp, err error)
	SetCommentAttrs(ctx context.Context, req *core_api.SetCommentAttrsReq) (resp *core_api.SetCommentAttrsResp, err error)
	GetCommentSubject(ctx context.Context, req *core_api.GetCommentSubjectReq) (resp *core_api.GetCommentSubjectResp, err error)
	UpdateCommentSubject(ctx context.Context, req *core_api.UpdateCommentSubjectReq) (resp *core_api.UpdateCommentSubjectResp, err error)
	DeleteCommentSubject(ctx context.Context, req *core_api.DeleteCommentSubjectReq) (resp *core_api.DeleteCommentSubjectResp, err error)
}

var CommentServiceSet = wire.NewSet(
	wire.Struct(new(CommentService), "*"),
	wire.Bind(new(ICommentService), new(*CommentService)),
)

type CommentService struct {
	Config               *config.Config
	PlatformComment      platform_comment.IPlatFormComment
	CommentDomainService service.ICommentDomainService
}

func (s *CommentService) GetComment(ctx context.Context, req *core_api.GetCommentReq) (resp *core_api.GetCommentResp, err error) {
	resp = new(core_api.GetCommentResp)
	userData := adaptor.ExtractUserMeta(ctx)
	var res *comment.GetCommentResp
	if res, err = s.PlatformComment.GetComment(ctx, &comment.GetCommentReq{CommentId: req.CommentId}); err != nil {
		return resp, err
	}
	resp.Comment = convertor.CommentInfoToCoreCommentInfo(res.Comment)
	_ = mr.Finish(func() error {
		s.CommentDomainService.LoadLikeCount(ctx, resp.Comment) // 点赞量
		return nil
	}, func() error {
		s.CommentDomainService.LoadAuthor(ctx, resp.Comment, userData.GetUserId()) // 作者
		return nil
	}, func() error {
		s.CommentDomainService.LoadLiked(ctx, resp.Comment, userData.GetUserId()) // 是否点赞
		return nil
	}, func() error {
		s.CommentDomainService.LoadHated(ctx, resp.Comment, userData.GetUserId()) // 是否点踩
		return nil
	}, func() error {
		s.CommentDomainService.LoadLabels(ctx, resp.Comment, res.Comment.Labels) // 标签集
		return nil
	})
	return resp, nil
}

func (s *CommentService) GetComments(ctx context.Context, req *core_api.GetCommentsReq) (resp *core_api.GetCommentsResp, err error) {
	resp = new(core_api.GetCommentsResp)
	userData := adaptor.ExtractUserMeta(ctx)
	var res *comment.GetCommentListResp
	p := convertor.MakePaginationOptions(req.Limit, req.Offset, req.LastToken, req.Backward)
	if res, err = s.PlatformComment.GetCommentList(ctx, &comment.GetCommentListReq{FilterOptions: &comment.CommentFilterOptions{OnlyUserId: req.OnlyUserId, OnlyAtUserId: req.OnlyAtUserId, OnlyCommentId: req.OnlyCommentId, OnlySubjectId: req.OnlySubjectId, OnlyRootId: req.OnlyRootId, OnlyFatherId: req.OnlyFatherId, OnlyState: req.OnlyState, OnlyAttrs: req.OnlyAttrs}, Pagination: p}); err != nil {
		return resp, err
	}
	resp.Comments = lo.Map(res.Comments, func(item *comment.CommentInfo, _ int) *core_api.CommentInfo {
		c := convertor.CommentInfoToCoreCommentInfo(item)
		_ = mr.Finish(func() error {
			s.CommentDomainService.LoadLikeCount(ctx, c) // 点赞量
			return nil
		}, func() error {
			s.CommentDomainService.LoadAuthor(ctx, c, userData.GetUserId()) // 作者
			return nil
		}, func() error {
			s.CommentDomainService.LoadLiked(ctx, c, userData.GetUserId()) // 是否点赞
			return nil
		}, func() error {
			s.CommentDomainService.LoadHated(ctx, c, userData.GetUserId()) // 是否点踩
			return nil
		}, func() error {
			s.CommentDomainService.LoadLabels(ctx, c, item.Labels) // 标签集
			return nil
		})
		return c
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
	req.Comment.UserId = userData.UserId
	req.Comment.Count = lo.ToPtr(consts.InitNumber)
	req.Comment.State = core_api.State_Normal
	req.Comment.Attrs = core_api.Attrs_None
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

	var res *comment.GetCommentResp
	if res, err = s.PlatformComment.GetComment(ctx, &comment.GetCommentReq{CommentId: req.Comment.Id}); err != nil {
		return resp, err
	}
	if res.Comment.UserId == userData.UserId {
		if _, err = s.PlatformComment.UpdateComment(ctx, &comment.UpdateCommentReq{Comment: convertor.CoreCommentToComment(req.Comment)}); err != nil {
			return resp, err
		}
	}
	return resp, nil
}

func (s *CommentService) DeleteComment(ctx context.Context, req *core_api.DeleteCommentReq) (resp *core_api.DeleteCommentResp, err error) {
	resp = new(core_api.DeleteCommentResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	var ok bool
	var res *comment.GetCommentResp
	if res, err = s.PlatformComment.GetComment(ctx, &comment.GetCommentReq{CommentId: req.CommentId}); err != nil {
		return resp, err
	}
	switch {
	case res.Comment.UserId == userData.UserId:
		ok = true
	default:
		var result *comment.GetCommentSubjectResp
		if result, err = s.PlatformComment.GetCommentSubject(ctx, &comment.GetCommentSubjectReq{Id: res.Comment.SubjectId}); err != nil {
			return resp, err
		}
		ok = result.Subject.UserId == userData.UserId
	}
	if ok {
		if _, err = s.PlatformComment.DeleteComment(ctx, &comment.DeleteCommentReq{Id: req.CommentId}); err != nil {
			return resp, err
		}
	}
	return resp, nil
}

func (s *CommentService) SetCommentAttrs(ctx context.Context, req *core_api.SetCommentAttrsReq) (resp *core_api.SetCommentAttrsResp, err error) {
	resp = new(core_api.SetCommentAttrsResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	var res *comment.GetCommentResp
	if res, err = s.PlatformComment.GetComment(ctx, &comment.GetCommentReq{CommentId: req.CommentId}); err != nil {
		return resp, err
	}
	if res.Comment.UserId == userData.UserId {
		if _, err = s.PlatformComment.SetCommentAttrs(ctx, &comment.SetCommentAttrsReq{Id: req.CommentId, Attrs: int64(req.Attrs), SubjectId: res.Comment.SubjectId, SortTime: res.Comment.CreateTime}); err != nil {
			return resp, err
		}
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

func (s *CommentService) UpdateCommentSubject(ctx context.Context, req *core_api.UpdateCommentSubjectReq) (resp *core_api.UpdateCommentSubjectResp, err error) {
	resp = new(core_api.UpdateCommentSubjectResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	var res *comment.GetCommentSubjectResp
	if res, err = s.PlatformComment.GetCommentSubject(ctx, &comment.GetCommentSubjectReq{Id: req.Subject.Id}); err != nil {
		return resp, err
	}
	if res.Subject.UserId == userData.UserId {
		subject := convertor.CoreSubjectToSubject(req.Subject)
		if _, err = s.PlatformComment.UpdateCommentSubject(ctx, &comment.UpdateCommentSubjectReq{Subject: subject}); err != nil {
			return resp, err
		}
	}
	return resp, nil
}

func (s *CommentService) DeleteCommentSubject(ctx context.Context, req *core_api.DeleteCommentSubjectReq) (resp *core_api.DeleteCommentSubjectResp, err error) {
	resp = new(core_api.DeleteCommentSubjectResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	var res *comment.GetCommentSubjectResp
	if res, err = s.PlatformComment.GetCommentSubject(ctx, &comment.GetCommentSubjectReq{Id: req.Id}); err != nil {
		return resp, err
	}
	if res.Subject.UserId == userData.UserId {
		if _, err = s.PlatformComment.DeleteCommentSubject(ctx, &comment.DeleteCommentSubjectReq{Id: req.Id}); err != nil {
			return resp, err
		}
	}
	return resp, nil
}
