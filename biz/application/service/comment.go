package service

import (
	"context"
	"github.com/CloudStriver/cloudmind-core-api/biz/adaptor"
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/cloudmind-core-api/biz/domain/service"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/consts"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/convertor"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_content"
	platformservice "github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/platform"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/platform"
	"github.com/google/wire"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/mr"
)

type ICommentService interface {
	GetComment(ctx context.Context, req *core_api.GetCommentReq) (resp *core_api.GetCommentResp, err error)
	GetComments(ctx context.Context, req *core_api.GetCommentsReq) (resp *core_api.GetCommentsResp, err error)
	GetCommentBlocks(ctx context.Context, req *core_api.GetCommentBlocksReq) (resp *core_api.GetCommentBlocksResp, err error)
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
	Config                *config.Config
	Platform              platformservice.IPlatForm
	CloudMindContent      cloudmind_content.ICloudMindContent
	CommentDomainService  service.ICommentDomainService
	RelationDomainService service.RelationDomainService
}

func (s *CommentService) GetComment(ctx context.Context, req *core_api.GetCommentReq) (resp *core_api.GetCommentResp, err error) {
	resp = new(core_api.GetCommentResp)
	var res *platform.GetCommentResp
	if res, err = s.Platform.GetComment(ctx, &platform.GetCommentReq{CommentId: req.CommentId}); err != nil {
		return resp, err
	}

	resp = &core_api.GetCommentResp{
		SubjectId:       res.SubjectId,
		RootId:          res.RootId,
		FatherId:        res.FatherId,
		Count:           res.Count,
		State:           res.State,
		Attrs:           res.Attrs,
		Labels:          res.LabelIds,
		UserId:          res.UserId,
		AtUserId:        res.AtUserId,
		Content:         res.Content,
		Meta:            res.Meta,
		CreateAt:        res.CreateTime,
		Author:          &core_api.SimpleUser{},
		CommentRelation: &core_api.CommentRelation{},
		ItemType:        core_api.TargetType(res.Type),
	}

	_ = mr.Finish(func() error {
		s.CommentDomainService.LoadLikeCount(ctx, &resp.LikedCount, req.CommentId) // 点赞量
		return nil
	}, func() error {
		s.CommentDomainService.LoadHateCount(ctx, &resp.HatedCount, req.CommentId) // 点赞量
		return nil
	}, func() error {
		s.CommentDomainService.LoadAuthor(ctx, resp.Author, resp.UserId) // 作者
		return nil
	}, func() error {
		s.CommentDomainService.LoadLiked(ctx, resp.CommentRelation, req.CommentId, resp.UserId) // 是否点赞
		return nil
	}, func() error {
		s.CommentDomainService.LoadHated(ctx, resp.CommentRelation, req.CommentId, resp.UserId) // 是否点踩
		return nil
	}, func() error {
		s.CommentDomainService.LoadLabels(ctx, &resp.Labels, res.LabelIds) // 标签集
		return nil
	})
	return resp, nil
}

func (s *CommentService) GetComments(ctx context.Context, req *core_api.GetCommentsReq) (resp *core_api.GetCommentsResp, err error) {
	resp = new(core_api.GetCommentsResp)
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	var res *platform.GetCommentListResp
	p := convertor.MakePaginationOptions(req.Limit, req.Offset, req.LastToken, req.Backward)
	if res, err = s.Platform.GetCommentList(ctx, &platform.GetCommentListReq{
		FilterOptions: &platform.CommentFilterOptions{
			OnlyUserId:   req.OnlyUserId,
			OnlyAtUserId: req.OnlyAtUserId,
			OnlyState:    req.OnlyState,
			OnlyAttrs:    req.OnlyAttrs,
		},
		Pagination: p,
	}); err != nil {
		return resp, err
	}
	resp.Comments = lo.Map(res.Comments, func(item *platform.Comment, _ int) *core_api.Comment {
		c := convertor.CommentToCoreComment(item)
		_ = mr.Finish(func() error {
			s.CommentDomainService.LoadLikeCount(ctx, &c.Like, c.CommentId) // 点赞量
			return nil
		}, func() error {
			s.CommentDomainService.LoadLikeCount(ctx, &c.Hate, c.CommentId) // 点赞量
			return nil
		}, func() error {
			s.CommentDomainService.LoadAuthor(ctx, c.Author, item.UserId) // 作者
			return nil
		}, func() error {
			s.CommentDomainService.LoadLiked(ctx, c.CommentRelation, c.CommentId, item.UserId) // 是否点赞
			return nil
		}, func() error {
			s.CommentDomainService.LoadHated(ctx, c.CommentRelation, c.CommentId, item.UserId) // 是否点踩
			return nil
		}, func() error {
			s.CommentDomainService.LoadLabels(ctx, &c.Labels, item.Labels) // 标签集
			return nil
		}, func() error {
			//switch item.Type {
			//case int64(core_api.TargetType_FileType):
			//	s.CommentDomainService.LoadFile(ctx, &c.ItemTitle, &c.ItemUserId, item.SubjectId)
			//case int64(core_api.TargetType_PostType):
			s.CommentDomainService.LoadPost(ctx, &c.ItemTitle, &c.ItemUserId, item.SubjectId)
			//}
			return nil
		})
		return c
	})
	resp.Token = res.Token
	resp.Total = res.Total
	return resp, nil
}

func (s *CommentService) GetCommentBlocks(ctx context.Context, req *core_api.GetCommentBlocksReq) (resp *core_api.GetCommentBlocksResp, err error) {
	resp = new(core_api.GetCommentBlocksResp)
	var res *platform.GetCommentBlocksResp
	p := convertor.MakePaginationOptions(req.Limit, req.Offset, req.LastToken, req.Backward)
	if res, err = s.Platform.GetCommentBlocks(ctx, &platform.GetCommentBlocksReq{
		SubjectId:  req.SubjectId,
		RootId:     req.RootId,
		Pagination: p,
	}); err != nil {
		return resp, err
	}

	resp.CommentBlocks = lo.Map(res.CommentBlocks, func(item *platform.CommentBlock, _ int) *core_api.CommentBlock {
		var rootComment *core_api.CommentNode
		if item.RootComment != nil {
			rootComment = convertor.CommentToCoreCommentNode(item.RootComment)
			_ = mr.Finish(func() error {
				s.CommentDomainService.LoadLikeCount(ctx, &rootComment.Like, rootComment.CommentId) // 点赞量
				return nil
			}, func() error {
				s.CommentDomainService.LoadAuthor(ctx, rootComment.Author, item.RootComment.UserId) // 作者
				return nil
			}, func() error {
				s.CommentDomainService.LoadLiked(ctx, rootComment.CommentRelation, rootComment.CommentId, item.RootComment.UserId) // 是否点赞
				return nil
			}, func() error {
				s.CommentDomainService.LoadHated(ctx, rootComment.CommentRelation, rootComment.CommentId, item.RootComment.UserId) // 是否点踩
				return nil
			}, func() error {
				s.CommentDomainService.LoadLabels(ctx, &rootComment.Labels, item.RootComment.Labels) // 标签集
				return nil
			}, func() error {
				s.CommentDomainService.LoadAuthor(ctx, rootComment.AtUser, rootComment.AtUserId) // 回复者
				return nil
			})
		}
		comments := lo.Map(item.ReplyList.Comments, func(comment *platform.Comment, _ int) *core_api.CommentNode {
			c := convertor.CommentToCoreCommentNode(comment)
			_ = mr.Finish(func() error {
				s.CommentDomainService.LoadLikeCount(ctx, &c.Like, c.CommentId) // 点赞量
				return nil
			}, func() error {
				s.CommentDomainService.LoadHateCount(ctx, &c.Hate, c.CommentId) // 点赞量
				return nil
			}, func() error {
				s.CommentDomainService.LoadAuthor(ctx, c.Author, comment.UserId) // 作者
				return nil
			}, func() error {
				s.CommentDomainService.LoadLiked(ctx, c.CommentRelation, c.CommentId, comment.UserId) // 是否点赞
				return nil
			}, func() error {
				s.CommentDomainService.LoadHated(ctx, c.CommentRelation, c.CommentId, comment.UserId) // 是否点踩
				return nil
			}, func() error {
				s.CommentDomainService.LoadLabels(ctx, &c.Labels, comment.Labels) // 标签集
				return nil
			}, func() error {
				s.CommentDomainService.LoadAuthor(ctx, c.AtUser, comment.AtUserId) // 回复者
				return nil
			})
			return c
		})

		return &core_api.CommentBlock{
			RootComment: rootComment,
			ReplyList: &core_api.ReplyList{
				Comments: comments,
				Total:    item.ReplyList.Total,
				Token:    item.ReplyList.Token,
			},
		}
	})
	resp.Token = res.Token
	resp.Total = res.Total
	return resp, nil
}

func (s *CommentService) CreateComment(ctx context.Context, req *core_api.CreateCommentReq) (resp *core_api.CreateCommentResp, err error) {
	resp = new(core_api.CreateCommentResp)
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	var res *platform.CreateCommentResp
	if res, err = s.Platform.CreateComment(ctx, &platform.CreateCommentReq{
		SubjectId: req.SubjectId,
		RootId:    req.RootId,
		FatherId:  req.FatherId,
		LabelIds:  req.LabelIds,
		UserId:    userData.UserId,
		AtUserId:  req.AtUserId,
		Content:   req.Content,
		Meta:      req.Meta,
		Type:      int64(req.ItemType),
	}); err != nil {
		return resp, err
	}
	resp.CommentId = res.CommentId

	_, err = s.CloudMindContent.CreateHot(ctx, &content.CreateHotReq{
		HotId: res.CommentId,
	})
	if err != nil {
		return resp, err
	}

	err = s.RelationDomainService.CreateRelation(ctx, &core_api.Relation{
		FromType:     core_api.TargetType_UserType,
		FromId:       userData.UserId,
		ToType:       core_api.TargetType_CommentContentType,
		ToId:         res.CommentId,
		RelationType: core_api.RelationType_CommentRelationType,
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *CommentService) UpdateComment(ctx context.Context, req *core_api.UpdateCommentReq) (resp *core_api.UpdateCommentResp, err error) {
	resp = new(core_api.UpdateCommentResp)
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	var res *platform.GetCommentResp
	if res, err = s.Platform.GetComment(ctx, &platform.GetCommentReq{CommentId: req.CommentId}); err != nil {
		return resp, err
	}
	if res.UserId == userData.UserId {
		if _, err = s.Platform.UpdateComment(ctx, &platform.UpdateCommentReq{
			CommentId: req.CommentId,
			State:     int64(req.State),
			LabelIds:  req.LabelIds,
		}); err != nil {
			return resp, err
		}
	}
	return resp, nil
}

func (s *CommentService) DeleteComment(ctx context.Context, req *core_api.DeleteCommentReq) (resp *core_api.DeleteCommentResp, err error) {
	resp = new(core_api.DeleteCommentResp)
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	var ok bool
	var res *platform.GetCommentResp
	if res, err = s.Platform.GetComment(ctx, &platform.GetCommentReq{CommentId: req.CommentId}); err != nil {
		return resp, err
	}
	switch {
	case res.UserId == userData.UserId:
		ok = true
	default:
		var result *platform.GetCommentSubjectResp
		if result, err = s.Platform.GetCommentSubject(ctx, &platform.GetCommentSubjectReq{SubjectId: res.SubjectId}); err != nil {
			return resp, err
		}
		ok = result.UserId == userData.UserId
	}
	if ok {
		if _, err = s.Platform.DeleteComment(ctx, &platform.DeleteCommentReq{CommentId: req.CommentId}); err != nil {
			return resp, err
		}
	}
	return resp, nil
}

func (s *CommentService) SetCommentAttrs(ctx context.Context, req *core_api.SetCommentAttrsReq) (resp *core_api.SetCommentAttrsResp, err error) {
	resp = new(core_api.SetCommentAttrsResp)
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	var res *platform.GetCommentResp
	if res, err = s.Platform.GetComment(ctx, &platform.GetCommentReq{CommentId: req.Id}); err != nil {
		return resp, err
	}
	if res.UserId == userData.UserId {
		if _, err = s.Platform.SetCommentAttrs(ctx, &platform.SetCommentAttrsReq{CommentId: req.Id, Attrs: int64(req.Attrs), SubjectId: res.SubjectId, SortTime: res.CreateTime}); err != nil {
			return resp, err
		}
	}
	return resp, nil
}

func (s *CommentService) GetCommentSubject(ctx context.Context, req *core_api.GetCommentSubjectReq) (resp *core_api.GetCommentSubjectResp, err error) {
	resp = new(core_api.GetCommentSubjectResp)
	var res *platform.GetCommentSubjectResp
	if res, err = s.Platform.GetCommentSubject(ctx, &platform.GetCommentSubjectReq{SubjectId: req.SubjectId}); err != nil {
		return resp, err
	}
	resp = &core_api.GetCommentSubjectResp{
		UserId:       res.UserId,
		TopCommentId: res.TopCommentId,
		RootCount:    res.RootCount,
		AllCount:     res.AllCount,
		State:        res.State,
		Attrs:        res.Attrs,
	}
	return resp, nil
}

func (s *CommentService) UpdateCommentSubject(ctx context.Context, req *core_api.UpdateCommentSubjectReq) (resp *core_api.UpdateCommentSubjectResp, err error) {
	resp = new(core_api.UpdateCommentSubjectResp)
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	var res *platform.GetCommentSubjectResp
	if res, err = s.Platform.GetCommentSubject(ctx, &platform.GetCommentSubjectReq{SubjectId: req.SubjectId}); err != nil {
		return resp, err
	}
	if res.UserId == userData.UserId {
		if _, err = s.Platform.UpdateCommentSubject(ctx, &platform.UpdateCommentSubjectReq{
			SubjectId: req.SubjectId,
			State:     int64(req.State),
			Attrs:     int64(req.Attrs),
		}); err != nil {
			return resp, err
		}
	}
	return resp, nil
}

func (s *CommentService) DeleteCommentSubject(ctx context.Context, req *core_api.DeleteCommentSubjectReq) (resp *core_api.DeleteCommentSubjectResp, err error) {
	resp = new(core_api.DeleteCommentSubjectResp)
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	var res *platform.GetCommentSubjectResp
	if res, err = s.Platform.GetCommentSubject(ctx, &platform.GetCommentSubjectReq{SubjectId: req.SubjectId}); err != nil {
		return resp, err
	}
	if res.UserId == userData.UserId {
		if _, err = s.Platform.DeleteCommentSubject(ctx, &platform.DeleteCommentSubjectReq{SubjectId: req.SubjectId}); err != nil {
			return resp, err
		}
	}
	return resp, nil
}
