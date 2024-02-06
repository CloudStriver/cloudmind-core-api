package service

import (
	"context"
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/consts"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_content"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/platform_relation"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/platform/relation"
	"github.com/google/wire"
)

type IPostDomainService interface {
	LoadAuthor(ctx context.Context, user *core_api.User, userId string)
	LoadLikeCount(ctx context.Context, likeCount *int64, postId string)
	LoadViewCount(ctx context.Context, viewCount *int64, postId string)
	LoadCollectCount(ctx context.Context, collectCount *int64, postId string)
	LoadLiked(ctx context.Context, liked *bool, userId, postId string)
	LoadCollected(ctx context.Context, collected *bool, userId, postId string)
}
type PostDomainService struct {
	CloudMindUser    cloudmind_content.ICloudMindContent
	PlatformRelation platform_relation.IPlatFormRelation
}

var PostDomainServiceSet = wire.NewSet(
	wire.Struct(new(PostDomainService), "*"),
	wire.Bind(new(IPostDomainService), new(*PostDomainService)),
)

func (s *PostDomainService) LoadAuthor(ctx context.Context, user *core_api.User, userId string) {
	if userId == "" {
		return
	}
	getUserResp, err := s.CloudMindUser.GetUser(ctx, &content.GetUserReq{UserId: userId})
	if err == nil {
		user.Name = getUserResp.User.Name
		user.Url = getUserResp.User.Url
		user.UserId = userId
	}
}

func (s *PostDomainService) LoadLikeCount(ctx context.Context, likeCount *int64, postId string) {
	getRelationCountResp, err := s.PlatformRelation.GetRelationCount(ctx, &relation.GetRelationCountReq{
		RelationFilterOptions: &relation.GetRelationCountReq_ToFilterOptions{
			ToFilterOptions: &relation.ToFilterOptions{
				ToType:   consts.RelationPostType,
				ToId:     postId,
				FromType: consts.RelationUserType,
			},
		},
		RelationType: consts.RelationLikeType,
	})
	if err == nil {
		*likeCount = getRelationCountResp.Total
	}
}

func (s *PostDomainService) LoadViewCount(ctx context.Context, viewCount *int64, postId string) {
	getRelationCountResp, err := s.PlatformRelation.GetRelationCount(ctx, &relation.GetRelationCountReq{
		RelationFilterOptions: &relation.GetRelationCountReq_ToFilterOptions{
			ToFilterOptions: &relation.ToFilterOptions{
				ToType:   consts.RelationPostType,
				ToId:     postId,
				FromType: consts.RelationUserType,
			},
		},
		RelationType: consts.RelationViewType,
	})
	if err == nil {
		*viewCount = getRelationCountResp.Total
	}
}

func (s *PostDomainService) LoadCollectCount(ctx context.Context, collectCount *int64, postId string) {
	getRelationCountResp, err := s.PlatformRelation.GetRelationCount(ctx, &relation.GetRelationCountReq{
		RelationFilterOptions: &relation.GetRelationCountReq_ToFilterOptions{
			ToFilterOptions: &relation.ToFilterOptions{
				ToType:   consts.RelationPostType,
				ToId:     postId,
				FromType: consts.RelationUserType,
			},
		},
		RelationType: consts.RelationCollectType,
	})
	if err == nil {
		*collectCount = getRelationCountResp.Total
	}
}

func (s *PostDomainService) LoadLiked(ctx context.Context, liked *bool, userId, postId string) {
	getRelationResp, err := s.PlatformRelation.GetRelation(ctx, &relation.GetRelationReq{
		Relation: &relation.Relation{
			FromType:     consts.RelationUserType,
			FromId:       userId,
			ToType:       consts.RelationPostType,
			ToId:         postId,
			RelationType: consts.RelationLikeType,
		},
	})
	if err == nil {
		*liked = getRelationResp.Ok
	}
}

func (s *PostDomainService) LoadCollected(ctx context.Context, collected *bool, userId, postId string) {
	getRelationResp, err := s.PlatformRelation.GetRelation(ctx, &relation.GetRelationReq{
		Relation: &relation.Relation{
			FromType:     consts.RelationUserType,
			FromId:       userId,
			ToType:       consts.RelationPostType,
			ToId:         postId,
			RelationType: consts.RelationCollectType,
		},
	})
	if err == nil {
		*collected = getRelationResp.Ok
	}
}
