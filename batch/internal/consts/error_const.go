package consts

import "errors"

var ErrOrgInfoFileSendError = errors.New("org info file send error")
var ErrUserDetailFileSendError = errors.New("user detail file send error")

// message gRPC error
var ErrSendFileInfoError = errors.New("send file info error")
