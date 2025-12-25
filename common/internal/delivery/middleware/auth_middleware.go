package middleware

import (
	"common/internal/consts"
	"common/internal/delivery/middleware/claims"
	"common/internal/infrastructure/config"
	commonConsts "common/pkg/consts"
	"common/pkg/response"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserHash string `json:"userHash"`
	jwt.RegisteredClaims
}

// 토큰 생성시 사용한 key와 동일해야함.

func AuthMiddleware(tokenConfig config.TokenHashConfig) gin.HandlerFunc {

	return func(c *gin.Context) {

		// 토큰 추출
		tokenStr, err := extractTokenFromHeader(c.Request.Header)
		if err != nil {
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_105, commonConsts.E_105_MSG)
			c.Abort() // 다음 핸들러 중단
			return
		}

		// 토큰 검증
		id, hash, err := verifyJWT(tokenStr, tokenConfig)
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

		// handler에서 값을 꺼낼 수 있게 하려면
		c.Set(consts.USER_ID, id)
		c.Set(consts.USER_HASH, hash)

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

func verifyJWT(tokenStr string, tokenHash config.TokenHashConfig) (string, string, error) {

	parser := jwt.NewParser(jwt.WithoutClaimsValidation())
	log.Println("[AuthMiddleware] tokenHash : ", tokenHash)

	token, err := parser.ParseWithClaims(tokenStr, &claims.DeviceJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenHash.AccessTokenHash), nil
	})

	if err != nil {
		return "", "", fmt.Errorf("token parsing failed: %w", err)
	}

	// 유효성 체크
	parsedClaims, ok := token.Claims.(*claims.DeviceJWTClaims) // ← 여기서도 JWTClaims로
	if !ok || !token.Valid {
		return "", "", consts.ErrInvalidClaims
	}

	// 만료 시간 검증
	if parsedClaims.ExpiresAt != nil && parsedClaims.ExpiresAt.Time.Before(time.Now()) {
		return "", "", consts.ErrTokenExpired
	}

	log.Println("[AuthMiddleware] verifyJWT id : ", parsedClaims.Id)

	return parsedClaims.Id, parsedClaims.Hash, nil
}
