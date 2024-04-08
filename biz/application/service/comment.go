package service

import (
	"context"
	"github.com/CloudStriver/cloudmind-core-api/biz/adaptor"
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/cloudmind-core-api/biz/domain/service"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/consts"
	platformservice "github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/platform"
	"github.com/google/wire"
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
	Platform             platformservice.IPlatForm
	CommentDomainService service.ICommentDomainService
}

func (s *CommentService) GetComment(ctx context.Context, req *core_api.GetCommentReq) (resp *core_api.GetCommentResp, err error) {
	resp = new(core_api.GetCommentResp)
	//var res *platform.GetCommentResp
	//if res, err = s.Platform.GetComment(ctx, &platform.GetCommentReq{CommentId: req.CommentId}); err != nil {
	//	return resp, err
	//}
	//resp.Comment = convertor.CommentInfoToCoreCommentInfo(res.Comment)
	//_ = mr.Finish(func() error {
	//	s.CommentDomainService.LoadLikeCount(ctx, resp.Comment) // 点赞量
	//	return nil
	//}, func() error {
	//	s.CommentDomainService.LoadAuthor(ctx, resp.Comment, resp.Comment.UserId) // 作者
	//	return nil
	//}, func() error {
	//	s.CommentDomainService.LoadLiked(ctx, resp.Comment, resp.Comment.UserId) // 是否点赞
	//	return nil
	//}, func() error {
	//	s.CommentDomainService.LoadHated(ctx, resp.Comment, resp.Comment.UserId) // 是否点踩
	//	return nil
	//}, func() error {
	//	s.CommentDomainService.LoadLabels(ctx, resp.Comment, res.Comment.Labels) // 标签集
	//	return nil
	//})
	return resp, nil
}

func (s *CommentService) GetComments(ctx context.Context, req *core_api.GetCommentsReq) (resp *core_api.GetCommentsResp, err error) {
	resp = new(core_api.GetCommentsResp)
	//var res *platform.GetCommentListResp
	//p := convertor.MakePaginationOptions(req.Limit, req.Offset, req.LastToken, req.Backward)
	//if res, err = s.Platform.GetCommentList(ctx, &platform.GetCommentListReq{FilterOptions: &platform.CommentFilterOptions{OnlyUserId: req.OnlyUserId, OnlyAtUserId: req.OnlyAtUserId, OnlyCommentId: req.OnlyCommentId, OnlySubjectId: req.OnlySubjectId, OnlyRootId: req.OnlyRootId, OnlyFatherId: req.OnlyFatherId, OnlyState: req.OnlyState, OnlyAttrs: req.OnlyAttrs}, Pagination: p}); err != nil {
	//	return resp, err
	//}
	//resp.Comments = lo.Map(res.Comments, func(item *platform.CommentInfo, _ int) *core_api.CommentInfo {
	//	c := convertor.CommentInfoToCoreCommentInfo(item)
	//	_ = mr.Finish(func() error {
	//		s.CommentDomainService.LoadLikeCount(ctx, c) // 点赞量
	//		return nil
	//	}, func() error {
	//		s.CommentDomainService.LoadAuthor(ctx, c, c.UserId) // 作者
	//		return nil
	//	}, func() error {
	//		s.CommentDomainService.LoadLiked(ctx, c, c.UserId) // 是否点赞
	//		return nil
	//	}, func() error {
	//		s.CommentDomainService.LoadHated(ctx, c, c.UserId) // 是否点踩
	//		return nil
	//	}, func() error {
	//		s.CommentDomainService.LoadLabels(ctx, c, item.Labels) // 标签集
	//		return nil
	//	})
	//	return c
	//})
	//resp.Token = res.Token
	//resp.Total = res.Total
	return resp, nil
}

func (s *CommentService) CreateComment(ctx context.Context, req *core_api.CreateCommentReq) (resp *core_api.CreateCommentResp, err error) {
	resp = new(core_api.CreateCommentResp)
	//userData, err := adaptor.ExtractUserMeta(ctx)
	//if err != nil || userData.GetUserId() == "" {
	//	return resp, consts.ErrNotAuthentication
	//}
	//var res *platform.CreateCommentResp
	//req.Comment.UserId = userData.UserId
	//req.Comment.Count = lo.ToPtr(consts.InitNumber)
	//req.Comment.State = core_api.State_Normal
	//req.Comment.Attrs = core_api.Attrs_None
	//if res, err = s.Platform.CreateComment(ctx, &platform.CreateCommentReq{Comment: convertor.CoreCommentToComment(req.Comment)}); err != nil {
	//	return resp, err
	//}
	//resp.CommentId = res.Id
	return resp, nil
}

func (s *CommentService) UpdateComment(ctx context.Context, req *core_api.UpdateCommentReq) (resp *core_api.UpdateCommentResp, err error) {
	resp = new(core_api.UpdateCommentResp)
	//userData, err := adaptor.ExtractUserMeta(ctx)
	//if err != nil || userData.GetUserId() == "" {
	//	return resp, consts.ErrNotAuthentication
	//}
	//
	//var res *platform.GetCommentResp
	//if res, err = s.Platform.GetComment(ctx, &platform.GetCommentReq{CommentId: req.Comment.Id}); err != nil {
	//	return resp, err
	//}
	//if res.Comment.UserId == userData.UserId {
	//	if _, err = s.Platform.UpdateComment(ctx, &platform.UpdateCommentReq{Comment: convertor.CoreCommentToComment(req.Comment)}); err != nil {
	//		return resp, err
	//	}
	//}
	return resp, nil
}

func (s *CommentService) DeleteComment(ctx context.Context, req *core_api.DeleteCommentReq) (resp *core_api.DeleteCommentResp, err error) {
	resp = new(core_api.DeleteCommentResp)
	//userData, err := adaptor.ExtractUserMeta(ctx)
	//if err != nil || userData.GetUserId() == "" {
	//	return resp, consts.ErrNotAuthentication
	//}
	//var ok bool
	//var res *platform.GetCommentResp
	//if res, err = s.Platform.GetComment(ctx, &platform.GetCommentReq{CommentId: req.CommentId}); err != nil {
	//	return resp, err
	//}
	//switch {
	//case res.Comment.UserId == userData.UserId:
	//	ok = true
	//default:
	//	var result *platform.GetCommentSubjectResp
	//	if result, err = s.Platform.GetCommentSubject(ctx, &platform.GetCommentSubjectReq{Id: res.Comment.SubjectId}); err != nil {
	//		return resp, err
	//	}
	//	ok = result.Subject.UserId == userData.UserId
	//}
	//if ok {
	//	if _, err = s.Platform.DeleteComment(ctx, &platform.DeleteCommentReq{Id: req.CommentId}); err != nil {
	//		return resp, err
	//	}
	//}
	return resp, nil
}

func (s *CommentService) SetCommentAttrs(ctx context.Context, req *core_api.SetCommentAttrsReq) (resp *core_api.SetCommentAttrsResp, err error) {
	resp = new(core_api.SetCommentAttrsResp)
	//userData, err := adaptor.ExtractUserMeta(ctx)
	//if err != nil || userData.GetUserId() == "" {
	//	return resp, consts.ErrNotAuthentication
	//}
	//var res *platform.GetCommentResp
	//if res, err = s.Platform.GetComment(ctx, &platform.GetCommentReq{CommentId: req.CommentId}); err != nil {
	//	return resp, err
	//}
	//if res.Comment.UserId == userData.UserId {
	//	if _, err = s.Platform.SetCommentAttrs(ctx, &platform.SetCommentAttrsReq{Id: req.CommentId, Attrs: int64(req.Attrs), SubjectId: res.Comment.SubjectId, SortTime: res.Comment.CreateTime}); err != nil {
	//		return resp, err
	//	}
	//}
	return resp, nil
}

func (s *CommentService) GetCommentSubject(ctx context.Context, req *core_api.GetCommentSubjectReq) (resp *core_api.GetCommentSubjectResp, err error) {
	resp = new(core_api.GetCommentSubjectResp)
	//userData, err := adaptor.ExtractUserMeta(ctx)
	//if err != nil || userData.GetUserId() == "" {
	//	return resp, consts.ErrNotAuthentication
	//}
	//var res *platform.GetCommentSubjectResp
	//if res, err = s.Platform.GetCommentSubject(ctx, &platform.GetCommentSubjectReq{Id: req.SubjectId}); err != nil {
	//	return resp, err
	//}
	//resp.Subject = convertor.SubjectDetailsToCoreSubjectDetails(res.Subject)
	return resp, nil
}

func (s *CommentService) UpdateCommentSubject(ctx context.Context, req *core_api.UpdateCommentSubjectReq) (resp *core_api.UpdateCommentSubjectResp, err error) {
	resp = new(core_api.UpdateCommentSubjectResp)
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	//var res *platform.GetCommentSubjectResp
	//if res, err = s.Platform.GetCommentSubject(ctx, &platform.GetCommentSubjectReq{Id: req.Subject.Id}); err != nil {
	//	return resp, err
	//}
	//if res.Subject.UserId == userData.UserId {
	//	subject := convertor.CoreSubjectToSubject(req.Subject)
	//	if _, err = s.Platform.UpdateCommentSubject(ctx, &platform.UpdateCommentSubjectReq{Subject: subject}); err != nil {
	//		return resp, err
	//	}
	//}
	return resp, nil
}

func (s *CommentService) DeleteCommentSubject(ctx context.Context, req *core_api.DeleteCommentSubjectReq) (resp *core_api.DeleteCommentSubjectResp, err error) {
	resp = new(core_api.DeleteCommentSubjectResp)
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	//var res *platform.GetCommentSubjectResp
	//if res, err = s.Platform.GetCommentSubject(ctx, &platform.GetCommentSubjectReq{Id: req.Id}); err != nil {
	//	return resp, err
	//}
	//if res.Subject.UserId == userData.UserId {
	//	if _, err = s.Platform.DeleteCommentSubject(ctx, &platform.DeleteCommentSubjectReq{Id: req.Id}); err != nil {
	//		return resp, err
	//	}
	//}
	return resp, nil
}
