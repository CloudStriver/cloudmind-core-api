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
	LoadAuthor(ctx context.Context, post *core_api.Post, userId string)
	LoadLikeCount(ctx context.Context, post *core_api.Post)
	LoadLiked(ctx context.Context, post *core_api.Post, userId string)
	LoadCollectCount(ctx context.Context, post *core_api.Post)
}
type PostDomainService struct {
	CloudMindUser    cloudmind_content.ICloudMindContent
	PlatformRelation platform_relation.IPlatFormRelation
}

var PostDomainServiceSet = wire.NewSet(
	wire.Struct(new(PostDomainService), "*"),
	wire.Bind(new(IPostDomainService), new(*PostDomainService)),
)

func (s *PostDomainService) LoadAuthor(ctx context.Context, post *core_api.Post, userId string) {
	if userId == "" {
		return
	}
	post.Author = &core_api.User{
		UserId: userId,
	}
	getUserResp, err := s.CloudMindUser.GetUser(ctx, &content.GetUserReq{UserId: userId})
	if err == nil {
		post.Author.Name = getUserResp.User.Name
		post.Author.Url = getUserResp.User.Url
	}
}

func (s *PostDomainService) LoadLikeCount(ctx context.Context, post *core_api.Post) {
	getRelationCountResp, err := s.PlatformRelation.GetRelationCount(ctx, &relation.GetRelationCountReq{
		RelationFilterOptions: &relation.GetRelationCountReq_ToFilterOptions{
			ToFilterOptions: &relation.ToFilterOptions{
				ToType:   consts.RelationPostType,
				ToId:     post.PostId,
				FromType: consts.RelationUserType,
			},
		},
		RelationType: consts.RelationLikeType,
	})
	if err == nil {
		post.PostCount.LikeCount = getRelationCountResp.Total
	}
}

func (s *PostDomainService) LoadLiked(ctx context.Context, post *core_api.Post, userId string) {
	getRelationResp, err := s.PlatformRelation.GetRelation(ctx, &relation.GetRelationReq{
		Relation: &relation.Relation{
			FromType:     consts.RelationUserType,
			FromId:       userId,
			ToType:       consts.RelationPostType,
			ToId:         post.PostId,
			RelationType: consts.RelationLikeType,
		},
	})
	if err == nil {
		post.PostRelation.Liked = getRelationResp.Ok
	}
}

func (s *PostDomainService) LoadCollectCount(ctx context.Context, post *core_api.Post) {
	getRelationCountResp, err := s.PlatformRelation.GetRelationCount(ctx, &relation.GetRelationCountReq{
		RelationFilterOptions: &relation.GetRelationCountReq_ToFilterOptions{
			ToFilterOptions: &relation.ToFilterOptions{
				ToType:   consts.RelationPostType,
				ToId:     post.PostId,
				FromType: consts.RelationUserType,
			},
		},
		RelationType: consts.RelationCollectType,
	})
	if err == nil {
		post.PostCount.CollectCount = getRelationCountResp.Total
	}
}
