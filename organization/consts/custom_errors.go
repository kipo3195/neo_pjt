package consts

import "errors"

var (
	ErrTokenParsing  = errors.New("token parsing failed")
	ErrInvalidClaims = errors.New("invalid token or claims")
	ErrTokenExpired  = errors.New("token has expired")
)
