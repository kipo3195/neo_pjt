package claims

import "github.com/golang-jwt/jwt/v5"

type DeviceJWTClaims struct {
	Id   string `json:"id"`
	Uuid string `json:"uuid"`
	jwt.RegisteredClaims
}
