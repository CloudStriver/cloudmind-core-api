package service

import (
	"context"
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/platform_comment"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/platform_relation"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/platform/comment"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/platform/relation"
	"github.com/google/wire"
	"github.com/samber/lo"
)

type IUserDomainService interface {
	LoadFollowCount(ctx context.Context, followCount *int64, userId string)
	LoadLabel(ctx context.Context, labels []string)
}
type UserDomainService struct {
	Config           *config.Config
	PlatFormRelation platform_relation.IPlatFormRelation
	PlatFormComment  platform_comment.IPlatFormComment
}

var UserDomainServiceSet = wire.NewSet(
	wire.Struct(new(UserDomainService), "*"),
	wire.Bind(new(IUserDomainService), new(*UserDomainService)),
)

func (s *UserDomainService) LoadLabel(ctx context.Context, labels []string) {
	getLabelsResp, err := s.PlatFormComment.GetLabelsInBatch(ctx, &comment.GetLabelsInBatchReq{
		LabelIds: labels,
	})
	if err == nil {
		lo.ForEach(getLabelsResp.Labels, func(value string, i int) {
			labels[i] = value
		})
	}
}
func (s *UserDomainService) LoadFollowCount(ctx context.Context, followCount *int64, userId string) {
	getRelationCountResp, err := s.PlatFormRelation.GetRelationCount(ctx, &relation.GetRelationCountReq{
		RelationFilterOptions: &relation.GetRelationCountReq_ToFilterOptions{
			ToFilterOptions: &relation.ToFilterOptions{
				ToType:   int64(core_api.TargetType_UserType),
				ToId:     userId,
				FromType: int64(core_api.TargetType_UserType),
			},
		},
		RelationType: int64(core_api.RelationType_FollowType),
	})
	if err == nil {
		*followCount = getRelationCountResp.Total
	}
}
