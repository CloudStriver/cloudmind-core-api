package service

import (
	"context"
	"github.com/CloudStriver/cloudmind-core-api/biz/adaptor"
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/cloudmind-core-api/biz/domain/service"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/consts"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_content"
	platformservice "github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/platform"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/basic"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/platform"
	"github.com/google/wire"
)

type IRelationService interface {
	GetRelation(ctx context.Context, req *core_api.GetRelationReq) (resp *core_api.GetRelationResp, err error)
	CreateRelation(ctx context.Context, req *core_api.CreateRelationReq) (resp *core_api.CreateRelationResp, err error)
	GetFromRelations(ctx context.Context, req *core_api.GetFromRelationsReq) (resp *core_api.GetFromRelationsResp, err error)
	GetToRelations(ctx context.Context, req *core_api.GetToRelationsReq) (resp *core_api.GetToRelationsResp, err error)
	DeleteRelation(ctx context.Context, req *core_api.DeleteRelationReq) (resp *core_api.DeleteRelationResp, err error)
	GetRelationPaths(ctx context.Context, req *core_api.GetRelationPathsReq) (resp *core_api.GetRelationPathsResp, err error)
}

var RelationServiceSet = wire.NewSet(
	wire.Struct(new(RelationService), "*"),
	wire.Bind(new(IRelationService), new(*RelationService)),
)

type RelationService struct {
	Config                *config.Config
	CloudMindContent      cloudmind_content.ICloudMindContent
	Platform              platformservice.IPlatForm
	UserDomainService     service.IUserDomainService
	PostDomainService     service.IPostDomainService
	RelationDomainService service.RelationDomainService
}

func (s *RelationService) GetRelationPaths(ctx context.Context, req *core_api.GetRelationPathsReq) (resp *core_api.GetRelationPathsResp, err error) {
	resp = new(core_api.GetRelationPathsResp)
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	relationPaths, err := s.Platform.GetRelationPaths(ctx, &platform.GetRelationPathsReq{
		FromId:    userData.UserId,
		FromType:  int64(core_api.TargetType_UserType),
		EdgeType1: int64(core_api.RelationType_FollowRelationType),
		EdgeType2: int64(req.RelationType),
		PaginationOptions: &basic.PaginationOptions{
			Limit:  req.Limit,
			Offset: req.Offset,
		},
	})
	if err != nil {
		return resp, err
	}

	switch req.RelationType {
	case core_api.RelationType_FollowRelationType:
		// 用户
		resp.Users = make([]*core_api.User, len(relationPaths.Relations))
		if err = s.RelationDomainService.GetUserByRelations(ctx, relationPaths.Relations, resp.Users, userData.GetUserId()); err != nil {
			return resp, err
		}
	case core_api.RelationType_PublishRelationType:
		// 文章
		resp.Posts = make([]*core_api.Post, len(relationPaths.Relations))
		if err = s.RelationDomainService.GetPostByRelations(ctx, relationPaths.Relations, resp.Posts, userData.GetUserId()); err != nil {
			return resp, err
		}
	default:
		return resp, consts.ErrNotSupportRelationType
	}
	return resp, nil
}

func (s *RelationService) GetFromRelations(ctx context.Context, req *core_api.GetFromRelationsReq) (resp *core_api.GetFromRelationsResp, err error) {
	resp = new(core_api.GetFromRelationsResp)
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil {
		return resp, consts.ErrNotAuthentication
	}
	getFromRelationsResp, err := s.Platform.GetRelations(ctx, &platform.GetRelationsReq{
		RelationFilterOptions: &platform.GetRelationsReq_FromFilterOptions{
			FromFilterOptions: &platform.FromFilterOptions{
				FromType: req.FromType,
				FromId:   req.FromId,
				ToType:   int64(req.ToType),
			},
		},
		RelationType: req.RelationType,
		PaginationOptions: &basic.PaginationOptions{
			Limit:  req.Limit,
			Offset: req.Offset,
		},
	})
	if err != nil {
		return resp, err
	}

	switch req.ToType {
	case core_api.TargetType_UserType:
		resp.Users = make([]*core_api.User, len(getFromRelationsResp.Relations))
		if err = s.RelationDomainService.GetUserByRelations(ctx, getFromRelationsResp.Relations, resp.Users, userData.GetUserId()); err != nil {
			return resp, err
		}
	case core_api.TargetType_PostType:
		// 文章
		resp.Posts = make([]*core_api.Post, len(getFromRelationsResp.Relations))
		if err = s.RelationDomainService.GetPostByRelations(ctx, getFromRelationsResp.Relations, resp.Posts, userData.GetUserId()); err != nil {
			return resp, err
		}
	default:
		return resp, consts.ErrNotSupportRelationType
	}
	return resp, nil
}

func (s *RelationService) GetToRelations(ctx context.Context, req *core_api.GetToRelationsReq) (resp *core_api.GetToRelationsResp, err error) {
	resp = new(core_api.GetToRelationsResp)
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil {
		return resp, consts.ErrNotAuthentication
	}
	getFromRelationsResp, err := s.Platform.GetRelations(ctx, &platform.GetRelationsReq{
		RelationFilterOptions: &platform.GetRelationsReq_ToFilterOptions{
			ToFilterOptions: &platform.ToFilterOptions{
				ToType:   req.ToType,
				ToId:     req.ToId,
				FromType: int64(req.FromType),
			},
		},
		RelationType: req.RelationType,
		PaginationOptions: &basic.PaginationOptions{
			Limit:  req.Limit,
			Offset: req.Offset,
		},
	})
	if err != nil {
		return resp, err
	}

	switch req.FromType {
	case core_api.TargetType_UserType:
		resp.Users = make([]*core_api.User, len(getFromRelationsResp.Relations))
		if err = s.RelationDomainService.GetUserByRelations(ctx, getFromRelationsResp.Relations, resp.Users, userData.GetUserId()); err != nil {
			return resp, err
		}
	default:
		return resp, consts.ErrNotSupportRelationType
	}
	return resp, nil
}

func (s *RelationService) DeleteRelation(ctx context.Context, req *core_api.DeleteRelationReq) (resp *core_api.DeleteRelationResp, err error) {
	user, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || user.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	if _, err = s.Platform.DeleteRelation(ctx, &platform.DeleteRelationReq{
		FromType:     int64(core_api.TargetType_UserType),
		FromId:       user.UserId,
		ToType:       req.ToType,
		ToId:         req.ToId,
		RelationType: req.RelationType,
	}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *RelationService) GetRelation(ctx context.Context, req *core_api.GetRelationReq) (resp *core_api.GetRelationResp, err error) {
	resp = new(core_api.GetRelationResp)
	getRelationResp, err := s.Platform.GetRelation(ctx, &platform.GetRelationReq{
		FromType:     req.FromType,
		FromId:       req.FromId,
		ToType:       req.ToType,
		ToId:         req.ToId,
		RelationType: req.RelationType,
	})
	if err != nil {
		return resp, err
	}

	resp.Ok = getRelationResp.Ok
	return resp, nil
}

func (s *RelationService) CreateRelation(ctx context.Context, req *core_api.CreateRelationReq) (resp *core_api.CreateRelationResp, err error) {
	user, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || user.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	if err = s.RelationDomainService.CreateRelation(ctx, &core_api.Relation{
		FromType:     core_api.TargetType_UserType,
		FromId:       user.GetUserId(),
		ToType:       req.ToType,
		ToId:         req.ToId,
		RelationType: req.RelationType,
	}); err != nil {
		return resp, err
	}
	return resp, nil
}
