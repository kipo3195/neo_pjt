package usecase

import (
	"auth/internal/application/usecase/input"
	"auth/internal/application/util"
	"auth/internal/claims"
	"auth/internal/consts"
	"auth/internal/delivery/dto/token"
	"auth/internal/domain/shared"
	"auth/internal/domain/token/entity"
	"auth/internal/domain/token/repository"
	"auth/internal/infrastructure/config"

	"context"
	"errors"
	"fmt"
	"log"

	"github.com/golang-jwt/jwt/v5"
)

type tokenUsecase struct {
	repo   repository.TokenRepository
	jwtCfg *config.JWTConfig
}

type TokenUsecase interface {
	AppTokenValidation(in input.AppTokenValidationInput, ctx context.Context) (bool, error)
	GenerateAppToken(body token.GenerateAppTokenRequestBody) (*token.GenerateAppTokenResponseDTO, error)
}

func NewTokenUsecase(repo repository.TokenRepository, jwtCfg *config.JWTConfig) TokenUsecase {
	return &tokenUsecase{
		repo:   repo,
		jwtCfg: jwtCfg,
	}
}

func (r *tokenUsecase) AppTokenValidation(in input.AppTokenValidationInput, ctx context.Context) (bool, error) {

	// 결과 데이터를 미리 정의함.
	var flag bool
	var err error

	entity := entity.NewAppTokenValidationEntity(in.Uuid, in.Token)

	if in.TokenType == "appToken" {
		flag, err = r.repo.GetValidationAppToken(entity)
	} else if in.TokenType == "accessToken" {
		// accessToken은 DB에 저장하여 관리하지 않으므로 만료 여부만 체크.
		// 아래 appTokenValidationCheck 로직...
	} else {
		err = fmt.Errorf("token type error")
	}

	if err != nil {
		// DB error, 조회 X
		return false, err
	}

	// 토큰 정보 불일치
	if !flag {
		return false, consts.ErrInvalidClaims
	}

	// 토큰 만료 점검
	err = appTokenValidationCheck(entity.Token)

	if err != nil {
		// 만료 에러 확인
		if errors.Is(err, consts.ErrTokenExpired) {
			log.Println("토큰이 만료되었습니다.")
			return false, consts.ErrTokenExpired
		}
		// 서명 오류
		if errors.Is(err, consts.ErrTokenSignatureInvalid) {
			log.Println("서명이 유효하지 않습니다.")
			return false, consts.ErrTokenSignatureInvalid
		}
		// 그 외 오류
		log.Println("토큰 파싱 오류:", err)
		return false, consts.ErrTokenParsing
	}

	return true, nil
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
		log.Println("토큰이 유효합니다.")
		log.Println("UUID:", claims.Uuid)
		log.Println("만료 시간:", claims.ExpiresAt.Time)
		return nil
	} else {
		log.Println("토큰이 유효하지 않습니다.")
		return errors.New(consts.ErrInvalidClaims.Error())
	}
}

func (r *tokenUsecase) GenerateAppToken(body token.GenerateAppTokenRequestBody) (*token.GenerateAppTokenResponseDTO, error) {

	// 토큰 발급
	appToken, err := util.GenerateDeviceTokenJWT(r.jwtCfg.AppTokenExp, body.Uuid)

	fmt.Printf("요청한 uuid : %s, 발급된 토큰 : %s \n", body.Uuid, appToken)

	if err != nil {
		return nil, consts.ErrServerError
	}

	refreshToken, err := util.GenerateDeviceTokenJWT(r.jwtCfg.AppRefreshTokenExp, body.Uuid)
	if err != nil {
		return nil, consts.ErrServerError
	}

	// entity 생성
	tokenEntity := &shared.AppTokenEntity{
		Uuid:         body.Uuid,
		AppToken:     appToken,
		RefreshToken: refreshToken,
	}

	// DB 저장 실패시
	result, err := r.repo.PutIssuedAppToken(tokenEntity)

	if !result || err != nil {
		return nil, err
	}

	return &token.GenerateAppTokenResponseDTO{
		Body: token.GenerateAppTokenResponseBody{
			Uuid:         tokenEntity.Uuid,
			AppToken:     tokenEntity.AppToken,
			RefreshToken: tokenEntity.RefreshToken,
		},
	}, nil
}
