package consts

import (
	"errors"
	"google.golang.org/grpc/codes"
)

var ErrNotAuthentication = errors.New("not authentication")
var ErrForbidden = errors.New("forbidden")
var PageSize int64 = 10
var NotContent = codes.Code(204)
