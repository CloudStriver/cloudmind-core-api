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
var InitNumber int64 = 0

const (
	PassCheckEmail      = "PassCheckEmail"
	BatcherSize         = 100
	BatcherBuffer       = 100
	BatcherWorker       = 10
	BatcherInterval     = time.Second
	NotificationRead    = true
	NotificationNotRead = false
	ObjectFile          = int64(1)
	FolderSize          = int64(-1)
	PublicSlider        = 1
	PrivateSlider       = 2
)
