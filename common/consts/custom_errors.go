package consts

import "errors"

var ErrInvalidType = errors.New("invalid type")
var ErrDB = errors.New("DB error")

var ErrTokenParsing = errors.New("token parsing failed")
var ErrInvalidClaims = errors.New("invalid token or claims")
var ErrTokenExpired = errors.New("token has expired")

var ErrServerError = errors.New("server error")

var ErrSkinHashInvalid = errors.New("skin hash invalid")
var ErrConfigHashInvalid = errors.New("config hash invalid")

var ErrFileSizeExceeded = errors.New("file size exceeded")
var ErrFileExtentionDetect = errors.New("file extension detect failed")
var ErrFileExtentionInvalid = errors.New("file extension invalid")
