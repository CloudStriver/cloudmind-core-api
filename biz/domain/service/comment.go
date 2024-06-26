package service

import (
	"context"
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_content"
	platformservice "github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/platform"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/platform"
	"github.com/google/wire"
)

type ICommentDomainService interface {
	LoadAuthor(ctx context.Context, c *core_api.SimpleUser, userId string)
	LoadFile(ctx context.Context, description, userId *string, id string)
	LoadPost(ctx context.Context, title, userId *string, id string)
	LoadLikeCount(ctx context.Context, c *int64, id string)
	LoadHateCount(ctx context.Context, c *int64, id string)
	LoadLiked(ctx context.Context, c *core_api.CommentRelation, id, userId string)
	LoadHated(ctx context.Context, c *core_api.CommentRelation, id, userId string)
	LoadLabels(ctx context.Context, c *[]string, labelIds []string)
}
type CommentDomainService struct {
	CloudMindUser cloudmind_content.ICloudMindContent
	Platform      platformservice.IPlatForm
}

var CommentDomainServiceSet = wire.NewSet(
	wire.Struct(new(CommentDomainService), "*"),
	wire.Bind(new(ICommentDomainService), new(*CommentDomainService)),
)

func (s *CommentDomainService) LoadAuthor(ctx context.Context, c *core_api.SimpleUser, userId string) {
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

func (s *CommentDomainService) LoadFile(ctx context.Context, description, userId *string, id string) {
	getRes, err := s.CloudMindUser.GetPublicFile(ctx, &content.GetPublicFileReq{FileId: id})
	if err == nil {
		*description = getRes.Description
		*userId = getRes.UserId
	}
}

func (s *CommentDomainService) LoadPost(ctx context.Context, title, userId *string, id string) {
	getRes, err := s.CloudMindUser.GetPost(ctx, &content.GetPostReq{PostId: id})
	if err == nil {
		*title = getRes.Title
		*userId = getRes.UserId
	}
}

func (s *CommentDomainService) LoadLikeCount(ctx context.Context, c *int64, id string) {
	getRelationCountResp, err := s.Platform.GetRelationCount(ctx, &platform.GetRelationCountReq{
		RelationFilterOptions: &platform.GetRelationCountReq_ToFilterOptions{
			ToFilterOptions: &platform.ToFilterOptions{
				ToType:   int64(core_api.TargetType_CommentContentType),
				ToId:     id,
				FromType: int64(core_api.TargetType_UserType),
			},
		},
		RelationType: int64(core_api.RelationType_LikeRelationType),
	})
	if err == nil {
		*c = getRelationCountResp.Total
	}
}

func (s *CommentDomainService) LoadHateCount(ctx context.Context, c *int64, id string) {
	getRelationCountResp, err := s.Platform.GetRelationCount(ctx, &platform.GetRelationCountReq{
		RelationFilterOptions: &platform.GetRelationCountReq_ToFilterOptions{
			ToFilterOptions: &platform.ToFilterOptions{
				ToType:   int64(core_api.TargetType_CommentContentType),
				ToId:     id,
				FromType: int64(core_api.TargetType_UserType),
			},
		},
		RelationType: int64(core_api.RelationType_HateRelationType),
	})
	if err == nil {
		*c = getRelationCountResp.Total
	}
}

func (s *CommentDomainService) LoadLiked(ctx context.Context, c *core_api.CommentRelation, id, userId string) {
	getRelationResp, err := s.Platform.GetRelation(ctx, &platform.GetRelationReq{
		FromType:     int64(core_api.TargetType_UserType),
		FromId:       userId,
		ToType:       int64(core_api.TargetType_CommentContentType),
		ToId:         id,
		RelationType: int64(core_api.RelationType_LikeRelationType),
	})
	if err == nil {
		c.Liked = getRelationResp.Ok
	}
}

func (s *CommentDomainService) LoadHated(ctx context.Context, c *core_api.CommentRelation, id, userId string) {
	getRelationResp, err := s.Platform.GetRelation(ctx, &platform.GetRelationReq{
		FromType:     int64(core_api.TargetType_UserType),
		FromId:       userId,
		ToType:       int64(core_api.TargetType_CommentContentType),
		ToId:         id,
		RelationType: int64(core_api.RelationType_HateRelationType),
	})
	if err == nil {
		c.Hated = getRelationResp.Ok
	}
}

func (s *CommentDomainService) LoadLabels(ctx context.Context, c *[]string, labelIds []string) {
	var labels *platform.GetLabelsInBatchResp
	labels, _ = s.Platform.GetLabelsInBatch(ctx, &platform.GetLabelsInBatchReq{LabelIds: labelIds})
	*c = labels.Labels
}
