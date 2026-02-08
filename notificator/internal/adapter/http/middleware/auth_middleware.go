package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"notificator/internal/adapter/http/middleware/claims"
	"notificator/internal/consts"
	"notificator/internal/domain/logger"
	"notificator/internal/infrastructure/config"
	commonConsts "notificator/pkg/consts"
	"notificator/pkg/response"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

// 토큰 생성시 사용한 key와 동일해야함.

func AuthMiddleware(logger logger.Logger, tokenConfig config.TokenHashConfig) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ctx := r.Context()

			// 토큰 추출
			tokenStr, err := extractTokenFromHeader(r.Header)

			if err != nil {
				response.SendError(w, commonConsts.UNAUTHORIZED, commonConsts.ERROR, commonConsts.E_105, commonConsts.E_105_MSG)
				logger.Error(ctx, "at_verification_fail",
					"detail_msg", err.Error(),
					"option", "not exist")
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
					response.SendError(w, commonConsts.UNAUTHORIZED, commonConsts.ERROR, commonConsts.E_107, commonConsts.E_107_MSG)
				} else {
					// 토큰 인증 실패 규격이 다르거나 정상적인 발급이 아님
					logger.Error(ctx, "at_verification_fail",
						"detail_msg", err.Error(),
						"option", "invalid")
					response.SendError(w, commonConsts.UNAUTHORIZED, commonConsts.ERROR, commonConsts.E_106, commonConsts.E_106_MSG)
				}
				return
			}

			// context 저장
			ctx = context.WithValue(ctx, consts.USER_ID, id)
			ctx = context.WithValue(ctx, consts.USER_HASH, hash)

			// context 적용
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})

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
