package service

import (
	"context"
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/consts"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_content"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_sts"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_system"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_trade"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/platform_comment"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/platform/comment"
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
	Config           *config.Config
	CloudMindContent cloudmind_content.ICloudMindContent
	CloudMindSts     cloudmind_sts.ICloudMindSts
	CloudMindTrade   cloudmind_trade.ICloudMindTrade
	CloudMindSystem  cloudmind_system.ICloudMindSystem
	PlatFormComment  platform_comment.IPlatFormComment
	Redis            *redis.Redis
}

func (s *HotRankService) GetHotRanks(ctx context.Context, req *core_api.GetHotRanksReq) (resp *core_api.GetHotRanksResp, err error) {
	resp = new(core_api.GetHotRanksResp)
	key := ""
	switch req.TargetType {
	case core_api.TargetType_UserType:
		key = consts.UserRankKey
	case core_api.TargetType_FileType:
		key = consts.FileRankKey
	case core_api.TargetType_PostType:
		key = consts.PostRankKey
	}

	values, err := s.Redis.ZrangeCtx(ctx, key, 0, -1)
	if err != nil {
		return resp, err
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
					UserId: item.UserId,
					Name:   item.Name,
					Url:    item.Url,
				}
				if err = mr.Finish(func() error {
					hotValue, _ := s.CloudMindContent.GetHotValue(ctx, &content.GetHotValueReq{
						HotId: item.UserId,
					})

					resp.Users[i].HotValue = hotValue.HotValue
					return nil
				}, func() error {
					tags, _ := s.PlatFormComment.GetLabelsInBatch(ctx, &comment.GetLabelsInBatchReq{
						LabelIds: item.Labels,
					})
					resp.Users[i].Tags = tags.Labels
					return nil
				}); err != nil {
					return err
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
		if err = mr.Finish(lo.Map(files.Files, func(item *content.FileInfo, i int) func() error {
			return func() error {
				resp.Files[i] = &core_api.HotFile{
					FileId: item.FileId,
					Name:   item.Name,
					Type:   item.Type,
				}
				hotValue, _ := s.CloudMindContent.GetHotValue(ctx, &content.GetHotValueReq{
					HotId: item.UserId,
				})
				resp.Files[i].HotValue = hotValue.HotValue
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
					PostId:   item.PostId,
					Title:    item.Title,
					Url:      item.Url,
					UserName: item.UserId,
				}
				if err = mr.Finish(func() error {
					hotValue, err := s.CloudMindContent.GetHotValue(ctx, &content.GetHotValueReq{
						HotId: item.UserId,
					})
					if err == nil {
						resp.Posts[i].HotValue = hotValue.HotValue
					}
					return nil
				}, func() error {
					user, err := s.CloudMindContent.GetUser(ctx, &content.GetUserReq{
						UserId: item.UserId,
					})
					if err == nil {
						resp.Posts[i].UserName = user.Name
					}
					return nil
				}); err != nil {
					return err
				}
				return nil
			}
		})...); err != nil {
			return resp, err
		}
	}
	return resp, nil
}
