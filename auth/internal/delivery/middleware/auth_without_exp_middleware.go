package middleware

import (
	"auth/internal/consts"
	"auth/internal/delivery/middleware/claims"
	"auth/internal/infrastructure/config"
	commonConsts "auth/pkg/consts"
	"auth/pkg/response"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// 토큰 생성시 사용한 key와 동일해야함.

func AuthWithoutExpMiddleware(tokenConfig config.TokenHashConfig) gin.HandlerFunc {

	return func(c *gin.Context) {

		// 토큰 추출
		tokenStr, err := extractTokenFromHeader(c.Request.Header)
		if err != nil {
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_105, commonConsts.E_105_MSG)
			c.Abort() // 다음 핸들러 중단
			return
		}

		// 토큰 검증
		id, err := verifyJWTWithoutExp(tokenStr, tokenConfig)
		if err != nil {
			log.Println(err, err.Error())
			log.Println("토큰 검증 실패")
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_106, commonConsts.E_106_MSG)
			c.Abort() // 다음 핸들러 중단
			return
		}

		// handler에서 값을 꺼낼 수 있게 하려면
		c.Set(consts.USER_ID, id)

		// 정상 처리 = 검증 성공 → 다음 핸들러 호출
		c.Next()
	}
}

func verifyJWTWithoutExp(tokenStr string, tokenHash config.TokenHashConfig) (string, error) {

	parser := jwt.NewParser(jwt.WithoutClaimsValidation())
	log.Println("111 tokenHash : ", tokenHash)

	token, err := parser.ParseWithClaims(tokenStr, &claims.DeviceJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenHash.AccessTokenHash), nil
	})

	log.Println("222")

	if err != nil {
		return "", fmt.Errorf("token parsing failed: %w", err)
	}

	// 유효성 체크
	parsedClaims, ok := token.Claims.(*claims.DeviceJWTClaims) // ← 여기서도 JWTClaims로
	if !ok || !token.Valid {
		return "", consts.ErrInvalidClaims
	}

	log.Println("토큰 검증 완료 verifyJWTWithoutExp id : ", parsedClaims.Id)

	return parsedClaims.Id, nil
}
