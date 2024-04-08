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

type ICommentDomainService interface {
	LoadAuthor(ctx context.Context, c *core_api.CommentInfo, userId string)
	LoadLikeCount(ctx context.Context, c *core_api.CommentInfo)
	LoadLiked(ctx context.Context, c *core_api.CommentInfo, userId string)
	LoadHated(ctx context.Context, c *core_api.CommentInfo, userId string)
	LoadLabels(ctx context.Context, c *core_api.CommentInfo, labelIds []string)
}
type CommentDomainService struct {
	CloudMindUser cloudmind_content.ICloudMindContent
	Platform      platformservice.IPlatForm
}

var CommentDomainServiceSet = wire.NewSet(
	wire.Struct(new(CommentDomainService), "*"),
	wire.Bind(new(ICommentDomainService), new(*CommentDomainService)),
)

func (s *CommentDomainService) LoadAuthor(ctx context.Context, c *core_api.CommentInfo, userId string) {
	if userId == "" {
		return
	}
	c.Author = &core_api.User{
		UserId: userId,
	}
	getUserResp, err := s.CloudMindUser.GetUser(ctx, &content.GetUserReq{UserId: userId})
	if err == nil {
		c.Author.Name = getUserResp.Name
		c.Author.Url = getUserResp.Url
	}
}

func (s *CommentDomainService) LoadLikeCount(ctx context.Context, c *core_api.CommentInfo) {
	getRelationCountResp, err := s.Platform.GetRelationCount(ctx, &platform.GetRelationCountReq{
		RelationFilterOptions: &platform.GetRelationCountReq_ToFilterOptions{
			ToFilterOptions: &platform.ToFilterOptions{
				ToType:   int64(core_api.TargetType_CommentContentType),
				ToId:     c.Id,
				FromType: int64(core_api.TargetType_UserType),
			},
		},
		RelationType: int64(core_api.RelationType_LikeRelationType),
	})
	if err == nil {
		c.Like = getRelationCountResp.Total
	}
}

func (s *CommentDomainService) LoadLiked(ctx context.Context, c *core_api.CommentInfo, userId string) {
	getRelationResp, err := s.Platform.GetRelation(ctx, &platform.GetRelationReq{
		FromType:     int64(core_api.TargetType_UserType),
		FromId:       userId,
		ToType:       int64(core_api.TargetType_CommentContentType),
		ToId:         c.Id,
		RelationType: int64(core_api.RelationType_LikeRelationType),
	})
	if err == nil {
		c.CommentRelation.Liked = getRelationResp.Ok
	}
}

func (s *CommentDomainService) LoadHated(ctx context.Context, c *core_api.CommentInfo, userId string) {
	getRelationResp, err := s.Platform.GetRelation(ctx, &platform.GetRelationReq{
		FromType:     int64(core_api.TargetType_UserType),
		FromId:       userId,
		ToType:       int64(core_api.TargetType_CommentContentType),
		ToId:         c.Id,
		RelationType: int64(core_api.RelationType_HateRelationType),
	})
	if err == nil {
		c.CommentRelation.Hated = getRelationResp.Ok
	}
}

func (s *CommentDomainService) LoadLabels(ctx context.Context, c *core_api.CommentInfo, labelIds []string) {
	var labels *platform.GetLabelsInBatchResp
	labels, _ = s.Platform.GetLabelsInBatch(ctx, &platform.GetLabelsInBatchReq{Ids: labelIds})
	c.Labels = lo.Map(labels.Labels, func(item string, index int) *core_api.Label {
		return &core_api.Label{LabelId: labelIds[index], Value: item}
	})
}
