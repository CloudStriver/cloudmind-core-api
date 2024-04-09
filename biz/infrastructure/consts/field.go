package consts

import (
	"google.golang.org/grpc/codes"
	"time"
)

const (
	PassCheckEmail            = "PassCheckEmail"
	BatcherSize               = 100
	BatcherBuffer             = 100
	BatcherWorker             = 10
	BatcherInterval           = time.Second
	NotificationRead          = true
	NotificationNotRead       = false
	ObjectFile                = int64(1)
	FolderSize                = int64(-1)
	PublicSlider              = 1
	PrivateSlider             = 2
	UserRankKey               = "cache:rank:user"
	FileRankKey               = "cache:rank:file"
	PostRankKey               = "cache:rank:post"
	BloomRelationKey          = "cache:bloom:relation"
	SexMan                    = 1
	SexWoman                  = 2
	WechatLoginKey            = "cache:wechat"
	ViewCountKey              = "cache:view"
	PageSize            int64 = 10
	NotContent                = codes.Code(204)
	NotDel              int64 = 1
	SoftDel             int64 = 2
	HardDel             int64 = 3
	InitNumber          int64 = 0
)
