package service

import (
	"context"
	"github.com/CloudStriver/cloudmind-core-api/biz/adaptor"
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/cloudmind-core-api/biz/domain/service"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/consts"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/kq"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_content"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_sts"
	platformservice "github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/platform"
	"github.com/CloudStriver/cloudmind-mq/app/util/message"
	"github.com/CloudStriver/go-pkg/utils/pconvertor"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/basic"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/sts"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/platform"
	"github.com/bytedance/sonic"
	"github.com/google/wire"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/mr"
)

type IPostService interface {
	CreatePost(ctx context.Context, req *core_api.CreatePostReq) (resp *core_api.CreatePostResp, err error)
	UpdatePost(ctx context.Context, req *core_api.UpdatePostReq) (resp *core_api.UpdatePostResp, err error)
	DeletePost(ctx context.Context, req *core_api.DeletePostReq) (resp *core_api.DeletePostResp, err error)
	GetPosts(ctx context.Context, c *core_api.GetPostsReq) (*core_api.GetPostsResp, error)
	GetPost(ctx context.Context, c *core_api.GetPostReq) (*core_api.GetPostResp, error)
}

var PostServiceSet = wire.NewSet(
	wire.Struct(new(PostService), "*"),
	wire.Bind(new(IPostService), new(*PostService)),
)

type PostService struct {
	Config                *config.Config
	CloudMindContent      cloudmind_content.ICloudMindContent
	PostDomainService     service.IPostDomainService
	Platform              platformservice.IPlatForm
	CloudMindSts          cloudmind_sts.ICloudMindSts
	RelationDomainService service.IRelationDomainService
	UserDomainService     service.IUserDomainService
	CreateItemKq          *kq.CreateItemKq
	UpdateItemKq          *kq.UpdateItemKq
	DeleteItemKq          *kq.DeleteItemKq
}

func (s *PostService) FiltetContet(ctx context.Context, IsSure bool, contents []*string) ([]*core_api.Keywords, error) {
	cts := lo.Map[*string, string](contents, func(item *string, index int) string {
		return *item
	})
	if IsSure {
		replaceContentResp, err := s.CloudMindSts.ReplaceContent(ctx, &sts.ReplaceContentReq{
			Contents: cts,
		})
		if err != nil {
			return nil, err
		}
		for i, c := range replaceContentResp.Content {
			*contents[i] = c
		}
		return nil, nil
	} else {
		// 内容检测
		findAllContentResp, err := s.CloudMindSts.FindAllContent(ctx, &sts.FindAllContentReq{
			Contents: cts,
		})
		if err != nil {
			return nil, err
		}
		keywords := make([]*core_api.Keywords, 0, len(findAllContentResp.Keywords))
		for _, keyword := range findAllContentResp.Keywords {
			if len(keyword.Keywords) != 0 {
				keywords = append(keywords, &core_api.Keywords{
					Keywords: keyword.Keywords,
				})
			}
		}
		if len(keywords) != 0 {
			return keywords, nil
		}
		return nil, nil
	}
}

func (s *PostService) CreatePost(ctx context.Context, req *core_api.CreatePostReq) (resp *core_api.CreatePostResp, err error) {
	resp = new(core_api.CreatePostResp)
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	if req.Status == int64(core_api.PostStatus_PublicPostStatus) {
		resp.Keywords, err = s.FiltetContet(ctx, req.IsSure, []*string{&req.Title, &req.Text})
		if err != nil {
			return resp, err
		}
		if resp.Keywords != nil {
			return resp, nil
		}
	}

	tags := lo.Map[*core_api.Tag, *content.Tag](req.Tags, func(item *core_api.Tag, index int) *content.Tag {
		return &content.Tag{
			TagId:  item.TagId,
			ZoneId: item.ZoneId,
		}
	})
	tagIds := lo.Map[*core_api.Tag, string](req.Tags, func(item *core_api.Tag, index int) string {
		return item.TagId
	})

	createPostResp, err := s.CloudMindContent.CreatePost(ctx, &content.CreatePostReq{
		UserId: userData.UserId,
		Title:  req.Title,
		Text:   req.Text,
		Tags:   tags,
		Status: req.Status,
		Url:    req.Url,
	})
	if err != nil {
		return resp, err
	}
	resp.PostId = createPostResp.PostId

	if err = mr.Finish(func() error {
		if _, err1 := s.Platform.CreateCommentSubject(ctx, &platform.CreateCommentSubjectReq{
			Subject: &platform.Subject{
				Id:        createPostResp.PostId,
				UserId:    userData.UserId,
				RootCount: lo.ToPtr(int64(0)),
				AllCount:  lo.ToPtr(int64(0)),
				State:     int64(platform.State_Normal),
				Attrs:     int64(platform.Attrs_None),
			},
		}); err1 != nil {
			return err1
		}
		return nil
	}, func() error {
		if _, err2 := s.CloudMindContent.CreateHot(ctx, &content.CreateHotReq{
			HotId: createPostResp.PostId,
		}); err2 != nil {
			return err2
		}
		return nil
	}, func() error {
		data, _ := sonic.Marshal(&message.CreateItemMessage{
			ItemId:   createPostResp.PostId,
			IsHidden: req.Status == int64(core_api.PostStatus_DraftPostStatus),
			Labels:   tagIds,
			Category: core_api.Category_name[int32(core_api.Category_PostCategory)],
		})
		if err3 := s.CreateItemKq.Push(pconvertor.Bytes2String(data)); err != nil {
			return err3
		}
		return nil
	}, func() error {
		if req.Status == int64(core_api.PostStatus_PublicPostStatus) {
			if err4 := s.RelationDomainService.CreateRelation(ctx, &core_api.Relation{
				FromType:     core_api.TargetType_UserType,
				FromId:       userData.UserId,
				ToType:       core_api.TargetType_PostType,
				ToId:         createPostResp.PostId,
				RelationType: core_api.RelationType_PublishRelationType,
			}); err4 != nil {
				return err4
			}
		}
		return nil
	}); err != nil {
		return resp, err
	}

	return resp, nil
}

func (s *PostService) UpdatePost(ctx context.Context, req *core_api.UpdatePostReq) (resp *core_api.UpdatePostResp, err error) {
	resp = new(core_api.UpdatePostResp)
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}

	post, err := s.CheckIsMyPost(ctx, req.PostId, userData.UserId)
	if err != nil {
		return resp, err
	}
	if post.Status == int64(core_api.PostStatus_PublicPostStatus) || req.Status == int64(core_api.PostStatus_PublicPostStatus) {
		if req.Title == "" {
			req.Title = post.Title
		}
		if req.Text == "" {
			req.Text = post.Text
		}
		resp.Keywords, err = s.FiltetContet(ctx, req.IsSure, []*string{&req.Title, &req.Text})
		if err != nil {
			return resp, err
		}
		if resp.Keywords != nil {
			return resp, nil
		}
	}

	tags := lo.Map[*core_api.Tag, *content.Tag](req.Tags, func(item *core_api.Tag, index int) *content.Tag {
		return &content.Tag{
			TagId:  item.TagId,
			ZoneId: item.ZoneId,
		}
	})
	tagIds := lo.Map[*core_api.Tag, string](req.Tags, func(item *core_api.Tag, index int) string {
		return item.TagId
	})

	if _, err = s.CloudMindContent.UpdatePost(ctx, &content.UpdatePostReq{
		PostId: req.PostId,
		Title:  req.Title,
		Text:   req.Text,
		Tags:   tags,
		Status: req.Status,
		Url:    req.Url,
	}); err != nil {
		return resp, err
	}

	if req.Status != 0 {
		if err = mr.Finish(func() error {
			data, _ := sonic.Marshal(&message.UpdateItemMessage{
				ItemId:   req.PostId,
				IsHidden: lo.ToPtr(req.Status != int64(core_api.PostStatus_PublicPostStatus)),
				Labels:   tagIds,
			})
			if err1 := s.UpdateItemKq.Push(pconvertor.Bytes2String(data)); err != nil {
				return err1
			}
			return nil
		}, func() error {
			if err2 := s.RelationDomainService.CreateRelation(ctx, &core_api.Relation{
				FromType:     core_api.TargetType_UserType,
				FromId:       userData.UserId,
				ToType:       core_api.TargetType_PostType,
				ToId:         req.PostId,
				RelationType: core_api.RelationType_PublishRelationType,
			}); err2 != nil {
				return err2
			}
			return nil
		}); err != nil {
			return resp, err
		}
	}
	return resp, nil
}

func (s *PostService) DeletePost(ctx context.Context, req *core_api.DeletePostReq) (resp *core_api.DeletePostResp, err error) {
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	// 只能删除自己的帖子
	if err = s.CheckIsMyPosts(ctx, req.PostIds, userData.UserId); err != nil {
		return resp, err
	}

	if _, err = s.CloudMindContent.DeletePost(ctx, &content.DeletePostReq{
		PostId: req.PostIds,
	}); err != nil {
		return resp, err
	}
	if err = mr.Finish(lo.Map(req.PostIds, func(item string, index int) func() (err error) {
		return func() (err error) {
			if err = mr.Finish(func() error {
				data, _ := sonic.Marshal(&message.DeleteItemMessage{
					ItemId: item,
				})
				if err1 := s.DeleteItemKq.Push(pconvertor.Bytes2String(data)); err != nil {
					return err1
				}
				return nil
			}, func() error {
				if _, err2 := s.Platform.DeleteRelation(ctx, &platform.DeleteRelationReq{
					FromType:     int64(core_api.TargetType_UserType),
					FromId:       userData.UserId,
					ToType:       int64(core_api.TargetType_PostType),
					ToId:         item,
					RelationType: int64(core_api.RelationType_PublishRelationType),
				}); err2 != nil {
					return err2
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
	return resp, nil
}

func (s *PostService) GetPost(ctx context.Context, req *core_api.GetPostReq) (resp *core_api.GetPostResp, err error) {
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil {
		return resp, consts.ErrNotAuthentication

	}
	var res *content.GetPostResp
	if res, err = s.CloudMindContent.GetPost(ctx, &content.GetPostReq{
		PostId: req.PostId,
	}); err != nil {
		return resp, err
	}

	// 如果该帖子非公开，并且不是他的，那么他没有权限查看
	if res.Status != int64(core_api.PostStatus_PublicPostStatus) && res.UserId != userData.GetUserId() {
		return resp, consts.ErrForbidden
	}

	tags := lo.Map[*content.Tag, *core_api.TagInfo](res.Tags, func(item *content.Tag, index int) *core_api.TagInfo {
		return &core_api.TagInfo{
			TagId:  item.TagId,
			ZoneId: item.ZoneId,
		}
	})
	tagIds := lo.Map[*content.Tag, string](res.Tags, func(item *content.Tag, index int) string {
		return item.TagId
	})

	resp = &core_api.GetPostResp{
		Title:      res.Title,
		Text:       res.Text,
		Url:        res.Url,
		Tags:       tags,
		CreateTime: res.CreateTime,
		UpdateTime: res.UpdateTime,
		Author:     &core_api.User{},
	}

	if err = mr.Finish(func() error {
		s.PostDomainService.LoadAuthor(ctx, resp.Author, res.UserId) // 作者
		return nil
	}, func() error {
		s.PostDomainService.LoadLikeCount(ctx, &resp.LikeCount, req.PostId) // 点赞量
		return nil
	}, func() error {
		s.PostDomainService.LoadViewCount(ctx, &resp.ViewCount, req.PostId) // 浏览量
		return nil
	}, func() error {
		s.PostDomainService.LoadCollectCount(ctx, &resp.CollectCount, req.PostId) // 收藏量
		return nil
	}, func() error {
		s.PostDomainService.LoadLiked(ctx, &resp.Liked, userData.GetUserId(), req.PostId) // 是否点赞
		return nil
	}, func() error {
		s.PostDomainService.LoadCollected(ctx, &resp.Collected, userData.GetUserId(), req.PostId) // 是否收藏
		return nil
	}, func() error {
		if userData.GetUserId() != "" {
			_, _ = s.Platform.CreateRelation(ctx, &platform.CreateRelationReq{
				FromType:     int64(core_api.TargetType_UserType),
				FromId:       userData.UserId,
				ToType:       int64(core_api.TargetType_PostType),
				ToId:         req.PostId,
				RelationType: int64(core_api.RelationType_ViewRelationType),
			})
		}
		return nil
	}, func() error {
		s.PostDomainService.LoadLabels(ctx, tagIds)
		for i := range resp.Tags {
			resp.Tags[i].Value = tagIds[i]
		}
		return nil
	}, func() error {
		s.UserDomainService.LoadFollowed(ctx, &resp.Author.Followed, userData.GetUserId(), res.UserId)
		return nil
	}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *PostService) GetPosts(ctx context.Context, req *core_api.GetPostsReq) (resp *core_api.GetPostsResp, err error) {
	resp = new(core_api.GetPostsResp)
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil {
		return resp, consts.ErrNotAuthentication
	}
	var (
		getPostsResp  *content.GetPostsResp
		searchOptions *content.SearchOptions
	)

	if req.AllFieldsKey != nil {
		searchOptions = &content.SearchOptions{
			Type: &content.SearchOptions_AllFieldsKey{
				AllFieldsKey: *req.AllFieldsKey,
			},
		}
	}
	if req.Text != nil || req.Title != nil || req.Tag != nil {
		searchOptions = &content.SearchOptions{
			Type: &content.SearchOptions_MultiFieldsKey{
				MultiFieldsKey: &content.SearchField{
					Tag:   req.Tag,
					Text:  req.Text,
					Title: req.Title,
				},
			},
		}
	}

	filter := &content.PostFilterOptions{
		OnlyUserId: req.OnlyUserId,
		OnlyTag:    req.OnlyTag,
		OnlyZoneId: req.OnlyZoneId,
	}

	// 查看所有人的，或者查看的不是自己的
	if req.OnlyUserId == nil || req.GetOnlyUserId() != userData.GetUserId() {
		filter.OnlyStatus = lo.ToPtr(int64(core_api.PostStatus_PublicPostStatus))
	}

	// 查看的自己的
	if req.GetOnlyUserId() != "" && req.GetOnlyUserId() == userData.GetUserId() {
		filter.OnlyStatus = req.OnlyStatus
	}

	if getPostsResp, err = s.CloudMindContent.GetPosts(ctx, &content.GetPostsReq{
		SearchOptions:     searchOptions,
		PostFilterOptions: filter,
		PaginationOptions: &basic.PaginationOptions{
			Limit:     req.Limit,
			LastToken: req.LastToken,
			Backward:  req.Backward,
			Offset:    req.Offset,
		},
	}); err != nil {
		return resp, err
	}

	resp.Posts = make([]*core_api.Post, len(getPostsResp.Posts))
	if err = mr.Finish(lo.Map(getPostsResp.Posts, func(item *content.Post, i int) func() error {
		return func() error {
			tags := lo.Map[*content.Tag, *core_api.TagInfo](item.Tags, func(item *content.Tag, index int) *core_api.TagInfo {
				return &core_api.TagInfo{
					TagId:  item.TagId,
					ZoneId: item.ZoneId,
				}
			})
			tagsId := lo.Map[*content.Tag, string](item.Tags, func(item *content.Tag, index int) string {
				return item.TagId
			})
			resp.Posts[i] = &core_api.Post{
				PostId: item.PostId,
				Title:  item.Title,
				Text:   item.Text,
				Url:    item.Url,
				Tags:   tags,
			}
			author := &core_api.User{}
			if err = mr.Finish(func() error {
				s.PostDomainService.LoadLikeCount(ctx, &resp.Posts[i].LikeCount, item.PostId) // 点赞量
				return nil
			}, func() error {
				// 加载评论量
				return nil
			}, func() error {
				if userData.GetUserId() != "" {
					s.PostDomainService.LoadLiked(ctx, &resp.Posts[i].Liked, userData.GetUserId(), item.PostId)
				}
				return nil
			}, func() error {
				s.PostDomainService.LoadAuthor(ctx, author, item.UserId)
				resp.Posts[i].UserName = author.Name
				return nil
			}, func() error {
				s.PostDomainService.LoadLabels(ctx, tagsId)
				for i := range tags {
					tags[i].Value = tagsId[i]
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
	resp.Total = getPostsResp.Total
	resp.Token = getPostsResp.Token
	return resp, nil
}

func (s *PostService) CheckIsMyPost(ctx context.Context, postId, userId string) (*content.GetPostResp, error) {
	getPostResp, err := s.CloudMindContent.GetPost(ctx, &content.GetPostReq{
		PostId: postId,
	})
	if err != nil {
		return nil, err
	}
	if getPostResp.UserId != userId {
		return nil, consts.ErrForbidden
	}
	return getPostResp, nil
}

func (s *PostService) CheckIsMyPosts(ctx context.Context, postIds []string, userId string) (err error) {
	post, err := s.CloudMindContent.GetPostsByPostIds(ctx, &content.GetPostsByPostIdsReq{
		PostIds: postIds,
	})
	if err != nil {
		return err
	}
	for i := range post.Posts {
		if post.Posts[i].UserId != userId {
			return consts.ErrForbidden
		}
	}
	return nil
}
