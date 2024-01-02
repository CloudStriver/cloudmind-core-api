package service

import (
	"context"
	"github.com/CloudStriver/cloudmind-core-api/biz/adaptor"
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/consts"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_content"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_sts"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/util/oauth"
	"github.com/CloudStriver/go-pkg/utils/util/log"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/sts"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/wire"
	"github.com/samber/lo"
	"strconv"
	"time"
)

type IAuthService interface {
	Register(ctx context.Context, req *core_api.RegisterReq) (resp *core_api.RegisterResp, err error)
	RefreshToken(ctx context.Context, req *core_api.RefreshTokenReq) (resp *core_api.RefreshTokenResp, err error)
	SendEmail(ctx context.Context, req *core_api.SendEmailReq) (resp *core_api.SendEmailResp, err error)
	SetPasswordByEmail(ctx context.Context, req *core_api.SetPasswordByEmailReq) (resp *core_api.SetPasswordByEmailResp, err error)
	SetPasswordByPassword(ctx context.Context, req *core_api.SetPasswordByPasswordReq) (resp *core_api.SetPasswordByPasswordResp, err error)
	GetCaptcha(ctx context.Context, req *core_api.GetCaptchaReq) (resp *core_api.GetCaptchaResp, err error)
	EmailLogin(ctx context.Context, req *core_api.EmailLoginReq) (resp *core_api.EmailLoginResp, err error)
	GithubLogin(ctx context.Context, req *core_api.GithubLoginReq) (resp *core_api.GithubLoginResp, err error)
	GiteeLogin(ctx context.Context, req *core_api.GiteeLoginReq) (resp *core_api.GiteeLoginResp, err error)
}

var AuthServiceSet = wire.NewSet(
	wire.Struct(new(AuthService), "*"),
	wire.Bind(new(IAuthService), new(*AuthService)),
)

type AuthService struct {
	Config           *config.Config
	CloudMindContent cloudmind_content.ICloudMindContent
	CloudMindSts     cloudmind_sts.ICloudMindSts
}

func (s *AuthService) GiteeLogin(ctx context.Context, req *core_api.GiteeLoginReq) (resp *core_api.GiteeLoginResp, err error) {
	resp = new(core_api.GiteeLoginResp)
	if resp.ShortToken, resp.LongToken, resp.UserId, err = s.ThirdLogin(ctx, req.Code, sts.AuthType_gitee); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *AuthService) EmailLogin(ctx context.Context, req *core_api.EmailLoginReq) (resp *core_api.EmailLoginResp, err error) {
	resp = new(core_api.EmailLoginResp)
	loginResp, err := s.CloudMindSts.Login(ctx, &sts.LoginReq{
		Auth:     &sts.AuthInfo{AuthType: sts.AuthType_email, AppId: req.Email},
		Password: req.Password,
		Captcha: &sts.Captcha{
			Point: &sts.Point{X: req.Point.X, Y: req.Point.Y},
			Key:   req.Key,
		},
	})
	if err != nil {
		return resp, err
	}

	resp.ShortToken, resp.LongToken, err = generateShortLongToken(ctx, s.Config.Auth.SecretKey, loginResp.UserId, s.Config.Auth.AccessExpire)
	if err != nil {
		return resp, consts.ErrAuthentication
	}
	resp.UserId = loginResp.UserId
	return resp, nil
}

func (s *AuthService) GithubLogin(ctx context.Context, req *core_api.GithubLoginReq) (resp *core_api.GithubLoginResp, err error) {
	resp = new(core_api.GithubLoginResp)
	if resp.ShortToken, resp.LongToken, resp.UserId, err = s.ThirdLogin(ctx, req.Code, sts.AuthType_github); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *AuthService) ThirdLogin(ctx context.Context, code string, authType sts.AuthType) (shortToken string, longToken string, userId string, err error) {
	// 第三方登录
	data, err := oauth.OauthLogin(s.Config.GiteeConf, authType, code)
	if err != nil {
		return "", "", "", consts.ErrThirdLogin
	}

	// 登录到系统
	var loginResp *sts.LoginResp
	if loginResp, err = s.CloudMindSts.Login(ctx, &sts.LoginReq{
		Auth: &sts.AuthInfo{
			AuthType: sts.AuthType_github,
			AppId:    strconv.FormatInt(data.Id, 10),
		},
	}); err != nil {
		return "", "", "", err

	}
	if loginResp.UserId == "" {
		// 第一次登录
		createAuthResp, err := s.CloudMindSts.CreateAuth(ctx, &sts.CreateAuthReq{
			AuthInfo: &sts.AuthInfo{
				AuthType: sts.AuthType_github,
				AppId:    strconv.FormatInt(data.Id, 10),
			},
			UserInfo: &sts.UserInfo{
				Role: sts.Role_user,
			},
		})
		if err != nil {
			return "", "", "", err
		}
		if _, err = s.CloudMindContent.CreateUser(ctx, &content.CreateUserReq{
			UserInfo: &content.UserInfo{
				UserId: createAuthResp.UserId,
				Name:   data.Name,
				Sex:    content.Sex_unKnownSex,
			},
		}); err != nil {
			return "", "", "", err
		}

		userId = createAuthResp.UserId
	} else {
		userId = loginResp.UserId
	}
	shortToken, longToken, err = generateShortLongToken(ctx, s.Config.Auth.SecretKey, userId, s.Config.Auth.AccessExpire)
	if err != nil {
		log.CtxError(ctx, "生成长短token异常[%v]\n", err)
		return "", "", "", consts.ErrAuthentication
	}
	return shortToken, longToken, userId, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, req *core_api.RefreshTokenReq) (resp *core_api.RefreshTokenResp, err error) {
	resp = new(core_api.RefreshTokenResp)
	claims := make(jwt.MapClaims)
	token, err := jwt.ParseWithClaims(req.LongToken, claims, func(_ *jwt.Token) (interface{}, error) {
		return jwt.ParseECPublicKeyFromPEM([]byte(s.Config.Auth.PublicKey))
	})
	if err != nil {
		log.CtxError(ctx, "解析token异常[%v]\n", err)
		return resp, consts.ErrParseToken
	}
	if !token.Valid {
		return resp, consts.ErrNotAuthentication
	}
	userId, ok := claims["userId"].(string)
	if !ok {
		return resp, consts.ErrNotAuthentication
	}
	if claims["expireTime"].(float64) <= float64(s.Config.Auth.AccessExpire) {
		return resp, consts.ErrNotLongToken
	}

	resp.ShortToken, resp.LongToken, err = generateShortLongToken(ctx, s.Config.Auth.SecretKey, userId, s.Config.Auth.AccessExpire)
	if err != nil {
		return resp, consts.ErrAuthentication
	}
	return resp, nil
}

func (s *AuthService) Register(ctx context.Context, req *core_api.RegisterReq) (resp *core_api.RegisterResp, err error) {
	resp = new(core_api.RegisterResp)
	createAuthResp, err := s.CloudMindSts.CreateAuth(ctx, &sts.CreateAuthReq{
		AuthInfo: &sts.AuthInfo{
			AuthType: sts.AuthType_email,
			AppId:    req.Email,
		},
		UserInfo: &sts.UserInfo{
			Role:     sts.Role_user,
			Password: lo.ToPtr(req.Password),
		},
		Code: req.Code,
	})

	if err != nil {
		return resp, err
	}

	if _, err = s.CloudMindContent.CreateUser(ctx, &content.CreateUserReq{
		UserInfo: &content.UserInfo{
			UserId: createAuthResp.UserId,
			Name:   req.Name,
			Sex:    content.Sex(req.Sex),
		},
	}); err != nil {
		return resp, err
	}

	resp.ShortToken, resp.LongToken, err = generateShortLongToken(ctx, s.Config.Auth.SecretKey, createAuthResp.UserId, s.Config.Auth.AccessExpire)
	if err != nil {
		return resp, consts.ErrAuthentication
	}
	resp.UserId = createAuthResp.UserId
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

func (s *AuthService) SendEmail(ctx context.Context, req *core_api.SendEmailReq) (resp *core_api.SendEmailResp, err error) {
	resp = new(core_api.SendEmailResp)
	if _, err = s.CloudMindSts.SendEmail(ctx, &sts.SendEmailReq{
		Email:   req.Email,
		Subject: req.Subject,
	}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *AuthService) SetPasswordByEmail(ctx context.Context, req *core_api.SetPasswordByEmailReq) (resp *core_api.SetPasswordByEmailResp, err error) {
	resp = new(core_api.SetPasswordByEmailResp)
	if _, err = s.CloudMindSts.SetPassword(ctx, &sts.SetPasswordReq{
		Key: &sts.SetPasswordReq_EmailOptions{
			EmailOptions: &sts.EmailOptions{
				Email: req.Email,
				Code:  req.Code,
			},
		},
		Password: req.Password,
	}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *AuthService) SetPasswordByPassword(ctx context.Context, req *core_api.SetPasswordByPasswordReq) (resp *core_api.SetPasswordByPasswordResp, err error) {
	resp = new(core_api.SetPasswordByPasswordResp)
	userData := adaptor.ExtractUserMeta(ctx)
	if userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	if _, err = s.CloudMindSts.SetPassword(ctx, &sts.SetPasswordReq{
		Key: &sts.SetPasswordReq_UserIdOptions{
			UserIdOptions: &sts.UserIdOptions{
				UserId:   userData.GetUserId(),
				Password: req.OldPassword,
			},
		},
		Password: req.Password,
	}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *AuthService) GetCaptcha(ctx context.Context, req *core_api.GetCaptchaReq) (resp *core_api.GetCaptchaResp, err error) {
	resp = new(core_api.GetCaptchaResp)
	createCaptchaResp, err := s.CloudMindSts.CreateCaptcha(ctx, &sts.CreateCaptchaReq{})
	if err != nil {
		return resp, err
	}

	resp.Key = createCaptchaResp.Key
	resp.JigsawImageBase64 = createCaptchaResp.JigsawImageBase64
	resp.OriginalImageBase64 = createCaptchaResp.OriginalImageBase64
	return resp, nil
}
