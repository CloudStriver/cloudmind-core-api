package consts

import (
	"google.golang.org/grpc/codes"
	"time"
)

var PageSize int64 = 10
var NotContent = codes.Code(204)
var NotDel int64 = 1
var SoftDel int64 = 2
var HardDel int64 = 3
var InitBrowseNumber int64 = 0

const (
	RelationUserType    = 1
	RelationFileType    = 2
	RelationProductType = 3
	RelationPostType    = 4
	RelationLikeType    = 1 // 点赞
	RelationFollowType  = 2 // 关注
	RelationCollectType = 3 // 收藏
	RelationViewType    = 4 // 浏览
	PassCheckEmail      = "PassCheckEmail"
	BatcherSize         = 100
	BatcherBuffer       = 100
	BatcherWorker       = 10
	BatcherInterval     = time.Second
	PostPublicStatus    = int64(1)
	PostPrivateStatus   = int64(2)
	NotificationRead    = true
	NotificationNotRead = false
	ObjectFile          = int64(1)
	FolderSize          = int64(-1)
)
