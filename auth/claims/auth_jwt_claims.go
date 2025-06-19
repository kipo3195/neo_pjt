package claims

import "github.com/golang-jwt/jwt/v5"

type AuthJWTClaims struct {
	UserHash string `json:"userHash"`
	jwt.RegisteredClaims
}
