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
	"github.com/zeromicro/go-zero/core/mr"
)

type IPostDomainService interface {
	LoadAuthor(ctx context.Context, user *core_api.PostUser, userId string)
	LoadLikeCount(ctx context.Context, likeCount *int64, postId string)
	LoadViewCount(ctx context.Context, viewCount *int64, postId string)
	LoadCollectCount(ctx context.Context, collectCount *int64, postId string)
	LoadCommentCount(ctx context.Context, commentCount *int64, postId string)
	LoadShareCount(ctx context.Context, shareCount *int64, postId string)
	LoadLiked(ctx context.Context, liked *bool, userId, postId string)
	LoadCollected(ctx context.Context, collected *bool, userId, postId string)
	LoadLabels(ctx context.Context, c *[]*core_api.Label, labels []string)
}
type PostDomainService struct {
	CloudMindContent cloudmind_content.ICloudMindContent
	Platform         platformservice.IPlatForm
}

var PostDomainServiceSet = wire.NewSet(
	wire.Struct(new(PostDomainService), "*"),
	wire.Bind(new(IPostDomainService), new(*PostDomainService)),
)

func (s *PostDomainService) LoadLabels(ctx context.Context, c *[]*core_api.Label, labels []string) {
	getLabelsInBatchResp, err := s.Platform.GetLabelsInBatch(ctx, &platform.GetLabelsInBatchReq{
		Ids: labels,
	})
	if err == nil {
		*c = lo.Map(getLabelsInBatchResp.Labels, func(value string, i int) *core_api.Label {
			return &core_api.Label{
				Id:    labels[i],
				Value: value,
			}
		})
	}
}

func (s *PostDomainService) LoadCommentCount(ctx context.Context, commentCount *int64, postId string) {
	getRelationCountResp, err := s.Platform.GetCommentSubject(ctx, &platform.GetCommentSubjectReq{Id: postId})
	if err == nil {
		*commentCount = getRelationCountResp.AllCount
	}
}

func (s *PostDomainService) LoadShareCount(ctx context.Context, shareCount *int64, postId string) {
	getRelationCountResp, err := s.Platform.GetRelationCount(ctx, &platform.GetRelationCountReq{
		RelationFilterOptions: &platform.GetRelationCountReq_ToFilterOptions{
			ToFilterOptions: &platform.ToFilterOptions{
				ToType:   int64(core_api.TargetType_PostType),
				ToId:     postId,
				FromType: int64(core_api.TargetType_UserType),
			},
		},
		RelationType: int64(core_api.RelationType_ShareRelationType),
	})
	if err == nil {
		*shareCount = getRelationCountResp.Total
	}
}

func (s *PostDomainService) LoadAuthor(ctx context.Context, user *core_api.PostUser, userId string) {
	if userId == "" {
		return
	}
	_ = mr.Finish(func() error {
		getUserResp, err := s.CloudMindContent.GetUser(ctx, &content.GetUserReq{UserId: user.UserId})
		if err == nil {
			user.Name = getUserResp.Name
			user.Url = getUserResp.Url
			user.Labels = getUserResp.Labels
			user.Description = getUserResp.Description
		}
		return nil
	}, func() error {
		getRelationCountResp, err := s.Platform.GetRelationCount(ctx, &platform.GetRelationCountReq{
			RelationFilterOptions: &platform.GetRelationCountReq_ToFilterOptions{
				ToFilterOptions: &platform.ToFilterOptions{
					ToType:   int64(core_api.TargetType_UserType),
					ToId:     userId,
					FromType: int64(core_api.TargetType_UserType),
				},
			},
			RelationType: int64(core_api.RelationType_FollowRelationType),
		})
		if err == nil {
			user.FollowedCount = getRelationCountResp.Total
		}
		return nil
	}, func() error {
		getRelationPathsCountResp, err := s.Platform.GetRelationPathsCount(ctx, &platform.GetRelationPathsCountReq{
			FromType1: int64(core_api.TargetType_UserType),
			FromId1:   userId,
			FromType2: int64(core_api.TargetType_UserType),
			EdgeType1: int64(core_api.RelationType_PublishRelationType),
			EdgeType2: int64(core_api.RelationType_LikeRelationType),
			ToType:    int64(core_api.TargetType_PostType),
		})
		if err == nil {
			user.LikedCount = getRelationPathsCountResp.Total
		}
		return nil
	}, func() error {
		if userId == "" {
			return nil
		}
		getRelation, err := s.Platform.GetRelation(ctx, &platform.GetRelationReq{
			FromType:     int64(core_api.TargetType_UserType),
			FromId:       userId,
			ToType:       int64(core_api.TargetType_UserType),
			ToId:         user.UserId,
			RelationType: int64(core_api.RelationType_FollowRelationType),
		})
		if err != nil {
			return err
		}
		user.Followed = getRelation.Ok
		return nil
	}, func() error {
		getRelationCountResp, err := s.Platform.GetRelationCount(ctx, &platform.GetRelationCountReq{
			RelationFilterOptions: &platform.GetRelationCountReq_FromFilterOptions{
				FromFilterOptions: &platform.FromFilterOptions{
					ToType:   int64(core_api.TargetType_PostType),
					FromId:   userId,
					FromType: int64(core_api.TargetType_UserType),
				},
			},
			RelationType: int64(core_api.RelationType_PublishRelationType),
		})
		if err == nil {
			user.PostCount = getRelationCountResp.Total
		}
		return nil
	})
}

func (s *PostDomainService) LoadLikeCount(ctx context.Context, likeCount *int64, postId string) {
	getRelationCountResp, err := s.Platform.GetRelationCount(ctx, &platform.GetRelationCountReq{
		RelationFilterOptions: &platform.GetRelationCountReq_ToFilterOptions{
			ToFilterOptions: &platform.ToFilterOptions{
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
	getRelationCountResp, err := s.Platform.GetRelationCount(ctx, &platform.GetRelationCountReq{
		RelationFilterOptions: &platform.GetRelationCountReq_ToFilterOptions{
			ToFilterOptions: &platform.ToFilterOptions{
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
	getRelationCountResp, err := s.Platform.GetRelationCount(ctx, &platform.GetRelationCountReq{
		RelationFilterOptions: &platform.GetRelationCountReq_ToFilterOptions{
			ToFilterOptions: &platform.ToFilterOptions{
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
	getRelationResp, err := s.Platform.GetRelation(ctx, &platform.GetRelationReq{
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
	getRelationResp, err := s.Platform.GetRelation(ctx, &platform.GetRelationReq{
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
