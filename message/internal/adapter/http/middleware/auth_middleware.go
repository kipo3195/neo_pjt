package middleware

import (
	"context"
	"errors"
	"fmt"
	"message/internal/consts"
	"message/internal/delivery/middleware/claims"
	"message/internal/domain/logger"
	"message/internal/infrastructure/config"
	commonConsts "message/pkg/consts"
	"message/pkg/response"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// 토큰 생성시 사용한 key와 동일해야함.

func AuthMiddleware(tokenConfig config.TokenHashConfig, logger logger.Logger) gin.HandlerFunc {

	return func(c *gin.Context) {

		// 이전 미들웨어에서 주입된 context
		ctx := c.Request.Context()

		// 토큰 추출
		tokenStr, err := extractTokenFromHeader(c.Request.Header)
		if err != nil {
			response.SendError(c, commonConsts.UNAUTHORIZED, commonConsts.ERROR, commonConsts.E_105, commonConsts.E_105_MSG)
			logger.Error(ctx, "at_verification_fail",
				"detail_msg", err.Error(),
				"option", "not exist")
			c.Abort() // 다음 핸들러 중단
			return
		}

		// 토큰 검증
		id, hash, err := verifyJWT(tokenStr, tokenConfig)
		if err != nil {
			if errors.Is(err, consts.ErrTokenExpired) {
				// 토큰 만료 ..
				logger.Error(ctx, "at_verification_fail",
					"detail_msg", err.Error(),
					"option", "expired")
				response.SendError(c, commonConsts.UNAUTHORIZED, commonConsts.ERROR, commonConsts.E_107, commonConsts.E_107_MSG)
			} else {
				// 토큰 인증 실패 규격이 다르거나 정상적인 발급이 아님
				logger.Error(ctx, "at_verification_fail",
					"detail_msg", err.Error(),
					"option", "invalid")
				response.SendError(c, commonConsts.UNAUTHORIZED, commonConsts.ERROR, commonConsts.E_106, commonConsts.E_106_MSG)
				c.Abort() // 다음 핸들러 중단
			}
			return
		}

		// 기존 ctx를 감싸서 user_hash 추가
		ctx = context.WithValue(ctx, "user_hash", hash)

		// 다시 request에 주입
		c.Request = c.Request.WithContext(ctx)

		// handler에서 값을 꺼낼 수 있게 하려면
		c.Set(consts.USER_ID, id)
		c.Set(consts.USER_HASH, hash)

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

	return parsedClaims.Id, parsedClaims.Hash, nil
}
