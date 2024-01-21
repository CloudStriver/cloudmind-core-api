package service

import (
	"context"
	"github.com/CloudStriver/cloudmind-core-api/biz/adaptor"
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/consts"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/convertor"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_content"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/platform_relation"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/basic"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/platform/relation"
	"github.com/google/wire"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/mr"
)

type IRelationService interface {
	GetRelation(ctx context.Context, req *core_api.GetRelationReq) (resp *core_api.GetRelationResp, err error)
	CreateRelation(ctx context.Context, req *core_api.CreateRelationReq) (resp *core_api.CreateRelationResp, err error)
	GetFromRelations(ctx context.Context, c *core_api.GetFromRelationsReq) (resp *core_api.GetFromRelationsResp, err error)
	GetToRelations(ctx context.Context, c *core_api.GetToRelationsReq) (resp *core_api.GetToRelationsResp, err error)
	DeleteRelation(ctx context.Context, c *core_api.DeleteRelationReq) (resp *core_api.DeleteRelationResp, err error)
}

var RelationServiceSet = wire.NewSet(
	wire.Struct(new(RelationService), "*"),
	wire.Bind(new(IRelationService), new(*RelationService)),
)

type RelationService struct {
	Config           *config.Config
	PlatFormRelation platform_relation.IPlatFormRelation
	CloudMindContent cloudmind_content.ICloudMindContent
}

func (s *RelationService) GetFromRelations(ctx context.Context, req *core_api.GetFromRelationsReq) (resp *core_api.GetFromRelationsResp, err error) {
	resp = new(core_api.GetFromRelationsResp)
	getFromRelationsResp, err := s.PlatFormRelation.GetRelations(ctx, &relation.GetRelationsReq{
		RelationFilterOptions: &relation.GetRelationsReq_FromFilterOptions{
			FromFilterOptions: &relation.FromFilterOptions{
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
		err = mr.Finish(lo.Map[*relation.Relation](getFromRelationsResp.Relations, func(relation *relation.Relation, i int) func() error {
			return func() error {
				user, err := s.CloudMindContent.GetUser(ctx, &content.GetUserReq{
					UserId: relation.ToId,
				})
				if err != nil {
					return err
				}
				resp.Users[i] = convertor.UserDetailToUser(user.User)
				return nil
			}
		})...)
		if err != nil {
			return resp, err
		}
	default:
		//return resp, consts.ErrNotSupportTargetType
	}
	return resp, nil
}

func (s *RelationService) GetToRelations(ctx context.Context, req *core_api.GetToRelationsReq) (resp *core_api.GetToRelationsResp, err error) {
	resp = new(core_api.GetToRelationsResp)
	getFromRelationsResp, err := s.PlatFormRelation.GetRelations(ctx, &relation.GetRelationsReq{
		RelationFilterOptions: &relation.GetRelationsReq_ToFilterOptions{
			ToFilterOptions: &relation.ToFilterOptions{
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
		err = mr.Finish(lo.Map[*relation.Relation](getFromRelationsResp.Relations, func(relation *relation.Relation, i int) func() error {
			return func() error {
				user, err := s.CloudMindContent.GetUser(ctx, &content.GetUserReq{
					UserId: relation.FromId,
				})
				if err != nil {
					return err
				}
				resp.Users[i] = convertor.UserDetailToUser(user.User)
				return nil
			}
		})...)
		if err != nil {
			return resp, err
		}
	default:
		//return resp, consts.ErrNotSupportTargetType
	}
	return resp, nil
}

func (s *RelationService) DeleteRelation(ctx context.Context, req *core_api.DeleteRelationReq) (resp *core_api.DeleteRelationResp, err error) {
	resp = new(core_api.DeleteRelationResp)
	user := adaptor.ExtractUserMeta(ctx)
	if user.GetUserId() != req.RelationInfo.FromId {
		return resp, consts.ErrNotPermission
	}
	if _, err = s.PlatFormRelation.DeleteRelation(ctx, &relation.DeleteRelationReq{
		RelationInfo: convertor.CoreApiRelationInfoToRelationInfo(req.RelationInfo),
	}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *RelationService) GetRelation(ctx context.Context, req *core_api.GetRelationReq) (resp *core_api.GetRelationResp, err error) {
	resp = new(core_api.GetRelationResp)
	getRelationResp, err := s.PlatFormRelation.GetRelation(ctx, &relation.GetRelationReq{
		RelationInfo: &relation.RelationInfo{
			FromType:     req.FromType,
			FromId:       req.FromId,
			ToType:       req.ToType,
			ToId:         req.ToId,
			RelationType: req.RelationType,
		},
	})
	if err != nil {
		return resp, err
	}
	resp.Ok = getRelationResp.Ok
	return resp, nil
}

func (s *RelationService) CreateRelation(ctx context.Context, req *core_api.CreateRelationReq) (resp *core_api.CreateRelationResp, err error) {
	resp = new(core_api.CreateRelationResp)
	user := adaptor.ExtractUserMeta(ctx)
	if user.GetUserId() != req.Relation.FromId {
		return resp, consts.ErrNotPermission
	}

	if _, err = s.PlatFormRelation.CreateRelation(ctx, &relation.CreateRelationReq{
		RelationInfo: convertor.CoreApiRelationInfoToRelationInfo(req.Relation),
	}); err != nil {
		return resp, err
	}
	return resp, nil
}
