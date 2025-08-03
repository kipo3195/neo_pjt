package consts

import "errors"

var ErrTokenParsing = errors.New("token parsing failed")
var ErrInvalidClaims = errors.New("invalid token or claims")
var ErrTokenExpired = errors.New("token has expired")

var ErrFileSizeExceeded = errors.New("file size exceeded")
var ErrFileExtentionDetect = errors.New("file extension detect failed")

var ErrServerApiCallError = errors.New("server api call error")
