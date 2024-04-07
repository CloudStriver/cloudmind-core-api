package service

import (
	"context"
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_content"
	platformservice "github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/platform"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/platform"
	"github.com/google/wire"
	"github.com/samber/lo"
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
	CloudMindUser cloudmind_content.ICloudMindContent
	Platform      platformservice.IPlatForm
}

var FileDomainServiceSet = wire.NewSet(
	wire.Struct(new(FileDomainService), "*"),
	wire.Bind(new(IFileDomainService), new(*FileDomainService)),
)

func (s *FileDomainService) LoadCollected(ctx context.Context, file *core_api.PublicFile, userId string) {
	getRelationResp, err := s.Platform.GetRelation(ctx, &platform.GetRelationReq{
		FromType:     int64(core_api.TargetType_UserType),
		FromId:       userId,
		ToType:       int64(core_api.TargetType_FileType),
		ToId:         file.FileId,
		RelationType: int64(core_api.RelationType_CollectRelationType),
	})
	if err == nil {
		file.FileRelation.Collected = getRelationResp.Ok
	}
}

func (s *FileDomainService) LoadViewCount(ctx context.Context, file *core_api.PublicFile) {
	getRelationCountResp, err := s.Platform.GetRelationCount(ctx, &platform.GetRelationCountReq{
		RelationFilterOptions: &platform.GetRelationCountReq_ToFilterOptions{
			ToFilterOptions: &platform.ToFilterOptions{
				ToType:   int64(core_api.TargetType_FileType),
				ToId:     file.FileId,
				FromType: int64(core_api.TargetType_UserType),
			},
		},
		RelationType: int64(core_api.RelationType_ViewRelationType),
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
	getRelationCountResp, err := s.Platform.GetRelationCount(ctx, &platform.GetRelationCountReq{
		RelationFilterOptions: &platform.GetRelationCountReq_ToFilterOptions{
			ToFilterOptions: &platform.ToFilterOptions{
				ToType:   int64(core_api.TargetType_FileType),
				ToId:     file.FileId,
				FromType: int64(core_api.TargetType_UserType),
			},
		},
		RelationType: int64(core_api.RelationType_LikeRelationType),
	})
	if err == nil {
		file.FileCount.LikeCount = getRelationCountResp.Total
	}
}

func (s *FileDomainService) LoadLiked(ctx context.Context, file *core_api.PublicFile, userId string) {
	if file.Zone == "" || file.SubZone == "" {
		return
	}
	getRelationResp, err := s.Platform.GetRelation(ctx, &platform.GetRelationReq{
		FromType:     int64(core_api.TargetType_UserType),
		FromId:       userId,
		ToType:       int64(core_api.TargetType_FileType),
		ToId:         file.FileId,
		RelationType: int64(core_api.RelationType_LikeRelationType),
	})
	if err == nil {
		file.FileRelation.Liked = getRelationResp.Ok
	}
}

func (s *FileDomainService) LoadCollectCount(ctx context.Context, file *core_api.PublicFile) {
	if file.Zone == "" || file.SubZone == "" {
		return
	}
	getRelationCountResp, err := s.Platform.GetRelationCount(ctx, &platform.GetRelationCountReq{
		RelationFilterOptions: &platform.GetRelationCountReq_ToFilterOptions{
			ToFilterOptions: &platform.ToFilterOptions{
				ToType:   int64(core_api.TargetType_FileType),
				ToId:     file.FileId,
				FromType: int64(core_api.TargetType_UserType),
			},
		},
		RelationType: int64(core_api.RelationType_CollectRelationType),
	})
	if err == nil {
		file.FileCount.CollectCount = getRelationCountResp.Total
	}
}

func (s *FileDomainService) LoadLabels(ctx context.Context, file *core_api.PublicFile, labelIds []string) {
	if file.Zone == "" || file.SubZone == "" {
		return
	}
	var labels *platform.GetLabelsInBatchResp
	labels, _ = s.Platform.GetLabelsInBatch(ctx, &platform.GetLabelsInBatchReq{LabelIds: labelIds})
	file.Labels = lo.Map(labels.Labels, func(item string, index int) *core_api.Label {
		return &core_api.Label{LabelId: labelIds[index], Value: item}
	})
}
