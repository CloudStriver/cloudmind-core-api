package service

import (
	"context"
	"fmt"
	"github.com/CloudStriver/cloudmind-core-api/biz/adaptor"
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/cloudmind-core-api/biz/domain/service"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/consts"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_content"
	platformservice "github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/platform"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/basic"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/platform"
	"github.com/google/wire"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/mr"
)

type IRelationService interface {
	GetRelation(ctx context.Context, req *core_api.GetRelationReq) (resp *core_api.GetRelationResp, err error)
	CreateRelation(ctx context.Context, req *core_api.CreateRelationReq) (resp *core_api.CreateRelationResp, err error)
	GetFromRelations(ctx context.Context, c *core_api.GetFromRelationsReq) (resp *core_api.GetFromRelationsResp, err error)
	GetToRelations(ctx context.Context, c *core_api.GetToRelationsReq) (resp *core_api.GetToRelationsResp, err error)
	DeleteRelation(ctx context.Context, c *core_api.DeleteRelationReq) (resp *core_api.DeleteRelationResp, err error)
}

var RelationServiceSet = wire.NewSet(
	wire.Struct(new(RelationService), "*"),
	wire.Bind(new(IRelationService), new(*RelationService)),
)

type RelationService struct {
	Config                *config.Config
	CloudMindContent      cloudmind_content.ICloudMindContent
	Platform              platformservice.IPlatForm
	UserDomainService     service.IUserDomainService
	PostDomainService     service.IPostDomainService
	RelationDomainService service.RelationDomainService
}

func (s *RelationService) GetFromRelations(ctx context.Context, req *core_api.GetFromRelationsReq) (resp *core_api.GetFromRelationsResp, err error) {
	resp = new(core_api.GetFromRelationsResp)
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil {
		return resp, consts.ErrNotAuthentication
	}
	getFromRelationsResp, err := s.Platform.GetRelations(ctx, &platform.GetRelationsReq{
		RelationFilterOptions: &platform.GetRelationsReq_FromFilterOptions{
			FromFilterOptions: &platform.FromFilterOptions{
				FromType: req.FromType,
				FromId:   req.FromId,
				ToType:   int64(req.ToType),
			},
		},
		RelationType: req.RelationType,
		PaginationOptions: &basic.PaginationOptions{
			Limit:  req.Limit,
			Offset: req.Offset,
		},
	})
	if err != nil {
		return resp, err
	}
	fmt.Println(getFromRelationsResp.Relations)

	switch req.ToType {
	case core_api.TargetType_UserType:
		resp.Users = make([]*core_api.User, len(getFromRelationsResp.Relations))
		err = mr.Finish(lo.Map[*platform.Relation](getFromRelationsResp.Relations, func(platform *platform.Relation, i int) func() error {
			return func() error {
				resp.Users[i] = &core_api.User{
					UserId: platform.ToId,
				}
				if err = mr.Finish(func() error {
					user, err := s.CloudMindContent.GetUser(ctx, &content.GetUserReq{
						UserId: platform.ToId,
					})
					if err != nil {
						return err
					}
					resp.Users[i].Name = user.Name
					resp.Users[i].Url = user.Url
					resp.Users[i].Tags = user.Labels
					s.UserDomainService.LoadLabel(ctx, resp.Users[i].Tags)
					return nil
				}, func() error {
					if userData.GetUserId() != "" && userData.UserId != resp.Users[i].UserId {
						s.UserDomainService.LoadFollowed(ctx, &resp.Users[i].Followed, userData.UserId, resp.Users[i].UserId)
					}
					return nil
				}, func() error {
					s.UserDomainService.LoadFollowedCount(ctx, &resp.Users[i].FollowedCount, resp.Users[i].UserId)
					return nil
				}); err != nil {
					return err
				}
				return nil
			}
		})...)
		if err != nil {
			return resp, err
		}
	case core_api.TargetType_PostType:
		resp.Posts = make([]*core_api.Post, len(getFromRelationsResp.Relations))
		if err = mr.Finish(lo.Map[*platform.Relation](getFromRelationsResp.Relations, func(val *platform.Relation, i int) func() error {
			return func() error {
				resp.Posts[i] = &core_api.Post{}
				if err = mr.Finish(func() error {
					post, err1 := s.CloudMindContent.GetPost(ctx, &content.GetPostReq{
						PostId: val.ToId,
					})
					if err1 != nil {
						return err1
					}

					tags := lo.Map[*content.Tag, *core_api.TagInfo](post.Tags, func(item *content.Tag, index int) *core_api.TagInfo {
						return &core_api.TagInfo{
							TagId:  item.TagId,
							ZoneId: item.ZoneId,
						}
					})
					tagsId := lo.Map[*content.Tag, string](post.Tags, func(item *content.Tag, index int) string {
						return item.TagId
					})

					resp.Posts[i].PostId = val.ToId
					resp.Posts[i].Title = post.Title
					resp.Posts[i].Text = post.Text
					resp.Posts[i].Url = post.Url
					resp.Posts[i].Tags = tags
					s.PostDomainService.LoadLabels(ctx, tagsId)
					for i := range tags {
						tags[i].Value = tagsId[i]
					}
					user, err1 := s.CloudMindContent.GetUser(ctx, &content.GetUserReq{
						UserId: post.UserId,
					})
					if err1 != nil {
						return err1
					}
					resp.Posts[i].UserName = user.Name
					return nil
				}, func() error {
					s.PostDomainService.LoadLikeCount(ctx, &resp.Posts[i].LikeCount, val.ToId)
					return nil
				}, func() error {
					if userData.GetUserId() != "" {
						s.PostDomainService.LoadLiked(ctx, &resp.Posts[i].Liked, userData.UserId, val.ToId)
					}
					return nil
				}, func() error {
					getCommentListResp, err2 := s.Platform.GetCommentList(ctx, &platform.GetCommentListReq{
						FilterOptions: &platform.CommentFilterOptions{
							OnlySubjectId: lo.ToPtr(val.ToId),
						},
						Pagination: &basic.PaginationOptions{
							Limit: lo.ToPtr(int64(1)),
						},
					})
					if err2 != nil {
						return err2
					}
					resp.Posts[i].CommentCount = getCommentListResp.Total
					return nil
				}); err != nil {
					return err
				}
				return nil
			}
		})...); err != nil {
			return resp, err
		}
	default:
		//return resp, consts.ErrNotSupportTargetType
	}
	return resp, nil
}

func (s *RelationService) GetToRelations(ctx context.Context, req *core_api.GetToRelationsReq) (resp *core_api.GetToRelationsResp, err error) {
	resp = new(core_api.GetToRelationsResp)
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil {
		return resp, consts.ErrNotAuthentication
	}
	getFromRelationsResp, err := s.Platform.GetRelations(ctx, &platform.GetRelationsReq{
		RelationFilterOptions: &platform.GetRelationsReq_ToFilterOptions{
			ToFilterOptions: &platform.ToFilterOptions{
				ToType:   req.ToType,
				ToId:     req.ToId,
				FromType: int64(req.FromType),
			},
		},
		RelationType: req.RelationType,
		PaginationOptions: &basic.PaginationOptions{
			Limit:  req.Limit,
			Offset: req.Offset,
		},
	})
	if err != nil {
		return resp, err
	}

	switch req.FromType {
	case core_api.TargetType_UserType:
		resp.Users = make([]*core_api.User, len(getFromRelationsResp.Relations))
		err = mr.Finish(lo.Map[*platform.Relation](getFromRelationsResp.Relations, func(platform *platform.Relation, i int) func() error {
			return func() error {
				resp.Users[i] = &core_api.User{
					UserId: platform.FromId,
				}
				if err = mr.Finish(func() error {
					user, err := s.CloudMindContent.GetUser(ctx, &content.GetUserReq{
						UserId: platform.FromId,
					})
					if err != nil {
						return err
					}
					resp.Users[i].Name = user.Name
					resp.Users[i].Url = user.Url
					resp.Users[i].Tags = user.Labels
					s.UserDomainService.LoadLabel(ctx, resp.Users[i].Tags)
					return nil
				}, func() error {
					if userData.GetUserId() != "" {
						s.UserDomainService.LoadFollowed(ctx, &resp.Users[i].Followed, userData.UserId, resp.Users[i].UserId)
					}
					return nil
				}, func() error {
					s.UserDomainService.LoadFollowedCount(ctx, &resp.Users[i].FollowedCount, resp.Users[i].UserId)
					return nil
				}); err != nil {
					return err
				}
				return nil
			}
		})...)
		if err != nil {
			return resp, err
		}
	default:
		//return resp, consts.ErrNotSupportTargetType
	}
	return resp, nil
}

func (s *RelationService) DeleteRelation(ctx context.Context, req *core_api.DeleteRelationReq) (resp *core_api.DeleteRelationResp, err error) {
	user, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || user.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	if _, err = s.Platform.DeleteRelation(ctx, &platform.DeleteRelationReq{
		FromType:     int64(core_api.TargetType_UserType),
		FromId:       user.UserId,
		ToType:       req.ToType,
		ToId:         req.ToId,
		RelationType: req.RelationType,
	}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *RelationService) GetRelation(ctx context.Context, req *core_api.GetRelationReq) (resp *core_api.GetRelationResp, err error) {
	resp = new(core_api.GetRelationResp)
	getRelationResp, err := s.Platform.GetRelation(ctx, &platform.GetRelationReq{
		FromType:     req.FromType,
		FromId:       req.FromId,
		ToType:       req.ToType,
		ToId:         req.ToId,
		RelationType: req.RelationType,
	})
	if err != nil {
		return resp, err
	}

	resp.Ok = getRelationResp.Ok
	return resp, nil
}

func (s *RelationService) CreateRelation(ctx context.Context, req *core_api.CreateRelationReq) (resp *core_api.CreateRelationResp, err error) {
	user, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || user.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	if err = s.RelationDomainService.CreateRelation(ctx, &core_api.Relation{
		FromType:     core_api.TargetType_UserType,
		FromId:       user.GetUserId(),
		ToType:       req.ToType,
		ToId:         req.ToId,
		RelationType: req.RelationType,
	}); err != nil {
		return resp, err
	}
	return resp, nil
}
