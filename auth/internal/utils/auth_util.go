package utils

import (
	"auth/internal/claims"
	"auth/internal/config"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// pkg/utils: 비즈니스 로직과 관련된 유틸

// 도메인 공통 로직
// 비즈니스 규칙이 포함된 경우

type AuthUtil struct {
	jwtCfg *config.JWTConfig
}

func NewAuthUtil(jwtCfg *config.JWTConfig) *AuthUtil {
	return &AuthUtil{
		jwtCfg: jwtCfg,
	}
}

func (a *AuthUtil) GenerateDeviceTokenJWT(appTokenExp int, uuid string) (string, error) {
	now := time.Now()
	// issuer
	const issuer = "auth"

	// Access 토큰 유효기간 설정
	accExpTime := now.Add(time.Duration(appTokenExp) * 24 * time.Hour)

	log.Println("jwt 토큰 생성  1 :", accExpTime)

	// 사용자 정보 포함한 Claims 생성
	uuidClaim := &claims.DeviceJWTClaims{
		Uuid: uuid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accExpTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    issuer,
		},
	}

	// 토큰 생성 (HS256 사용)
	accToken := jwt.NewWithClaims(jwt.SigningMethodHS256, uuidClaim)

	secret := []byte("neo-test-secret-key")

	// 서명 및 문자열 반환
	token, err := accToken.SignedString(secret)

	log.Println("jwt 토큰 생성  2 :", token)
	return token, err
}
