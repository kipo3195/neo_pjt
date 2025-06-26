package consts

import "errors"

var ErrInvalidType = errors.New("invalid type")
var ErrDB = errors.New("DB error")

var ErrTokenParsing = errors.New("token parsing failed")
var ErrInvalidClaims = errors.New("invalid token or claims")
var ErrTokenExpired = errors.New("token has expired")

var ErrDbRowNotFound = errors.New("db row not found")

var ErrTokenSignatureInvalid = errors.New("token signature invalid")

var ErrTokenMismatch = errors.New("token mismatch")
