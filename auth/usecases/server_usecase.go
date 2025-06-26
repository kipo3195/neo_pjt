package usecases

import (
	"auth/claims"
	consts "auth/consts"
	svDto "auth/dto/server"
	"auth/entities"
	"auth/repositories"
	"context"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type serverUsecase struct {
	repo     repositories.ServerRepository
	authRepo repositories.AuthRepository
}

type ServerUsecase interface {
	AppTokenValidation(req svDto.SvAppTokenValidationRequest, ctx context.Context) (bool, error)
}

func NewServerUsecase(repo repositories.ServerRepository, authRepo repositories.AuthRepository) ServerUsecase {
	return &serverUsecase{
		repo:     repo,
		authRepo: authRepo,
	}
}

func (r *serverUsecase) AppTokenValidation(req svDto.SvAppTokenValidationRequest, ctx context.Context) (bool, error) {

	// authUsecase를 주입받아 사용.
	flag, err := r.authRepo.GetValidation(toAppTokenValidationEntity(req.Uuid, req.AppToken))
	if err != nil {
		switch {
		case errors.Is(err, consts.ErrDbRowNotFound):
			// 매핑된 hash 정보가 없음
			return false, err
		default:
			// 기타 DB 에러
			return false, err
		}
	}

	// 토큰 정보 불일치
	if !flag {
		return false, consts.ErrInvalidClaims
	}

	err = appTokenValidationCheck(req.AppToken)

	if err != nil {
		// 만료 에러 확인
		if errors.Is(err, jwt.ErrTokenExpired) {
			fmt.Println("토큰이 만료되었습니다.")
			return false, consts.ErrTokenExpired
		}
		// 서명 오류
		if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			fmt.Println("서명이 유효하지 않습니다.")
			return false, consts.ErrTokenSignatureInvalid
		}
		// 그 외 오류
		fmt.Println("토큰 파싱 오류:", err)
		return false, consts.ErrTokenParsing
	}

	return true, nil
}

func toAppTokenValidationEntity(uuid string, appToken string) entities.AppTokenValidationEntity {
	return entities.AppTokenValidationEntity{
		Uuid:     uuid,
		AppToken: appToken,
	}
}

func appTokenValidationCheck(appToken string) error {
	// 파싱하면서 검증 진행
	token, err := jwt.ParseWithClaims(appToken, &claims.DeviceJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// HS256 서명 방식 검증
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("neo-test-secret-key"), nil
	})

	// 에러 핸들링
	// 에러 처리
	if err != nil {
		return err
	}

	// 유효한 토큰인 경우
	if claims, ok := token.Claims.(*claims.DeviceJWTClaims); ok && token.Valid {
		fmt.Println("토큰이 유효합니다.")
		fmt.Println("UUID:", claims.Uuid)
		fmt.Println("만료 시간:", claims.ExpiresAt.Time)
		return nil
	} else {
		fmt.Println("토큰이 유효하지 않습니다.")
		return errors.New(consts.ErrInvalidClaims.Error())
	}
}
