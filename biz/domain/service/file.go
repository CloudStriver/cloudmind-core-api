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

type IFileDomainService interface {
	LoadAuthor(ctx context.Context, file *core_api.PublicFile, userId string)
	LoadLikeCount(ctx context.Context, file *core_api.PublicFile)
	LoadViewCount(ctx context.Context, file *core_api.PublicFile)
	LoadLiked(ctx context.Context, file *core_api.PublicFile, userId string)
	LoadCollectCount(ctx context.Context, file *core_api.PublicFile)
	LoadCollected(ctx context.Context, file *core_api.PublicFile, userId string)
}
type FileDomainService struct {
	CloudMindUser    cloudmind_content.ICloudMindContent
	PlatformRelation platform_relation.IPlatFormRelation
}

var FileDomainServiceSet = wire.NewSet(
	wire.Struct(new(FileDomainService), "*"),
	wire.Bind(new(IFileDomainService), new(*FileDomainService)),
)

func (s *FileDomainService) LoadCollected(ctx context.Context, file *core_api.PublicFile, userId string) {
	getRelationResp, err := s.PlatformRelation.GetRelation(ctx, &relation.GetRelationReq{
		Relation: &relation.Relation{
			FromType:     consts.RelationUserType,
			FromId:       userId,
			ToType:       consts.RelationFileType,
			ToId:         file.FileId,
			RelationType: consts.RelationCollectType,
		},
	})
	if err == nil {
		file.FileRelation.Collected = getRelationResp.Ok
	}
}

func (s *FileDomainService) LoadViewCount(ctx context.Context, file *core_api.PublicFile) {
	getRelationCountResp, err := s.PlatformRelation.GetRelationCount(ctx, &relation.GetRelationCountReq{
		RelationFilterOptions: &relation.GetRelationCountReq_ToFilterOptions{
			ToFilterOptions: &relation.ToFilterOptions{
				ToType:   consts.RelationFileType,
				ToId:     file.FileId,
				FromType: consts.RelationUserType,
			},
		},
		RelationType: consts.RelationViewType,
	})
	if err == nil {
		file.FileCount.ViewCount = getRelationCountResp.Total
	}
}
func (s *FileDomainService) LoadAuthor(ctx context.Context, file *core_api.PublicFile, userId string) {
	if userId == "" || file.Zone == "" || file.SubZone == "" {
		return
	}
	file.Author = &core_api.User{
		UserId: userId,
	}
	getUserResp, err := s.CloudMindUser.GetUser(ctx, &content.GetUserReq{UserId: userId})
	if err == nil {
		file.Author.Name = getUserResp.Name
		file.Author.Url = getUserResp.Url
	}
}

func (s *FileDomainService) LoadLikeCount(ctx context.Context, file *core_api.PublicFile) {
	if file.Zone == "" || file.SubZone == "" {
		return
	}
	getRelationCountResp, err := s.PlatformRelation.GetRelationCount(ctx, &relation.GetRelationCountReq{
		RelationFilterOptions: &relation.GetRelationCountReq_ToFilterOptions{
			ToFilterOptions: &relation.ToFilterOptions{
				ToType:   consts.RelationFileType,
				ToId:     file.FileId,
				FromType: consts.RelationUserType,
			},
		},
		RelationType: consts.RelationLikeType,
	})
	if err == nil {
		file.FileCount.LikeCount = getRelationCountResp.Total
	}
}

func (s *FileDomainService) LoadLiked(ctx context.Context, file *core_api.PublicFile, userId string) {
	if file.Zone == "" || file.SubZone == "" {
		return
	}
	getRelationResp, err := s.PlatformRelation.GetRelation(ctx, &relation.GetRelationReq{
		Relation: &relation.Relation{
			FromType:     consts.RelationUserType,
			FromId:       userId,
			ToType:       consts.RelationFileType,
			ToId:         file.FileId,
			RelationType: consts.RelationLikeType,
		},
	})
	if err == nil {
		file.FileRelation.Liked = getRelationResp.Ok
	}
}

func (s *FileDomainService) LoadCollectCount(ctx context.Context, file *core_api.PublicFile) {
	if file.Zone == "" || file.SubZone == "" {
		return
	}
	getRelationCountResp, err := s.PlatformRelation.GetRelationCount(ctx, &relation.GetRelationCountReq{
		RelationFilterOptions: &relation.GetRelationCountReq_ToFilterOptions{
			ToFilterOptions: &relation.ToFilterOptions{
				ToType:   consts.RelationFileType,
				ToId:     file.FileId,
				FromType: consts.RelationUserType,
			},
		},
		RelationType: consts.RelationCollectType,
	})
	if err == nil {
		file.FileCount.CollectCount = getRelationCountResp.Total
	}
}
