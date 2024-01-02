package consts

import (
	"errors"
	"google.golang.org/grpc/status"
)

var ErrNotAuthentication = errors.New("not authentication")
var ErrForbidden = errors.New("forbidden")

var (
	ErrAuthentication = status.Error(10001, "生成token失败")
	ErrParseToken     = status.Error(10002, "解析token失败")
	ErrNotLongToken   = status.Error(10003, "请使用长token刷新")
	ErrThirdLogin     = status.Error(10004, "第三方登录失败")
)

var (
	ErrPasswordNotEqual = status.Error(20001, "密码错误")
	ErrCodeNotFound     = status.Error(20002, "验证码已过期")
	ErrCodeNotEqual     = status.Error(20003, "验证码错误")
	ErrHaveExist        = status.Error(20004, "邮箱已被注册")
)
