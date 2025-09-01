package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"org/internal/delivery/consts"
	commonConsts "org/pkg/consts"
	"org/pkg/response"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecretKey = []byte("neo-test-secret-key")

type JWTClaims struct {
	UserHash string `json:"userHash"`
	jwt.RegisteredClaims
}

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		// 토큰 추출
		tokenStr, err := extractTokenFromHeader(c.Request.Header)
		if err != nil {
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_105, commonConsts.E_105_MSG)
			c.Abort() // 다음 핸들러 중단
			return
		}

		// 토큰 검증
		_, err = verifyJWT(tokenStr)
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

		// 정상 처리 = 검증 성공 → 다음 핸들러 호출
		c.Next()
	}
}

// Authorization 헤더에서 "Bearer <token>" 형태로 된 토큰 추출
func extractTokenFromHeader(header http.Header) (string, error) {
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

func verifyJWT(tokenStr string) (string, error) {

	parser := jwt.NewParser(jwt.WithoutClaimsValidation())

	token, err := parser.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		return "", fmt.Errorf("token parsing failed: %w", err)
	}

	// 유효성 체크
	parsedClaims, ok := token.Claims.(*JWTClaims) // ← 여기서도 JWTClaims로
	if !ok || !token.Valid {
		return "", consts.ErrInvalidClaims
	}

	// 만료 시간 검증
	if parsedClaims.ExpiresAt != nil && parsedClaims.ExpiresAt.Time.Before(time.Now()) {
		return "", consts.ErrTokenExpired
	}

	return parsedClaims.UserHash, nil
}
