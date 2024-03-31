package service

import (
	"context"
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_content"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_trade"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/platform_relation"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/trade"
	"github.com/google/wire"
)

type IProductDomainService interface {
	LoadAuthor(ctx context.Context, user *core_api.User, userId string)
	LoadLikeCount(ctx context.Context, likeCount *int64, productId string)
	LoadViewCount(ctx context.Context, viewCount *int64, productId string)
	LoadCollectCount(ctx context.Context, collectCount *int64, productId string)
	LoadLiked(ctx context.Context, liked *bool, userId, productId string)
	LoadCollected(ctx context.Context, collected *bool, userId, productId string)
	LoadPurchaseCount(ctx context.Context, purchaseCount *int64, productId string)
	LoadStock(ctx context.Context, stock *int64, productId string)
}
type ProductDomainService struct {
	CloudMindUser    cloudmind_content.ICloudMindContent
	PlatformRelation platform_relation.IPlatFormRelation
	CloudMindTrade   cloudmind_trade.ICloudMindTrade
}

var ProductDomainServiceSet = wire.NewSet(
	wire.Struct(new(ProductDomainService), "*"),
	wire.Bind(new(IProductDomainService), new(*ProductDomainService)),
)

func (s *ProductDomainService) LoadStock(ctx context.Context, stock *int64, productId string) {
	getStockResp, err := s.CloudMindTrade.GetStock(ctx, &trade.GetStockReq{ProductId: productId})
	if err == nil {
		*stock = getStockResp.Stock
	}
}
func (s *ProductDomainService) LoadPurchaseCount(ctx context.Context, purchaseCount *int64, productId string) {
	//getRelationCountResp, err := s.PlatformRelation.GetRelationCount(ctx, &relation.GetRelationCountReq{
	//	RelationFilterOptions: &relation.GetRelationCountReq_ToFilterOptions{
	//		ToFilterOptions: &relation.ToFilterOptions{
	//			ToType:   int64(core_api.TargetType_ProductType),
	//			ToId:     productId,
	//			FromType: int64(core_api.TargetType_UserType),
	//		},
	//	},
	//	RelationType: int64(core_api.RelationType_PurchaseType),
	//})
	//if err == nil {
	//	*purchaseCount = getRelationCountResp.Total
	//}
}

func (s *ProductDomainService) LoadAuthor(ctx context.Context, user *core_api.User, userId string) {
	if userId == "" {
		return
	}
	getUserResp, err := s.CloudMindUser.GetUser(ctx, &content.GetUserReq{UserId: userId})
	if err == nil {
		user.Name = getUserResp.Name
		user.Url = getUserResp.Url
		user.UserId = userId
	}
}

func (s *ProductDomainService) LoadLikeCount(ctx context.Context, likeCount *int64, productId string) {
	//getRelationCountResp, err := s.PlatformRelation.GetRelationCount(ctx, &relation.GetRelationCountReq{
	//	RelationFilterOptions: &relation.GetRelationCountReq_ToFilterOptions{
	//		ToFilterOptions: &relation.ToFilterOptions{
	//			ToType:   int64(core_api.TargetType_ProductType),
	//			ToId:     productId,
	//			FromType: int64(core_api.TargetType_UserType),
	//		},
	//	},
	//	RelationType: int64(core_api.RelationType_LikeRelationType),
	//})
	//if err == nil {
	//	*likeCount = getRelationCountResp.Total
	//}
}

func (s *ProductDomainService) LoadViewCount(ctx context.Context, viewCount *int64, productId string) {
	//getRelationCountResp, err := s.PlatformRelation.GetRelationCount(ctx, &relation.GetRelationCountReq{
	//	RelationFilterOptions: &relation.GetRelationCountReq_ToFilterOptions{
	//		ToFilterOptions: &relation.ToFilterOptions{
	//			ToType:   int64(core_api.TargetType_ProductType),
	//			ToId:     productId,
	//			FromType: int64(core_api.TargetType_UserType),
	//		},
	//	},
	//	RelationType: int64(core_api.RelationType_ViewRelationType),
	//})
	//if err == nil {
	//	*viewCount = getRelationCountResp.Total
	//}
}

func (s *ProductDomainService) LoadCollectCount(ctx context.Context, collectCount *int64, productId string) {
	//getRelationCountResp, err := s.PlatformRelation.GetRelationCount(ctx, &relation.GetRelationCountReq{
	//	RelationFilterOptions: &relation.GetRelationCountReq_ToFilterOptions{
	//		ToFilterOptions: &relation.ToFilterOptions{
	//			ToType:   int64(core_api.TargetType_ProductType),
	//			ToId:     productId,
	//			FromType: int64(core_api.TargetType_UserType),
	//		},
	//	},
	//	RelationType: int64(core_api.RelationType_CollectRelationType),
	//})
	//if err == nil {
	//	*collectCount = getRelationCountResp.Total
	//}
}

func (s *ProductDomainService) LoadLiked(ctx context.Context, liked *bool, userId, productId string) {
	//getRelationResp, err := s.PlatformRelation.GetRelation(ctx, &relation.GetRelationReq{
	//	FromType:     int64(core_api.TargetType_UserType),
	//	FromId:       userId,
	//	ToType:       int64(core_api.TargetType_ProductType),
	//	ToId:         productId,
	//	RelationType: int64(core_api.RelationType_LikeRelationType),
	//})
	//if err == nil {
	//	*liked = getRelationResp.Ok
	//}
}

func (s *ProductDomainService) LoadCollected(ctx context.Context, collected *bool, userId, productId string) {
	//getRelationResp, err := s.PlatformRelation.GetRelation(ctx, &relation.GetRelationReq{
	//	FromType:     int64(core_api.TargetType_UserType),
	//	FromId:       userId,
	//	ToType:       int64(core_api.TargetType_ProductType),
	//	ToId:         productId,
	//	RelationType: int64(core_api.RelationType_CollectRelationType),
	//})
	//if err == nil {
	//	*collected = getRelationResp.Ok
	//}
}
