package service

import (
	"context"
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_content"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/platform_comment"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/platform_relation"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/platform/comment"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/platform/relation"
	"github.com/google/wire"
	"github.com/samber/lo"
)

type IPostDomainService interface {
	LoadAuthor(ctx context.Context, user *core_api.User, userId string)
	LoadLikeCount(ctx context.Context, likeCount *int64, postId string)
	LoadViewCount(ctx context.Context, viewCount *int64, postId string)
	LoadCollectCount(ctx context.Context, collectCount *int64, postId string)
	LoadLiked(ctx context.Context, liked *bool, userId, postId string)
	LoadCollected(ctx context.Context, collected *bool, userId, postId string)
	LoadLabels(ctx context.Context, labels []string)
}
type PostDomainService struct {
	CloudMindContent cloudmind_content.ICloudMindContent
	PlatformRelation platform_relation.IPlatFormRelation
	PlatFormComment  platform_comment.IPlatFormComment
}

var PostDomainServiceSet = wire.NewSet(
	wire.Struct(new(PostDomainService), "*"),
	wire.Bind(new(IPostDomainService), new(*PostDomainService)),
)

func (s *PostDomainService) LoadLabels(ctx context.Context, labels []string) {
	getLabelsInBatchResp, err := s.PlatFormComment.GetLabelsInBatch(ctx, &comment.GetLabelsInBatchReq{
		LabelIds: labels,
	})
	if err == nil {
		lo.ForEach(getLabelsInBatchResp.Labels, func(label string, i int) {
			labels[i] = label
		})
	}
}
func (s *PostDomainService) LoadAuthor(ctx context.Context, user *core_api.User, userId string) {
	if userId == "" {
		return
	}
	getUserResp, err := s.CloudMindContent.GetUser(ctx, &content.GetUserReq{UserId: userId})
	if err == nil {
		user.Name = getUserResp.Name
		user.Url = getUserResp.Url
		user.UserId = userId
	}
}

func (s *PostDomainService) LoadLikeCount(ctx context.Context, likeCount *int64, postId string) {
	getRelationCountResp, err := s.PlatformRelation.GetRelationCount(ctx, &relation.GetRelationCountReq{
		RelationFilterOptions: &relation.GetRelationCountReq_ToFilterOptions{
			ToFilterOptions: &relation.ToFilterOptions{
				ToType:   int64(core_api.TargetType_PostType),
				ToId:     postId,
				FromType: int64(core_api.TargetType_UserType),
			},
		},
		RelationType: int64(core_api.RelationType_LikeRelationType),
	})
	if err == nil {
		*likeCount = getRelationCountResp.Total
	}
}

func (s *PostDomainService) LoadViewCount(ctx context.Context, viewCount *int64, postId string) {
	getRelationCountResp, err := s.PlatformRelation.GetRelationCount(ctx, &relation.GetRelationCountReq{
		RelationFilterOptions: &relation.GetRelationCountReq_ToFilterOptions{
			ToFilterOptions: &relation.ToFilterOptions{
				ToType:   int64(core_api.TargetType_PostType),
				ToId:     postId,
				FromType: int64(core_api.TargetType_UserType),
			},
		},
		RelationType: int64(core_api.RelationType_ViewRelationType),
	})
	if err == nil {
		*viewCount = getRelationCountResp.Total
	}
}

func (s *PostDomainService) LoadCollectCount(ctx context.Context, collectCount *int64, postId string) {
	getRelationCountResp, err := s.PlatformRelation.GetRelationCount(ctx, &relation.GetRelationCountReq{
		RelationFilterOptions: &relation.GetRelationCountReq_ToFilterOptions{
			ToFilterOptions: &relation.ToFilterOptions{
				ToType:   int64(core_api.TargetType_PostType),
				ToId:     postId,
				FromType: int64(core_api.TargetType_UserType),
			},
		},
		RelationType: int64(core_api.RelationType_CollectRelationType),
	})
	if err == nil {
		*collectCount = getRelationCountResp.Total
	}
}

func (s *PostDomainService) LoadLiked(ctx context.Context, liked *bool, userId, postId string) {
	getRelationResp, err := s.PlatformRelation.GetRelation(ctx, &relation.GetRelationReq{
		FromType:     int64(core_api.TargetType_UserType),
		FromId:       userId,
		ToType:       int64(core_api.TargetType_PostType),
		ToId:         postId,
		RelationType: int64(core_api.RelationType_LikeRelationType),
	})
	if err == nil {
		*liked = getRelationResp.Ok
	}
}

func (s *PostDomainService) LoadCollected(ctx context.Context, collected *bool, userId, postId string) {
	getRelationResp, err := s.PlatformRelation.GetRelation(ctx, &relation.GetRelationReq{
		FromType:     int64(core_api.TargetType_UserType),
		FromId:       userId,
		ToType:       int64(core_api.TargetType_PostType),
		ToId:         postId,
		RelationType: int64(core_api.RelationType_CollectRelationType),
	})
	if err == nil {
		*collected = getRelationResp.Ok
	}
}
