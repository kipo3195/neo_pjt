package consts

import "errors"

var ErrProfileImgSaveError = errors.New("profile image save error")
var ErrProfileImgDBSaveError = errors.New("profile image db save error")
var ErrProfileImgNotRegist = errors.New("profile image not regist")
var ErrProfileImgDBDeleteError = errors.New("profile image db delete error")
var ErrProfileImgDBRoleBackError = errors.New("profile image db rollback error")
var ErrProfileImgNotExist = errors.New("profile image not exist")
var ErrProfileImgRemoveError = errors.New("profile image remove error")

var ErrFileSizeExceeded = errors.New("file size exceeded")
var ErrFileExtentionDetect = errors.New("file extension detect failed")
var ErrFileExtentionInvalid = errors.New("file extension invalid")

// token
var ErrTokenParsing = errors.New("token parsing failed")
var ErrTokenExpired = errors.New("token has expired")
var ErrTokenSignatureInvalid = errors.New("token signature invalid")
var ErrInvalidClaims = errors.New("invalid token or claims")
var ErrTokenMismatch = errors.New("token mismatch")
var ErrRefreshTokenAuthError = errors.New("refreshToken authentication failed")

// batch
var ErrUnzipAndGetJSONError = errors.New("unzip and get json error")
var ErrInvalidUserDetailJSONError = errors.New("invalid user detail json error")
