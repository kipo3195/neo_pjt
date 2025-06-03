package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"org/claims"
	"org/consts"
	dto "org/dto/common"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecretKey = []byte("neo-test-secret-key")

func AuthMiddleware(next http.Handler) http.Handler {

	// response
	var res dto.Response

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 토큰 추출 및 검증 로직
		tokenStr, err := extractTokenFromHeader(r.Header)
		if err != nil {
			fmt.Println("토큰 전달 형식이 맞지 않음. header check.")
			res.Result = consts.ERROR
			res.Data = dto.ErrorResponse{
				Code:    consts.E_105,
				Message: consts.E_105_MSG,
			}
			w.WriteHeader(http.StatusUnauthorized) //401
			json.NewEncoder(w).Encode(res)
			return
		}

		_, err = verifyJWT(tokenStr)
		if err != nil {
			fmt.Println(err, err.Error())
			if errors.Is(err, consts.ErrTokenExpired) {
				fmt.Println("토큰 만료")
				res.Result = consts.ERROR
				res.Data = dto.ErrorResponse{
					Code:    consts.E_107,
					Message: consts.E_107_MSG,
				}

			} else {
				fmt.Println("토큰 검증 실패")
				res.Result = consts.ERROR
				res.Data = dto.ErrorResponse{
					Code:    consts.E_106,
					Message: consts.E_106_MSG,
				}
			}

			w.WriteHeader(http.StatusUnauthorized) //401
			json.NewEncoder(w).Encode(res)
			return

		}

		next.ServeHTTP(w, r)
	})
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

func verifyJWT(tokenStr string) (bool, error) {

	parser := jwt.NewParser(jwt.WithoutClaimsValidation())

	token, err := parser.ParseWithClaims(tokenStr, &claims.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		return false, fmt.Errorf("token parsing failed: %w", err)
	}

	// 유효성 체크
	parsedClaims, ok := token.Claims.(*claims.JWTClaims) // ← 여기서도 JWTClaims로
	if !ok || !token.Valid {
		return false, consts.ErrInvalidClaims
	}

	// 만료 시간 검증
	if parsedClaims.ExpiresAt != nil && parsedClaims.ExpiresAt.Time.Before(time.Now()) {
		return false, consts.ErrTokenExpired
	}

	return true, nil
}
