package consts

import "errors"

var ErrInvalidType = errors.New("invalid type")
var ErrDB = errors.New("DB error")
var ErrServerError = errors.New("server error")

var ErrDbRowNotFound = errors.New("db row not found")

var ErrTokenParsing = errors.New("token parsing failed")
var ErrTokenExpired = errors.New("token has expired")
var ErrTokenSignatureInvalid = errors.New("token signature invalid")
var ErrInvalidClaims = errors.New("invalid token or claims")

var ErrUnregisteredUuid = errors.New("unregistered uuid")
var ErrTokenMismatch = errors.New("token mismatch")
var ErrUnregisteredUser = errors.New("unregistered user")
var ErrAuthenticationFailed = errors.New("authentication failed")

var ErrSaltNotRegist = errors.New("salt is not regist.")

var ErrUserAuthChallengeExpired = errors.New("user auth challenge expired")
