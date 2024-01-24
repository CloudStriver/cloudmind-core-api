package consts

import (
	"google.golang.org/grpc/codes"
)

var PageSize int64 = 10
var NotContent = codes.Code(204)

const (
	RelationUserType    = 1
	RelationPostType    = 4
	RelationLikeType    = 1 // 点赞
	RelationCollectType = 3 // 收藏
	RelationViewType    = 4 // 浏览
	PassCheckEmail      = "PassCheckEmail"
)
