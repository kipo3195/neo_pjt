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

// 앱 토큰 재발급 API에서 refreshToken이 검증되지 않았을때 (일치하지 않을때)
var ErrRefreshTokenAuthInvalid = errors.New("refresh token auth invalid")

// 앱 토큰 재발급 API에서 refreshToken이 만료되었을때
var ErrRefreshTokenAuthExpired = errors.New("refresh token expired")

var ErrProfileImgSaveError = errors.New("profile image save error")
var ErrProfileImgDBSaveError = errors.New("profile image db save error")
