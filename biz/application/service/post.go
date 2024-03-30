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
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/platform_comment"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/platform_relation"
	"github.com/CloudStriver/cloudmind-mq/app/util/message"
	"github.com/CloudStriver/go-pkg/utils/pconvertor"
	"github.com/CloudStriver/go-pkg/utils/util/log"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/basic"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/sts"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/platform/comment"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/platform/relation"
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
	Config            *config.Config
	CloudMindContent  cloudmind_content.ICloudMindContent
	PostDomainService service.IPostDomainService
	PlatFormRelation  platform_relation.IPlatFormRelation
	PlatFormComment   platform_comment.IPlatFormComment
	CloudMindSts      cloudmind_sts.ICloudMindSts
	UserDomainService service.IUserDomainService
	CreateItemKq      *kq.CreateItemKq
	UpdateItemKq      *kq.UpdateItemKq
	DeleteItemKq      *kq.DeleteItemKq
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
		for i, content := range replaceContentResp.Content {
			*contents[i] = content
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

	createPostResp, err := s.CloudMindContent.CreatePost(ctx, &content.CreatePostReq{
		UserId: userData.UserId,
		Title:  req.Title,
		Text:   req.Text,
		Tags:   req.Tags,
		Status: req.Status,
		Url:    req.Url,
	})
	if err != nil {
		return resp, err
	}

	if _, err = s.PlatFormComment.CreateCommentSubject(ctx, &comment.CreateCommentSubjectReq{
		Subject: &comment.Subject{
			Id:        createPostResp.PostId,
			UserId:    userData.UserId,
			RootCount: lo.ToPtr(int64(0)),
			AllCount:  lo.ToPtr(int64(0)),
			State:     int64(comment.State_Normal),
			Attrs:     int64(comment.Attrs_None),
		},
	}); err != nil {
		return resp, err
	}

	if _, err = s.CloudMindContent.CreateHot(ctx, &content.CreateHotReq{
		HotId: createPostResp.PostId,
	}); err != nil {
		return resp, err
	}

	data, _ := sonic.Marshal(&message.CreateItemMessage{
		ItemId:   createPostResp.PostId,
		IsHidden: req.Status == int64(core_api.PostStatus_DraftPostStatus),
		Labels:   req.Tags,
		Category: core_api.Category_name[int32(core_api.Category_PostCategory)],
	})
	if err = s.CreateItemKq.Push(pconvertor.Bytes2String(data)); err != nil {
		return resp, err
	}

	resp.PostId = createPostResp.PostId
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
	if _, err = s.CloudMindContent.UpdatePost(ctx, &content.UpdatePostReq{
		PostId: req.PostId,
		Title:  req.Title,
		Text:   req.Text,
		Tags:   req.Tags,
		Status: req.Status,
		Url:    req.Url,
	}); err != nil {
		return resp, err
	}

	if req.Status != 0 || req.Tags != nil {
		var isHidden *bool
		if req.Status != 0 {
			isHidden = lo.ToPtr(req.Status != int64(core_api.PostStatus_PublicPostStatus))
		}

		data, _ := sonic.Marshal(&message.UpdateItemMessage{
			ItemId:   req.PostId,
			IsHidden: isHidden,
			Labels:   req.Tags,
		})
		if err = s.UpdateItemKq.Push(pconvertor.Bytes2String(data)); err != nil {
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

	for i := range req.PostIds {
		go func(i int) {
			data, _ := sonic.Marshal(&message.DeleteItemMessage{
				ItemId: req.PostIds[i],
			})
			if err = s.DeleteItemKq.Push(pconvertor.Bytes2String(data)); err != nil {
				log.CtxError(ctx, "DeleteItemKq.Push", err)
			}
		}(i)
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

	resp = &core_api.GetPostResp{
		Title:      res.Title,
		Text:       res.Text,
		Url:        res.Url,
		Tags:       res.Tags,
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
			_, _ = s.PlatFormRelation.CreateRelation(ctx, &relation.CreateRelationReq{
				FromType:     int64(core_api.TargetType_UserType),
				FromId:       userData.UserId,
				ToType:       int64(core_api.TargetType_PostType),
				ToId:         req.PostId,
				RelationType: int64(core_api.RelationType_ViewType),
			})
		}
		return nil
	}, func() error {
		s.PostDomainService.LoadLabels(ctx, resp.Tags)
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
		OnlyUserId:      req.OnlyUserId,
		OnlyTags:        req.OnlyTags,
		OnlySetRelation: req.OnlySetRelation,
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
			resp.Posts[i] = &core_api.Post{
				PostId: item.PostId,
				Title:  item.Title,
				Text:   item.Text,
				Url:    item.Url,
				Tags:   item.Tags,
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
				s.PostDomainService.LoadLabels(ctx, resp.Posts[i].Tags)
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
