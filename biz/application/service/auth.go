package service

import (
	"context"
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/consts"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_sts"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_user"
	"github.com/CloudStriver/go-pkg/utils/util/log"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/sts"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/user"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/wire"
	"google.golang.org/grpc/status"
	"time"
)

type IAuthService interface {
	Register(ctx context.Context, req *core_api.RegisterReq) (resp *core_api.RegisterResp, err error)
	Login(ctx context.Context, req *core_api.LoginReq) (resp *core_api.LoginResp, err error)
	RefreshToken(ctx context.Context, req *core_api.RefreshTokenReq) (resp *core_api.RefreshTokenResp, err error)
}

type AuthService struct {
	Config        *config.Config
	CloudMindUser cloudmind_user.ICloudMindUser
	CloudMindSts  cloudmind_sts.ICloudMindSts
}

var AuthServiceSet = wire.NewSet(
	wire.Struct(new(AuthService), "*"),
	wire.Bind(new(IAuthService), new(*AuthService)),
)

func (s *AuthService) Login(ctx context.Context, req *core_api.LoginReq) (resp *core_api.LoginResp, err error) {
	resp = new(core_api.LoginResp)
	loginResp, err := s.CloudMindSts.Login(ctx, &sts.LoginReq{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		log.CtxError(ctx, "调用sts.Login异常[%v]\n", err)
		return resp, err
	}
	if loginResp.Error != "" {
		return resp, status.New(consts.NotContent, loginResp.Error).Err()
	}

	resp.ShortToken, resp.LongToken, err = generateShortLongToken(ctx, s.Config.Auth.SecretKey, loginResp.UserId, s.Config.Auth.AccessExpire)
	if err != nil {
		log.CtxError(ctx, "生成长短token异常[%v]\n", err)
		return resp, err
	}
	resp.UserId = loginResp.UserId
	return resp, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, req *core_api.RefreshTokenReq) (resp *core_api.RefreshTokenResp, err error) {
	resp = new(core_api.RefreshTokenResp)
	claims := make(jwt.MapClaims)
	token, err := jwt.ParseWithClaims(req.LongToken, claims, func(_ *jwt.Token) (interface{}, error) {
		return jwt.ParseECPublicKeyFromPEM([]byte(s.Config.Auth.PublicKey))
	})
	if err != nil {
		log.CtxError(ctx, "解析token异常[%v]\n", err)
		return resp, err
	}
	if !token.Valid {
		return resp, consts.ErrNotAuthentication
	}
	userId, ok := claims["userId"].(string)
	if !ok {
		return resp, consts.ErrNotAuthentication
	}
	if claims["expireTime"].(float64) <= float64(s.Config.Auth.AccessExpire) {
		return resp, status.New(consts.NotContent, "请使用LongToken刷新").Err()
	}

	resp.ShortToken, resp.LongToken, err = generateShortLongToken(ctx, s.Config.Auth.SecretKey, userId, s.Config.Auth.AccessExpire)
	if err != nil {
		log.CtxError(ctx, "生成长短token异常[%v]\n", err)
		return resp, err
	}
	return resp, nil
}

func (s *AuthService) Register(ctx context.Context, req *core_api.RegisterReq) (resp *core_api.RegisterResp, err error) {
	resp = new(core_api.RegisterResp)
	createUserResp, err := s.CloudMindUser.CreateUser(ctx, &user.CreateUserReq{
		Name: req.Name,
		Sex:  user.Sex(req.Sex),
	})
	if err != nil {
		log.CtxError(ctx, "调用user.CreateUser异常[%v]\n", err)
		return resp, err
	}
	if createUserResp.Error != "" {
		return resp, status.New(consts.NotContent, createUserResp.Error).Err()
	}

	sresp, err := s.CloudMindSts.CreateAuth(ctx, &sts.CreateAuthReq{
		Type:     sts.AuthType(req.RegisterType),
		Key:      req.AuthKey,
		UserId:   createUserResp.UserId,
		Role:     sts.Role_User,
		Password: req.Password,
	})
	if err != nil {
		log.CtxError(ctx, "调用sts.CreateAuth异常[%v]\n", err)
		return resp, err
	}
	if sresp.Error != "" {
		// 删除已经创建的user
		if _, err = s.CloudMindUser.DeleteUser(ctx, &user.DeleteUserReq{
			UserId: createUserResp.UserId,
		}); err != nil {
			log.CtxError(ctx, "调用user.DeleteUser异常[%v]\n", err)
			return resp, err
		}
		return resp, status.New(consts.NotContent, sresp.Error).Err()
	}

	resp.ShortToken, resp.LongToken, err = generateShortLongToken(ctx, s.Config.Auth.SecretKey, createUserResp.UserId, s.Config.Auth.AccessExpire)
	if err != nil {
		log.CtxError(ctx, "生成长短token异常[%v]\n", err)
		return resp, err
	}
	resp.UserId = createUserResp.UserId
	return resp, nil
}

func generateShortLongToken(ctx context.Context, secretKey, userId string, accessExpire int64) (shortToken, longToken string, err error) {
	shortToken, _, err = generateJwtToken(userId, secretKey, accessExpire)
	if err != nil {
		log.CtxError(ctx, "生成短token异常[%v]\n", err)
		return "", "", err
	}
	longToken, _, err = generateJwtToken(userId, secretKey, 24*30*accessExpire)
	if err != nil {
		log.CtxError(ctx, "生成长token异常[%v]\n", err)
		return "", "", err
	}
	return shortToken, longToken, nil
}
func generateJwtToken(userId, secret string, expire int64) (string, int64, error) {
	key, err := jwt.ParseECPrivateKeyFromPEM([]byte(secret))
	if err != nil {
		return "", 0, err
	}
	iat := time.Now().Unix()
	exp := iat + expire
	claims := make(jwt.MapClaims)
	claims["exp"] = exp
	claims["iat"] = iat
	claims["expireTime"] = expire
	claims["userId"] = userId
	token := jwt.New(jwt.SigningMethodES256)
	token.Claims = claims
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", 0, err
	}
	return tokenString, exp, nil
}
