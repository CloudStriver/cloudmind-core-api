package service

import (
	"context"
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/kq"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_content"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/platform_relation"
	"github.com/CloudStriver/cloudmind-mq/app/util/message"
	"github.com/CloudStriver/go-pkg/utils/pconvertor"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/platform/relation"
	"github.com/bytedance/sonic"
	"github.com/google/wire"
)

type IRelationDomainService interface {
	CreateRelation(ctx context.Context, r *core_api.Relation) (err error)
}
type RelationDomainService struct {
	Config               *config.Config
	PlatFormRelation     platform_relation.IPlatFormRelation
	CloudMindContent     cloudmind_content.ICloudMindContent
	CreateNotificationKq *kq.CreateNotificationsKq
	CreateFeedBackKq     *kq.CreateFeedBackKq
}

var RelationDomainServiceSet = wire.NewSet(
	wire.Struct(new(RelationDomainService), "*"),
	wire.Bind(new(IRelationDomainService), new(*RelationDomainService)),
)

func (s *RelationDomainService) CreateRelation(ctx context.Context, r *core_api.Relation) (err error) {
	ok, err := s.PlatFormRelation.CreateRelation(ctx, &relation.CreateRelationReq{
		FromType:     int64(r.FromType),
		FromId:       r.FromId,
		ToType:       int64(r.ToType),
		ToId:         r.ToId,
		RelationType: int64(r.RelationType),
	})
	if err != nil {
		return err
	}

	if !ok.Ok {
		return nil
	}

	act := r.RelationType
	if r.RelationType == core_api.RelationType_HateRelationType {
		act = core_api.RelationType(content.Action_HateType)
	}

	userId := ""
	toName := ""
	var reqs *content.IncrHotValueReq
	switch r.ToType {
	case core_api.TargetType_UserType:
		reqs = &content.IncrHotValueReq{
			Action:     content.Action(act),
			HotId:      r.ToId,
			TargetType: content.TargetType_UserType,
		}
	case core_api.TargetType_FileType:
		reqs = &content.IncrHotValueReq{
			Action:     content.Action(act),
			HotId:      r.ToId,
			TargetType: content.TargetType_FileType,
		}
	case core_api.TargetType_PostType:
		reqs = &content.IncrHotValueReq{
			Action:     content.Action(act),
			HotId:      r.ToId,
			TargetType: content.TargetType_PostType,
		}
	}
	if _, err = s.CloudMindContent.IncrHotValue(ctx, reqs); err != nil {
		return err
	}

	if r.ToId == r.FromId {
		return nil
	}

	switch r.ToType {
	case core_api.TargetType_UserType:
		userId = r.ToId
	case core_api.TargetType_FileType:
		getFileResp, err := s.CloudMindContent.GetFile(ctx, &content.GetFileReq{
			FileId: r.ToId,
		})
		if err != nil {
			return err
		}

		toName = getFileResp.File.Name
		userId = getFileResp.File.UserId

	case core_api.TargetType_PostType:
		getPostResp, err := s.CloudMindContent.GetPost(ctx, &content.GetPostReq{
			PostId: r.ToId,
		})
		if err != nil {
			return err
		}
		toName = getPostResp.Title
		userId = getPostResp.UserId
	}

	userinfo, err := s.CloudMindContent.GetUser(ctx, &content.GetUserReq{
		UserId: r.FromId,
	})
	if err != nil {
		return err
	}

	// 创建通知
	msg, _ := sonic.Marshal(Msg{
		FromName: userinfo.Name,
		ToName:   toName,
	})
	data, _ := sonic.Marshal(&message.CreateNotificationMessage{
		TargetUserId:    userId,
		SourceUserId:    r.FromId,
		SourceContentId: r.ToId,
		TargetType:      int64(r.ToType),
		Type:            int64(r.RelationType),
		Text:            pconvertor.Bytes2String(msg),
	})
	if err = s.CreateNotificationKq.Push(pconvertor.Bytes2String(data)); err != nil {
		return err
	}

	data, _ = sonic.Marshal(&message.CreateFeedBackMessage{
		FeedbackType: core_api.RelationType_name[int32(r.RelationType)],
		UserId:       r.FromId,
		ItemId:       r.ToId,
	})
	if err = s.CreateFeedBackKq.Push(pconvertor.Bytes2String(data)); err != nil {
		return err
	}

	return nil
}

type Msg struct {
	FromName string
	ToName   string
}