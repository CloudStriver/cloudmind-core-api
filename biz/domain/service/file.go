package service

import (
	"context"
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_content"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/platform_comment"
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
	LoadLabels(ctx context.Context, file *core_api.PublicFile, labelIds []string)
}
type FileDomainService struct {
	CloudMindUser    cloudmind_content.ICloudMindContent
	PlatformRelation platform_relation.IPlatFormRelation
	PlatformComment  platform_comment.IPlatFormComment
}

var FileDomainServiceSet = wire.NewSet(
	wire.Struct(new(FileDomainService), "*"),
	wire.Bind(new(IFileDomainService), new(*FileDomainService)),
)

func (s *FileDomainService) LoadCollected(ctx context.Context, file *core_api.PublicFile, userId string) {
	getRelationResp, err := s.PlatformRelation.GetRelation(ctx, &relation.GetRelationReq{
		FromType:     int64(core_api.TargetType_UserType),
		FromId:       userId,
		ToType:       int64(core_api.TargetType_FileType),
		ToId:         file.FileId,
		RelationType: int64(core_api.RelationType_CollectType),
	})
	if err == nil {
		file.FileRelation.Collected = getRelationResp.Ok
	}
}

func (s *FileDomainService) LoadViewCount(ctx context.Context, file *core_api.PublicFile) {
	getRelationCountResp, err := s.PlatformRelation.GetRelationCount(ctx, &relation.GetRelationCountReq{
		RelationFilterOptions: &relation.GetRelationCountReq_ToFilterOptions{
			ToFilterOptions: &relation.ToFilterOptions{
				ToType:   int64(core_api.TargetType_FileType),
				ToId:     file.FileId,
				FromType: int64(core_api.TargetType_UserType),
			},
		},
		RelationType: int64(core_api.RelationType_ViewType),
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
				ToType:   int64(core_api.TargetType_FileType),
				ToId:     file.FileId,
				FromType: int64(core_api.TargetType_UserType),
			},
		},
		RelationType: int64(core_api.RelationType_LikeType),
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
		FromType:     int64(core_api.TargetType_UserType),
		FromId:       userId,
		ToType:       int64(core_api.TargetType_FileType),
		ToId:         file.FileId,
		RelationType: int64(core_api.RelationType_LikeType),
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
				ToType:   int64(core_api.TargetType_FileType),
				ToId:     file.FileId,
				FromType: int64(core_api.TargetType_UserType),
			},
		},
		RelationType: int64(core_api.RelationType_CollectType),
	})
	if err == nil {
		file.FileCount.CollectCount = getRelationCountResp.Total
	}
}

func (s *FileDomainService) LoadLabels(ctx context.Context, file *core_api.PublicFile, labelIds []string) {
	if file.Zone == "" || file.SubZone == "" {
		return
	}
	var labels *comment.GetLabelsInBatchResp
	labels, _ = s.PlatformComment.GetLabelsInBatch(ctx, &comment.GetLabelsInBatchReq{LabelIds: labelIds})
	file.Labels = lo.Map(labels.Labels, func(item string, index int) *core_api.Label {
		return &core_api.Label{LabelId: labelIds[index], Value: item}
	})
}
