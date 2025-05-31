package consts

import "errors"

var ErrInvalidType = errors.New("invalid type")
var ErrInvalidMappingServer = errors.New("invalid mapping server")
var ErrDB = errors.New("DB error")

var ErrTokenParsing = errors.New("token parsing failed")
var ErrInvalidClaims = errors.New("invalid token or claims")
var ErrTokenExpired = errors.New("token has expired")
