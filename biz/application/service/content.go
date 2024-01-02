package service

import (
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_sts"
	"github.com/google/wire"
)

type IContentService interface {
	//Register(ctx context.Context, req *core_api.RegisterReq) (resp *core_api.RegisterResp, err error)
}

type ContentService struct {
	Config       *config.Config
	CloudMindSts cloudmind_sts.ICloudMindSts
}

var ContentServiceSet = wire.NewSet(
	wire.Struct(new(ContentService), "*"),
	wire.Bind(new(IContentService), new(*ContentService)),
)

//
//func (s *ContentService) Register(ctx context.Context, req *core_api.RegisterReq) (resp *core_api.RegisterResp, err error) {
//	resp = new(core_api.RegisterResp)
//
//	uresp, err := s.CloudMindUser.CreateUser(ctx, &user.CreateUserReq{
//		Name: req.Name,
//		Sex:  user.Sex(req.Sex),
//	})
//	if err != nil {
//		return resp, err
//	}
//	if uresp.Error != "" {
//		return resp, status.New(consts.NotContent, uresp.Error).Err()
//	}
//
//	sresp, err := s.CloudMindSts.CreateAuth(ctx, &sts.CreateAuthReq{
//		Type:     sts.AuthType(req.RegisterType),
//		Key:      req.AuthKey,
//		UserId:   uresp.UserId,
//		Role:     sts.Role_User,
//		Password: req.Password,
//	})
//	if err != nil {
//		return resp, err
//	}
//	if sresp.Error != "" {
//		// 删除已经创建的user
//		if _, err = s.CloudMindUser.DeleteUser(ctx, &user.DeleteUserReq{
//			UserId: uresp.UserId,
//		}); err != nil {
//			return resp, err
//		}
//		return resp, status.New(consts.NotContent, sresp.Error).Err()
//	}
//
//	return resp, nil
//}
