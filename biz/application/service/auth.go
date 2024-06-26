package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/CloudStriver/cloudmind-core-api/biz/adaptor"
	"github.com/CloudStriver/cloudmind-core-api/biz/application/dto/cloudmind/core_api"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/consts"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/kq"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_content"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_sts"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_system"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/rpc/cloudmind_trade"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/utils/oauth"
	"github.com/CloudStriver/cloudmind-mq/app/util/message"
	"github.com/CloudStriver/go-pkg/utils/pconvertor"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/sts"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/system"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/trade"
	"github.com/bytedance/sonic"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/wire"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"net/http"
	"time"
)

type IAuthService interface {
	Register(ctx context.Context, req *core_api.RegisterReq) (resp *core_api.RegisterResp, err error)
	RefreshToken(ctx context.Context, req *core_api.RefreshTokenReq) (resp *core_api.RefreshTokenResp, err error)
	SendEmail(ctx context.Context, req *core_api.SendEmailReq) (resp *core_api.SendEmailResp, err error)
	SetPasswordByEmail(ctx context.Context, req *core_api.SetPasswordByEmailReq) (resp *core_api.SetPasswordByEmailResp, err error)
	SetPasswordByPassword(ctx context.Context, req *core_api.SetPasswordByPasswordReq) (resp *core_api.SetPasswordByPasswordResp, err error)
	EmailLogin(ctx context.Context, req *core_api.EmailLoginReq) (resp *core_api.EmailLoginResp, err error)
	CheckEmail(ctx context.Context, c *core_api.CheckEmailReq) (resp *core_api.CheckEmailResp, err error)
	AskUploadAvatar(ctx context.Context, req *core_api.AskUploadAvatarReq) (resp *core_api.AskUploadAvatarResp, err error)
	WeixinLogin(ctx context.Context, req *core_api.WeixinLoginReq) (resp *core_api.WeixinLoginResp, err error)
	WeixinCallBack(ctx context.Context, req *core_api.WeixinCallBackReq) (resp *core_api.WeixinCallBackResp, err error)
	WeixinIsLogin(ctx context.Context, req *core_api.WeixinIsLoginReq) (resp *core_api.WeixinIsLoginResp, err error)
	QQLogin(ctx context.Context, req *core_api.QQLoginReq) (resp *core_api.QQLoginResp, err error)
}

var AuthServiceSet = wire.NewSet(
	wire.Struct(new(AuthService), "*"),
	wire.Bind(new(IAuthService), new(*AuthService)),
)

type AuthService struct {
	Config           *config.Config
	CloudMindContent cloudmind_content.ICloudMindContent
	CloudMindSts     cloudmind_sts.ICloudMindSts
	CloudMindTrade   cloudmind_trade.ICloudMindTrade
	CloudMindSystem  cloudmind_system.ICloudMindSystem
	CreateItemKq     *kq.CreateItemKq
	Redis            *redis.Redis
}

func (s *AuthService) QQLogin(ctx context.Context, req *core_api.QQLoginReq) (resp *core_api.QQLoginResp, err error) {
	resp = new(core_api.QQLoginResp)
	fmt.Println(req.Code)
	if resp.ShortToken, resp.LongToken, resp.UserId, err = s.ThirdLogin(ctx, req.Code, core_api.LoginType_QQLoginType, "", "", "", consts.SexMan); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *AuthService) WeixinLogin(ctx context.Context, req *core_api.WeixinLoginReq) (resp *core_api.WeixinLoginResp, err error) {
	url := fmt.Sprintf("https://yd.jylt.cc/api/wxLogin/tempUserId?secret=%s", s.Config.WechatConf.Secret)
	var (
		reqs       *http.Request
		res        *http.Response
		httpClient http.Client
		WechatInfo oauth.WechatInfo
	)
	if reqs, err = http.NewRequest(http.MethodGet, url, nil); err != nil {
		return resp, err
	}
	reqs.Header.Set("accept", "application/json")
	if res, err = httpClient.Do(reqs); err != nil {
		return resp, err
	}
	if err = json.NewDecoder(res.Body).Decode(&WechatInfo); err != nil {
		return resp, err
	}
	if WechatInfo.Code != 0 {
		return resp, err
	}

	if err = s.Redis.SetexCtx(ctx, fmt.Sprintf("%s:%s", consts.WechatLoginKey, WechatInfo.Data.TempUserId), "Login", 3000); err != nil {
		return resp, err
	}

	return &core_api.WeixinLoginResp{
		QrUrl:      WechatInfo.Data.QrUrl,
		TempUserId: WechatInfo.Data.TempUserId,
	}, nil
}

func (s *AuthService) WeixinCallBack(ctx context.Context, req *core_api.WeixinCallBackReq) (resp *core_api.WeixinCallBackResp, err error) {
	val, err := s.Redis.GetCtx(ctx, fmt.Sprintf("%s:%s", consts.WechatLoginKey, req.TempUserId))
	if err != nil {
		return resp, err
	}

	if val == "" {
		return resp, consts.ErrThirdLogin
	}

	if req.ScanSuccess {
		if err = s.Redis.SetexCtx(ctx, fmt.Sprintf("%s:%s", consts.WechatLoginKey, req.TempUserId), "ScanSuccess", 3000); err != nil {
			return resp, err
		}
		return resp, nil
	}
	if req.CancelLogin {
		if err = s.Redis.SetexCtx(ctx, fmt.Sprintf("%s:%s", consts.WechatLoginKey, req.TempUserId), "CancelLogin", 3000); err != nil {
			return resp, err
		}
		return resp, nil
	}
	if err = s.Redis.SetexCtx(ctx, fmt.Sprintf("%s:%s", consts.WechatLoginKey, req.TempUserId), "LoginSuccess", 3000); err != nil {
		return resp, err
	}

	var (
		userId string
	)

	if _, _, userId, err = s.ThirdLogin(ctx, "", core_api.LoginType_WeixinLoginType, req.WxMaUserInfo.OpenId, req.WxMaUserInfo.NickName, req.WxMaUserInfo.AvatarUrl, consts.SexMan); err != nil {
		return resp, err
	}

	if err = s.Redis.SetexCtx(ctx, fmt.Sprintf("%s:%s:temp", consts.WechatLoginKey, req.TempUserId), userId, 3000); err != nil {
		return resp, err
	}

	return &core_api.WeixinCallBackResp{
		Code: 0,
		Msg:  "登录成功",
	}, nil
}

func (s *AuthService) WeixinIsLogin(ctx context.Context, req *core_api.WeixinIsLoginReq) (resp *core_api.WeixinIsLoginResp, err error) {
	resp = new(core_api.WeixinIsLoginResp)
	val, err := s.Redis.GetCtx(ctx, fmt.Sprintf("%s:%s", consts.WechatLoginKey, req.TempUserId))
	if err != nil {
		return resp, err
	}

	if val == "LoginSuccess" {
		val, err = s.Redis.GetCtx(ctx, fmt.Sprintf("%s:%s:temp", consts.WechatLoginKey, req.TempUserId))
		if err != nil {
			return resp, err
		}
		resp.ShortToken, resp.LongToken, err = generateShortLongToken(s.Config.Auth.SecretKey, val, s.Config.Auth.ShortTokenExpire, s.Config.Auth.LongTokenExpire)
		if err != nil {
			return resp, err
		}
		resp.UserId = val
		return resp, nil
	} else {
		resp.Status = val
		return resp, nil
	}
}

func (s *AuthService) CheckEmail(ctx context.Context, req *core_api.CheckEmailReq) (resp *core_api.CheckEmailResp, err error) {
	resp = new(core_api.CheckEmailResp)
	checkEmailResp, err := s.CloudMindSts.CheckEmail(ctx, &sts.CheckEmailReq{
		Email: req.Email,
		Code:  req.Code,
	})
	if err != nil {
		return resp, err
	}
	resp.Ok = checkEmailResp.Ok
	return resp, nil
}

func (s *AuthService) EmailLogin(ctx context.Context, req *core_api.EmailLoginReq) (resp *core_api.EmailLoginResp, err error) {
	resp = new(core_api.EmailLoginResp)
	loginResp, err := s.CloudMindSts.Login(ctx, &sts.LoginReq{
		AuthType: int64(core_api.LoginType_EmailLoginType),
		AppId:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return resp, err
	}

	if loginResp.UserId == "" {
		return resp, consts.ErrEmailNotFound
	}

	resp.ShortToken, resp.LongToken, err = generateShortLongToken(s.Config.Auth.SecretKey, loginResp.UserId, s.Config.Auth.ShortTokenExpire, s.Config.Auth.LongTokenExpire)
	if err != nil {
		return resp, consts.ErrAuthentication
	}
	resp.UserId = loginResp.UserId
	return resp, nil
}

func (s *AuthService) ThirdLogin(ctx context.Context, code string, authType core_api.LoginType, appId string, name string, url string, sex int64) (shortToken string, longToken string, userId string, err error) {
	// 第三方登录
	switch authType {
	case core_api.LoginType_QQLoginType:
		conf := s.Config.QQConf
		data, err := oauth.QQLogin(conf, code)
		if err != nil {
			return "", "", "", consts.ErrThirdLogin
		}
		appId = data.Openid
		name = data.Nickname
		url = data.FigureurlQq1
		if data.Gender != "男" {
			sex = consts.SexWoman
		}
	}

	// 登录到系统
	var loginResp *sts.LoginResp
	if loginResp, err = s.CloudMindSts.Login(ctx, &sts.LoginReq{
		AuthType: int64(authType),
		AppId:    appId,
	}); err != nil {
		return "", "", "", err
	}
	if loginResp.UserId == "" {
		// 第一次登录
		createAuthResp, err := s.CloudMindSts.CreateAuth(ctx, &sts.CreateAuthReq{
			AuthType: int64(authType),
			AppId:    appId,
			Role:     int64(core_api.RoleType_UserRoleType),
		})
		if err != nil {
			return "", "", "", err
		}
		if err = s.UserInit(ctx, createAuthResp.UserId, name, sex, url); err != nil {
			return "", "", "", err
		}
		userId = createAuthResp.UserId
	} else {
		userId = loginResp.UserId
	}

	shortToken, longToken, err = generateShortLongToken(s.Config.Auth.SecretKey, userId, s.Config.Auth.ShortTokenExpire, s.Config.Auth.LongTokenExpire)
	if err != nil {
		return "", "", "", consts.ErrAuthentication
	}
	return shortToken, longToken, userId, nil
}

func (s *AuthService) RefreshToken(_ context.Context, req *core_api.RefreshTokenReq) (resp *core_api.RefreshTokenResp, err error) {
	resp = new(core_api.RefreshTokenResp)
	claims := make(jwt.MapClaims)
	token, err := jwt.ParseWithClaims(req.LongToken, claims, func(_ *jwt.Token) (interface{}, error) {
		return jwt.ParseECPublicKeyFromPEM([]byte(s.Config.Auth.PublicKey))
	})
	if err != nil {
		return resp, consts.ErrParseToken
	}
	if !token.Valid {
		return resp, consts.ErrNotAuthentication
	}
	userId, ok := claims["userId"].(string)
	if !ok {
		return resp, consts.ErrNotAuthentication
	}
	if claims["expireTime"].(float64) <= float64(s.Config.Auth.ShortTokenExpire) {
		return resp, consts.ErrNotLongToken
	}

	resp.ShortToken, resp.LongToken, err = generateShortLongToken(s.Config.Auth.SecretKey, userId, s.Config.Auth.ShortTokenExpire, s.Config.Auth.LongTokenExpire)
	if err != nil {
		return resp, consts.ErrAuthentication
	}
	return resp, nil
}

func (s *AuthService) FiltetContet(ctx context.Context, IsSure bool, contents []*string) ([]string, error) {
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
		return findAllContentResp.Keywords, nil
	}
}

func (s *AuthService) Register(ctx context.Context, req *core_api.RegisterReq) (resp *core_api.RegisterResp, err error) {
	resp = new(core_api.RegisterResp)
	value := ""
	if value, err = s.Redis.GetCtx(ctx, fmt.Sprintf("%s:%s", consts.PassCheckEmail, req.Email)); err != nil {
		return resp, err
	}
	if value != "true" {
		return resp, consts.ErrNotEmailCheck
	}

	if _, err = s.Redis.DelCtx(ctx, fmt.Sprintf("%s:%s", consts.PassCheckEmail, req.Email)); err != nil {
		return resp, err
	}

	createAuthResp, err := s.CloudMindSts.CreateAuth(ctx, &sts.CreateAuthReq{
		AuthType: int64(core_api.LoginType_EmailLoginType),
		AppId:    req.Email,
		Role:     int64(core_api.RoleType_UserRoleType),
		Password: req.Password,
	})
	if err != nil {
		return resp, err
	}
	userId := createAuthResp.UserId
	if err = s.UserInit(ctx, createAuthResp.UserId, req.Name, consts.SexMan, ""); err != nil {
		return resp, err
	}
	resp.ShortToken, resp.LongToken, err = generateShortLongToken(s.Config.Auth.SecretKey, userId, s.Config.Auth.ShortTokenExpire, s.Config.Auth.LongTokenExpire)
	if err != nil {
		return resp, consts.ErrAuthentication
	}
	resp.UserId = userId
	return resp, nil
}

func (s *AuthService) UserInit(ctx context.Context, UserId, Name string, Sex int64, Url string) error {
	if _, err := s.CloudMindContent.CreateUser(ctx, &content.CreateUserReq{
		UserId: UserId,
		Name:   Name,
		Sex:    Sex,
		Url:    Url,
	}); err != nil {
		return err
	}

	if _, err := s.CloudMindTrade.CreateBalance(ctx, &trade.CreateBalanceReq{
		UserId: UserId,
	}); err != nil {
		return err
	}

	if _, err := s.CloudMindSystem.CreateNotificationCount(ctx, &system.CreateNotificationCountReq{
		UserId: UserId,
	}); err != nil {
		return err
	}
	if _, err := s.CloudMindContent.CreateHot(ctx, &content.CreateHotReq{
		HotId: UserId,
	}); err != nil {
		return err
	}

	data, _ := sonic.Marshal(&message.CreateItemMessage{
		ItemId:   UserId,
		Category: core_api.Category_name[int32(core_api.Category_UserCategory)],
	})
	if err := s.CreateItemKq.Push(pconvertor.Bytes2String(data)); err != nil {
		return err
	}
	return nil
}

func generateShortLongToken(secretKey, userId string, shortTokenExpire, longTokenExpire int64) (shortToken, longToken string, err error) {
	shortToken, _, err = generateJwtToken(userId, secretKey, shortTokenExpire)
	if err != nil {
		return "", "", err
	}
	longToken, _, err = generateJwtToken(userId, secretKey, longTokenExpire)
	if err != nil {
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
			EmailOptions: &sts.EmailOptions{Email: req.Email},
		},
		Password: req.Password,
	}); err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *AuthService) SetPasswordByPassword(ctx context.Context, req *core_api.SetPasswordByPasswordReq) (resp *core_api.SetPasswordByPasswordResp, err error) {
	resp = new(core_api.SetPasswordByPasswordResp)
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || userData.GetUserId() == "" {
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

func (s *AuthService) AskUploadAvatar(ctx context.Context, req *core_api.AskUploadAvatarReq) (resp *core_api.AskUploadAvatarResp, err error) {
	resp = new(core_api.AskUploadAvatarResp)
	userData, err := adaptor.ExtractUserMeta(ctx)
	if err != nil || userData.GetUserId() == "" {
		return resp, consts.ErrNotAuthentication
	}
	genCosStsResp, err := s.CloudMindSts.GenCosSts(ctx, &sts.GenCosStsReq{
		Path:   "users/" + req.Name,
		IsFile: false,
		Time:   req.AvatarSize / (1024 * 1024),
	})
	if err != nil {
		return resp, err
	}
	resp.SessionToken = genCosStsResp.SessionToken
	resp.TmpSecretId = genCosStsResp.SecretId
	resp.TmpSecretKey = genCosStsResp.SecretKey
	resp.StartTime = genCosStsResp.StartTime
	resp.ExpiredTime = genCosStsResp.ExpiredTime
	return resp, nil
}
