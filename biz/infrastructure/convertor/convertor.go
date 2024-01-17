package convertor

import (
	dto_basic "github.com/CloudStriver/cloudmind-core-api/biz/application/dto/basic"
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/basic"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/platform/relation"
)

func UserToUserDetailInfo(req *core_api.UserDetail) *content.UserDetailInfo {
	return &content.UserDetailInfo{
		Name:        req.Name,
		Sex:         int64(req.Sex),
		FullName:    req.FullName,
		IdCard:      req.IdCard,
		Description: req.Description,
		Url:         req.Url,
	}
}

func UserDetailToUser(req *content.User) *core_api.User {
	return &core_api.User{
		UserId: req.UserId,
		Name:   req.Name,
		Url:    req.Url,
	}
}

func PaginationOptionsToPaginationOptions(req *dto_basic.PaginationOptions) *basic.PaginationOptions {
	return &basic.PaginationOptions{
		Limit:     req.Limit,
		LastToken: req.LastToken,
		Backward:  req.Backward,
		Offset:    req.Offset,
	}
}

func CoreApiRelationInfoToRelationInfo(req *core_api.RelationInfo) *relation.RelationInfo {
	return &relation.RelationInfo{
		FromType:     req.FromType,
		FromId:       req.FromId,
		ToType:       req.ToType,
		ToId:         req.ToId,
		RelationType: int64(req.RelationType),
	}
}
