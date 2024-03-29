package adaptor

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/go-pkg/utils/pconvertor"
	"github.com/CloudStriver/go-pkg/utils/util"
	"github.com/CloudStriver/go-pkg/utils/util/log"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/basic"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/golang-jwt/jwt/v4"
)

type CommonString string

const hertzContext CommonString = "hertzContext"

func InjectContext(ctx context.Context, c *app.RequestContext) context.Context {
	return context.WithValue(ctx, hertzContext, c)
}

func ExtractContext(ctx context.Context) (*app.RequestContext, error) {
	c, ok := ctx.Value(hertzContext).(*app.RequestContext)
	if !ok {
		return nil, errors.New("hertz context not found")
	}
	return c, nil
}

func ExtractUserMeta(ctx context.Context) (user *basic.UserMeta, err error) {
	user = new(basic.UserMeta)
	c, err := ExtractContext(ctx)
	if err != nil {
		return
	}
	tokenString := c.GetHeader("Authorization")
	if pconvertor.Bytes2String(tokenString) == "" {
		return
	}
	token, err := jwt.Parse(pconvertor.Bytes2String(tokenString), func(_ *jwt.Token) (interface{}, error) {
		return jwt.ParseECPublicKeyFromPEM([]byte(config.GetConfig().Auth.PublicKey))
	})
	if err != nil {
		return
	}
	if !token.Valid {
		err = errors.New("token is not valid")
		return
	}
	data, err := json.Marshal(token.Claims)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, user)
	if err != nil {
		return
	}
	if user.SessionUserId == "" {
		user.SessionUserId = user.UserId
	}
	if user.SessionAppId == 0 {
		user.SessionAppId = user.AppId
	}
	if user.SessionDeviceId == "" {
		user.SessionDeviceId = user.DeviceId
	}
	log.CtxInfo(ctx, "userMeta=%s", util.JSONF(user))
	return
}

func ExtractExtra(ctx context.Context) (extra *basic.Extra) {
	extra = new(basic.Extra)
	var err error
	defer func() {
		if err != nil {
			log.CtxInfo(ctx, "extract extra fail, err=%v", err)
		}
	}()
	c, err := ExtractContext(ctx)
	if err != nil {
		return
	}
	extra.ClientIP = c.ClientIP()
	err = c.Bind(extra)
	if err != nil {
		return
	}
	log.CtxInfo(ctx, "extra=%s", util.JSONF(extra))
	return
}
