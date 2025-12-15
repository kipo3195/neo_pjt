package consts

import "errors"

var ErrInvalidClaims = errors.New("invalid token or claims")
var ErrTokenExpired = errors.New("token has expired")

var ErrUnzipAndGetJSONError = errors.New("unzip and get json error")
var ErrInvalidOrgJSONError = errors.New("invalid org json error")

var ErrOrgCodeNotExist = errors.New("org code not exist")

var ErrOrgFileNotFound = errors.New("org file not found")
