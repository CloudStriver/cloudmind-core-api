package oauth

import (
	"encoding/json"
	"fmt"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/consts"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/sts"
	"net/http"
)

type Token struct {
	AccessToken string `json:"access_token"`
	//TokenType   string `json:"token_type"` // 这个字段没用到
	//Scope       string `json:"scope"`      // 这个字段也没用到
}

type GiteeInfo struct {
	Name      string `json:"name"`
	Id        int64  `json:"id"`
	AvatarUrl string `json:"avatar_url"`
}

type QQInfo struct {
	Ret             int    `json:"ret"`
	Msg             string `json:"msg"`
	IsLost          int    `json:"is_lost"`
	Nickname        string `json:"nickname"`
	Gender          string `json:"gender"`
	GenderType      int    `json:"gender_type"`
	Province        string `json:"province"`
	City            string `json:"city"`
	Year            string `json:"year"`
	Figureurl       string `json:"figureurl"`
	Figureurl1      string `json:"figureurl_1"`
	Figureurl2      string `json:"figureurl_2"`
	FigureurlQq1    string `json:"figureurl_qq_1"`
	FigureurlQq2    string `json:"figureurl_qq_2"`
	FigureurlQq     string `json:"figureurl_qq"`
	IsYellowVip     string `json:"is_yellow_vip"`
	Vip             string `json:"vip"`
	YellowVipLevel  string `json:"yellow_vip_level"`
	Level           string `json:"level"`
	IsYellowYearVip string `json:"is_yellow_year_vip"`
	ClientId        string `json:"client_id"`
	Openid          string `json:"openid"`
}

type WechatInfo struct {
	Code int         `json:"code"`
	Msg  interface{} `json:"msg"`
	Data struct {
		QrUrl      string `json:"qrUrl"`
		TempUserId string `json:"tempUserId"`
	} `json:"data"`
}

func QQLogin(conf config.OauthConf, code string) (*QQInfo, error) {
	url := getTokenUrl(conf, sts.AuthType_qq, code)
	var (
		req        *http.Request
		err        error
		httpClient http.Client
		res        *http.Response
		token      Token
	)
	if req, err = http.NewRequest(http.MethodGet, url, nil); err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")
	if res, err = httpClient.Do(req); err != nil {
		return nil, err
	}

	if err = json.NewDecoder(res.Body).Decode(&token); err != nil {
		return nil, err
	}
	url = fmt.Sprintf("https://graph.qq.com/oauth2.0/me?access_token=%s&fmt=json", token.AccessToken)
	if req, err = http.NewRequest(http.MethodGet, url, nil); err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")
	if res, err = httpClient.Do(req); err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, consts.ErrThirdLogin
	}
	userInfo := &QQInfo{}
	if err = json.NewDecoder(res.Body).Decode(&userInfo); err != nil {
		return nil, err
	}
	url = fmt.Sprintf("https://graph.qq.com/user/get_user_info?access_token=%s&oauth_consumer_key=%s&openid=%s&fmt=json", token.AccessToken, conf.ClientId, userInfo.Openid)
	if req, err = http.NewRequest(http.MethodGet, url, nil); err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")
	if res, err = httpClient.Do(req); err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, consts.ErrThirdLogin
	}
	if err = json.NewDecoder(res.Body).Decode(&userInfo); err != nil {
		return nil, err
	}
	return userInfo, nil
}

func getTokenUrl(conf config.OauthConf, authType sts.AuthType, code string) string {
	switch authType {
	case sts.AuthType_qq:
		return fmt.Sprintf("https://graph.qq.com/oauth2.0/token?grant_type=authorization_code&client_id=%s&client_secret=%s&code=%s&redirect_uri=%s&fmt=json", conf.ClientId, conf.Secret, code, conf.Redirect)
	default:
		return ""
	}
}
