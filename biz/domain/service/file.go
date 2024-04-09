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
	LoadAuthor(ctx context.Context, c *core_api.FileUser, userId string)
	LoadLikeCount(ctx context.Context, c *core_api.FileCount, id string)
	LoadViewCount(ctx context.Context, c *core_api.FileCount, id string)
	LoadLiked(ctx context.Context, c *core_api.FileRelation, id, userId string)
	LoadCollectCount(ctx context.Context, c *core_api.FileCount, id string)
	LoadCollected(ctx context.Context, c *core_api.FileRelation, id, userId string)
	LoadLabels(ctx context.Context, c *[]*core_api.Label, labelIds []string)
}
type FileDomainService struct {
	CloudMindUser cloudmind_content.ICloudMindContent
	Platform      platformservice.IPlatForm
}

var FileDomainServiceSet = wire.NewSet(
	wire.Struct(new(FileDomainService), "*"),
	wire.Bind(new(IFileDomainService), new(*FileDomainService)),
)

func (s *FileDomainService) LoadAuthor(ctx context.Context, c *core_api.FileUser, userId string) {
	if userId == "" {
		return
	}
	c.UserId = userId
	getUserResp, err := s.CloudMindUser.GetUser(ctx, &content.GetUserReq{UserId: userId})
	if err == nil {
		c.Name = getUserResp.Name
		c.Url = getUserResp.Url
	}
}

func (s *FileDomainService) LoadLikeCount(ctx context.Context, c *core_api.FileCount, id string) {
	getRelationCountResp, err := s.Platform.GetRelationCount(ctx, &platform.GetRelationCountReq{
		RelationFilterOptions: &platform.GetRelationCountReq_ToFilterOptions{
			ToFilterOptions: &platform.ToFilterOptions{
				ToType:   int64(core_api.TargetType_FileType),
				ToId:     id,
				FromType: int64(core_api.TargetType_UserType),
			},
		},
		RelationType: int64(core_api.RelationType_LikeRelationType),
	})
	if err == nil {
		c.LikedCount = getRelationCountResp.Total
	}
}

func (s *FileDomainService) LoadViewCount(ctx context.Context, c *core_api.FileCount, id string) {
	getRelationCountResp, err := s.Platform.GetRelationCount(ctx, &platform.GetRelationCountReq{
		RelationFilterOptions: &platform.GetRelationCountReq_ToFilterOptions{
			ToFilterOptions: &platform.ToFilterOptions{
				ToType:   int64(core_api.TargetType_FileType),
				ToId:     id,
				FromType: int64(core_api.TargetType_UserType),
			},
		},
		RelationType: int64(core_api.RelationType_ViewRelationType),
	})
	if err == nil {
		c.ViewCount = getRelationCountResp.Total
	}
}

func (s *FileDomainService) LoadLiked(ctx context.Context, c *core_api.FileRelation, id, userId string) {
	getRelationResp, err := s.Platform.GetRelation(ctx, &platform.GetRelationReq{
		FromType:     int64(core_api.TargetType_UserType),
		FromId:       userId,
		ToType:       int64(core_api.TargetType_FileType),
		ToId:         id,
		RelationType: int64(core_api.RelationType_LikeRelationType),
	})
	if err == nil {
		c.Liked = getRelationResp.Ok
	}
}

func (s *FileDomainService) LoadCollected(ctx context.Context, c *core_api.FileRelation, id, userId string) {
	getRelationResp, err := s.Platform.GetRelation(ctx, &platform.GetRelationReq{
		FromType:     int64(core_api.TargetType_UserType),
		FromId:       userId,
		ToType:       int64(core_api.TargetType_FileType),
		ToId:         id,
		RelationType: int64(core_api.RelationType_CollectRelationType),
	})
	if err == nil {
		c.Collected = getRelationResp.Ok
	}
}

func (s *FileDomainService) LoadCollectCount(ctx context.Context, c *core_api.FileCount, id string) {
	getRelationCountResp, err := s.Platform.GetRelationCount(ctx, &platform.GetRelationCountReq{
		RelationFilterOptions: &platform.GetRelationCountReq_ToFilterOptions{
			ToFilterOptions: &platform.ToFilterOptions{
				ToType:   int64(core_api.TargetType_FileType),
				ToId:     id,
				FromType: int64(core_api.TargetType_UserType),
			},
		},
		RelationType: int64(core_api.RelationType_CollectRelationType),
	})
	if err == nil {
		c.CollectCount = getRelationCountResp.Total
	}
}

func (s *FileDomainService) LoadLabels(ctx context.Context, c *[]*core_api.Label, labelIds []string) {
	var labels *platform.GetLabelsInBatchResp
	labels, _ = s.Platform.GetLabelsInBatch(ctx, &platform.GetLabelsInBatchReq{Ids: labelIds})
	*c = lo.Map(labels.Labels, func(item string, index int) *core_api.Label {
		return &core_api.Label{Id: labelIds[index], Value: item}
	})
}
