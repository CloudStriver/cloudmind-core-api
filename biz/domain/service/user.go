package service

import (
	"context"
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	platformservice "github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/platform"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/platform"
	"github.com/google/wire"
	"github.com/samber/lo"
)

type IUserDomainService interface {
	LoadFollowedCount(ctx context.Context, followedCount *int64, userId string)
	LoadFollowCount(ctx context.Context, followCount *int64, userId string)
	LoadLabel(ctx context.Context, labels []string)
	LoadFollowed(ctx context.Context, followed *bool, fromUserId string, toUserId string)
}
type UserDomainService struct {
	Config   *config.Config
	Platform platformservice.IPlatForm
}

var UserDomainServiceSet = wire.NewSet(
	wire.Struct(new(UserDomainService), "*"),
	wire.Bind(new(IUserDomainService), new(*UserDomainService)),
)

func (s *UserDomainService) LoadFollowed(ctx context.Context, followed *bool, fromUserId string, toUserId string) {
	getRelationResp, err := s.Platform.GetRelation(ctx, &platform.GetRelationReq{
		FromType:     int64(core_api.TargetType_UserType),
		FromId:       fromUserId,
		ToType:       int64(core_api.TargetType_UserType),
		ToId:         toUserId,
		RelationType: int64(core_api.RelationType_FollowRelationType),
	})
	if err == nil {
		*followed = getRelationResp.Ok
	}
}
func (s *UserDomainService) LoadLabel(ctx context.Context, labels []string) {
	getLabelsResp, err := s.Platform.GetLabelsInBatch(ctx, &platform.GetLabelsInBatchReq{
		LabelIds: labels,
	})
	if err == nil {
		lo.ForEach(getLabelsResp.Labels, func(label string, i int) {
			labels[i] = label
		})
	}
}

func (s *UserDomainService) LoadFollowCount(ctx context.Context, followCount *int64, userId string) {
	getRelationCountResp, err := s.Platform.GetRelationCount(ctx, &platform.GetRelationCountReq{
		RelationFilterOptions: &platform.GetRelationCountReq_FromFilterOptions{
			FromFilterOptions: &platform.FromFilterOptions{
				ToType:   int64(core_api.TargetType_UserType),
				FromId:   userId,
				FromType: int64(core_api.TargetType_UserType),
			},
		},
		RelationType: int64(core_api.RelationType_FollowRelationType),
	})
	if err == nil {
		*followCount = getRelationCountResp.Total
	}
}

func (s *UserDomainService) LoadFollowedCount(ctx context.Context, followedCount *int64, userId string) {
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
		*followedCount = getRelationCountResp.Total
	}
}
