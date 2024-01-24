package consts

import (
	"google.golang.org/grpc/codes"
)

var PageSize int64 = 10
var NotContent = codes.Code(204)

const (
	RelationUserType    = 1
	RelationFileType    = 2
	RelationProductType = 3
	RelationPostType    = 4
	RelationLikeType    = 1 // 点赞
	RelationFollowType  = 2 // 关注
	RelationCollectType = 3 // 收藏
	RelationViewType    = 4 // 浏览
)
