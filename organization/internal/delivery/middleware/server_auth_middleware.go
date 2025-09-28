package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"org/internal/consts"
	commonConsts "org/pkg/consts"
	"org/pkg/response"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var serverJwtSecretKey = []byte("neo-test-secret-key")

func ServerAuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		// 토큰 추출 및 검증 로직
		tokenStr, err := serverExtractTokenFromHeader(c.Request.Header)
		if err != nil {
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_105, commonConsts.E_105_MSG)
			c.Abort() // 다음 핸들러 중단
			return
		}

		_, err = serverVerifyJWT(tokenStr)
		if err != nil {
			log.Println(err, err.Error())
			if errors.Is(err, consts.ErrTokenExpired) {
				log.Println("토큰 만료")
				response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_107, commonConsts.E_107_MSG)
			} else {
				log.Println("토큰 검증 실패")
				response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_106, commonConsts.E_106_MSG)
				c.Abort() // 다음 핸들러 중단
			}
			return
		}

		c.Next()
	}
}

// Authorization 헤더에서 "Bearer <token>" 형태로 된 토큰 추출
func serverExtractTokenFromHeader(header http.Header) (string, error) {
	authHeader := header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is missing")
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("authorization header format must be Bearer {token}")
	}
	return parts[1], nil
}

func serverVerifyJWT(tokenStr string) (bool, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 서명 알고리즘 확인
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// 시크릿 키 반환
		return serverJwtSecretKey, nil
	})

	if err != nil {
		return false, fmt.Errorf("token parsing failed: %w", err)
	}

	// 유효성 체크
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return false, errors.New("invalid token or claims")
	}

	// 만료 시간 검증
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return false, errors.New("token has expired")
	}

	return true, nil
}
