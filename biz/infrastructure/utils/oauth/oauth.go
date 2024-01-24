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
	TokenType   string `json:"token_type"` // 这个字段没用到
	Scope       string `json:"scope"`      // 这个字段也没用到
}

type UserInfo struct {
	Name string `json:"name"`
	Id   int64  `json:"id"`
}

func OauthLogin(conf config.OauthConf, authType sts.AuthType, code string) (*UserInfo, error) {
	url := getTokenUrl(conf, authType, code)
	var req *http.Request
	var err error
	if req, err = http.NewRequest(http.MethodPost, url, nil); err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")
	var httpClient = http.Client{}
	var res *http.Response
	if res, err = httpClient.Do(req); err != nil {
		return nil, err
	}
	var token Token
	if err = json.NewDecoder(res.Body).Decode(&token); err != nil {
		return nil, err
	}
	url = getUserUrl(authType, token.AccessToken)

	if req, err = http.NewRequest(http.MethodGet, url, nil); err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")
	if authType == sts.AuthType_gitee {
		req.Header.Set("Authorization", fmt.Sprintf("token %s", token.AccessToken))
	}

	// 发送请求并获取响应
	var client = http.Client{}
	if res, err = client.Do(req); err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, consts.ErrThirdLogin
	}
	userInfo := &UserInfo{}
	if err = json.NewDecoder(res.Body).Decode(&userInfo); err != nil {
		return nil, err
	}
	return userInfo, nil
}

func getTokenUrl(conf config.OauthConf, authType sts.AuthType, code string) string {
	switch authType {
	case sts.AuthType_github:
		return fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s", conf.ClientId, conf.Secret, code)
	case sts.AuthType_gitee:
		return fmt.Sprintf("https://gitee.com/oauth/token?grant_type=authorization_code&code=%s&client_id=%s&redirect_uri=%s&client_secret=%s", code, conf.ClientId, "http://apisix.cloudmind.top/auth/giteeLogin", conf.Secret)
	default:
		return ""
	}
}

func getUserUrl(authType sts.AuthType, accessToken string) string {
	switch authType {
	case sts.AuthType_github:
		return "https://api.github.com/user"
	case sts.AuthType_gitee:
		return fmt.Sprintf("https://gitee.com/api/v5/user?access_token=%s", accessToken) // github用户信息获取接口
	default:
		return ""
	}
}
