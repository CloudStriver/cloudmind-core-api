package consts

import (
	"errors"
	"google.golang.org/grpc/status"
)

var ErrNotAuthentication = errors.New("not authentication")
var ErrForbidden = errors.New("forbidden")

var (
	ErrAuthentication         = status.Error(10001, "生成token失败")
	ErrParseToken             = status.Error(10002, "解析token失败")
	ErrNotLongToken           = status.Error(10003, "请使用长token刷新")
	ErrThirdLogin             = status.Error(10004, "第三方登录失败")
	ErrNotSupportRelationType = status.Error(10005, "不支持的参数类型")
)

var (
	ErrNotEmailCheck = status.Error(20006, "请先通过邮箱验证")
	ErrEmailNotFound = status.Error(20007, "邮箱不存在")
	ErrFlowNotEnough = status.Error(20008, "流量不足")
)

var (
	ErrFileIsNotDir      = status.Error(30001, "目标文件不是文件夹")
	ErrNoAccessFile      = status.Error(30002, "您无权访问该文件")
	ErrFileNotExist      = status.Error(30003, "查询的文件不存在")
	ErrIllegalOperation  = status.Error(30004, "非法操作")
	ErrShareFileKey      = status.Error(10110, "分享文件提取码错误")
	ErrShareCodeNotExist = status.Error(10111, "分享文件链接已失效")
	ErrNoAuditStatus     = status.Error(30002, "该文件无需审核")
)
