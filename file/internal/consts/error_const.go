package consts

import "errors"

// db
var ErrDBresultNotFound = errors.New("db result not found")
var ErrDB = errors.New("db error")
var ErrDBResultNotUpdate = errors.New("db result not update")

// token
var ErrTokenParsing = errors.New("token parsing failed")
var ErrTokenExpired = errors.New("token has expired")
var ErrTokenSignatureInvalid = errors.New("token signature invalid")
var ErrInvalidClaims = errors.New("invalid token or claims")
var ErrTokenMismatch = errors.New("token mismatch")
var ErrRefreshTokenAuthError = errors.New("refreshToken authentication failed")
