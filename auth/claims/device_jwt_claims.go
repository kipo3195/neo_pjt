package claims

import "github.com/golang-jwt/jwt/v5"

type DeviceJWTClaims struct {
	Uuid string `json:"uuid"`
	jwt.RegisteredClaims
}
