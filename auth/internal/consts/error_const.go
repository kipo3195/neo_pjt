package consts

import "errors"

// user
var ErrUserIdDoesNotExist = errors.New("user id does not exist")
var ErrUnregisteredUser = errors.New("unregistered user")

// common
var ErrInvalidType = errors.New("invalid type")
var ErrServerError = errors.New("server error")

// db
var ErrDB = errors.New("DB error")
var ErrDbRowNotFound = errors.New("db row not found")

// token
var ErrTokenParsing = errors.New("token parsing failed")
var ErrTokenExpired = errors.New("token has expired")
var ErrTokenSignatureInvalid = errors.New("token signature invalid")
var ErrInvalidClaims = errors.New("invalid token or claims")
var ErrTokenMismatch = errors.New("token mismatch")
var ErrRefreshTokenAuthError = errors.New("refreshToken authentication failed")

// uuid
var ErrUnregisteredUuid = errors.New("unregistered uuid")

// auth
var ErrAuthenticationFailed = errors.New("authentication failed")
var ErrUserAuthChallengeExpired = errors.New("user auth challenge expired")
var ErrUserAuthFvMismatch = errors.New("user auth fv mismatch")

// salt
var ErrSaltNotRegist = errors.New("salt is not regist")

// device
var ErrDeviceNotRegist = errors.New("device not regist")
var ErrDeviceChallengeExpired = errors.New("device challenge expired")
var ErrDeviceChallengeMismatch = errors.New("device challenge mismatch")
var ErrDeviceAlreadyRegist = errors.New("device already regist")

// parsing
var ErrTimeDataParsingError = errors.New("time data parsing error")

// http status
var ErrHttpStatusError = errors.New("http status is not ok")
