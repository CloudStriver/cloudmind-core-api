package service

import (
	"context"
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/platform_relation"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/platform/relation"
	"github.com/google/wire"
)

type IUserDomainService interface {
	LoadFollowCount(ctx context.Context, followCount *int64, userId string)
}
type UserDomainService struct {
	Config           *config.Config
	PlatFormRelation platform_relation.IPlatFormRelation
}

var UserDomainServiceSet = wire.NewSet(
	wire.Struct(new(UserDomainService), "*"),
	wire.Bind(new(IUserDomainService), new(*UserDomainService)),
)

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
