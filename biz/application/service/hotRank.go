package service

import (
	"context"
	"github.com/CloudStriver/cloudmind-core-api/biz/adaptor"
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/cloudmind-core-api/biz/domain/service"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/consts"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_content"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_sts"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_system"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_trade"
	platformservice "github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/platform"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/basic"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/platform"
	"github.com/google/wire"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/mr"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type IHotRankService interface {
	GetHotRanks(ctx context.Context, req *core_api.GetHotRanksReq) (resp *core_api.GetHotRanksResp, err error)
}

var HotRankServiceSet = wire.NewSet(
	wire.Struct(new(HotRankService), "*"),
	wire.Bind(new(IHotRankService), new(*HotRankService)),
)

type HotRankService struct {
	Config            *config.Config
	CloudMindContent  cloudmind_content.ICloudMindContent
	CloudMindSts      cloudmind_sts.ICloudMindSts
	CloudMindTrade    cloudmind_trade.ICloudMindTrade
	CloudMindSystem   cloudmind_system.ICloudMindSystem
	Platform          platformservice.IPlatForm
	UserDomainService service.IUserDomainService
	Redis             *redis.Redis
}

func (s *HotRankService) GetHotRanks(ctx context.Context, req *core_api.GetHotRanksReq) (resp *core_api.GetHotRanksResp, err error) {
	resp = new(core_api.GetHotRanksResp)
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil {
		return resp, consts.ErrNotAuthentication
	}
	key := ""
	switch req.TargetType {
	case core_api.TargetType_UserType:
		key = consts.UserRankKey
	case core_api.TargetType_FileType:
		key = consts.FileRankKey
	case core_api.TargetType_PostType:
		key = consts.PostRankKey
	case core_api.TargetType_CommentContentType:
		key = consts.CommentRankKey
	}

	values, err := s.Redis.ZrevrangeCtx(ctx, key, req.Offset, req.Offset+req.Limit-1)
	if err != nil {
		return resp, err
	}
	if len(values) == 0 {
		return resp, nil
	}
	switch req.TargetType {
	case core_api.TargetType_UserType:
		users, err := s.CloudMindContent.GetUsersByUserIds(ctx, &content.GetUsersByUserIdsReq{
			UserIds: values,
		})
		if err != nil {
			return resp, err
		}
		resp.Users = make([]*core_api.HotUser, len(values))
		if err = mr.Finish(lo.Map(users.Users, func(item *content.User, i int) func() error {
			return func() error {
				resp.Users[i] = &core_api.HotUser{
					UserId:      item.UserId,
					Name:        item.Name,
					Url:         item.Url,
					Description: item.Description,
				}
				if userData.GetUserId() != "" || userData.UserId == item.UserId {
					s.UserDomainService.LoadFollowed(ctx, &resp.Users[i].Followed, userData.UserId, item.UserId)
				}
				return nil
			}
		})...); err != nil {
			return resp, err
		}
	case core_api.TargetType_FileType:
		files, err := s.CloudMindContent.GetFilesByIds(ctx, &content.GetFilesByIdsReq{
			FileIds: values,
		})
		if err != nil {
			return resp, err
		}
		resp.Files = make([]*core_api.HotFile, len(values))
		if err = mr.Finish(lo.Map(files.Files, func(item *content.File, i int) func() error {
			return func() error {
				resp.Files[i] = &core_api.HotFile{
					FileId: item.FileId,
					Name:   item.Name,
					Type:   item.Type,
				}

				return nil
			}
		})...); err != nil {
			return resp, err
		}
	case core_api.TargetType_PostType:
		posts, err := s.CloudMindContent.GetPostsByPostIds(ctx, &content.GetPostsByPostIdsReq{
			PostIds: values,
		})
		if err != nil {
			return resp, err
		}
		resp.Posts = make([]*core_api.HotPost, len(posts.Posts))
		if err = mr.Finish(lo.Map(posts.Posts, func(item *content.Post, i int) func() error {
			return func() error {
				resp.Posts[i] = &core_api.HotPost{
					PostId: item.PostId,
					Title:  item.Title,
				}
				return nil
			}
		})...); err != nil {
			return resp, err
		}
	case core_api.TargetType_CommentContentType:
		comments, err := s.Platform.GetCommentList(ctx, &platform.GetCommentListReq{
			FilterOptions: &platform.CommentFilterOptions{
				OnlyCommentIds: values,
			},
			Pagination: &basic.PaginationOptions{
				Limit: lo.ToPtr(int64(1000)),
			}})
		if err != nil {
			return resp, err
		}
		resp.Comments = make([]*core_api.HotComment, len(comments.Comments))
		if err = mr.Finish(lo.Map(comments.Comments, func(item *platform.Comment, i int) func() error {
			return func() error {
				resp.Comments[i] = &core_api.HotComment{
					CommentId:  item.CommentId,
					Content:    item.Content,
					UserId:     item.UserId,
					ItemId:     item.SubjectId,
					CreateTime: item.CreateTime,
				}
				user, err := s.CloudMindContent.GetUser(ctx, &content.GetUserReq{
					UserId: item.UserId,
				})
				if err != nil {
					return err
				}
				resp.Comments[i].UserName = user.Name
				return nil
			}
		})...); err != nil {
			return resp, err
		}
	}
	return resp, nil
}
