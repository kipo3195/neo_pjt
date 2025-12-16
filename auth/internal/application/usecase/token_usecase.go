package usecase

import (
	"auth/internal/application/usecase/input"
	"auth/internal/application/usecase/output"
	"auth/internal/application/util"
	"auth/internal/consts"
	"auth/internal/delivery/dto/token"
	"auth/internal/delivery/middleware/claims"
	"auth/internal/domain/shared"
	"auth/internal/domain/token/entity"
	"auth/internal/domain/token/repository"
	"auth/internal/infrastructure/config"
	"auth/internal/infrastructure/storage"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type tokenUsecase struct {
	repo               repository.TokenRepository
	storage            storage.AuthTokenStorage
	serviceUserStorage storage.ServiceUserStorage
	jwtCfg             *config.JWTConfig
	tokenCfg           config.TokenHashConfig
}

type TokenUsecase interface {
	AppTokenValidation(in input.AppTokenValidationInput, ctx context.Context) (bool, error)
	GenerateAppToken(ctx context.Context, in input.GenerateAppTokenInput) (*token.GenerateAppTokenResponseDTO, error)
	GenerateAuthToken(ctx context.Context, in input.GenerateAuthTokenInput) (output.GenerateAuthTokenOutput, error)
	CheckRefreshTokenWithExp(in input.RefreshTokenCheckInput, ctx context.Context) (bool, string, error)
	CheckRefreshToken(in input.RefreshTokenCheckInput, ctx context.Context) (string, error)
	ReIssueAccessToken(in input.ReIssueAccessTokenInput, ctx context.Context) (string, error)
	ReIssueAccessTokenSaved(ctx context.Context, in input.ReIssueAccessTokenSavedInput) error
}

func NewTokenUsecase(repo repository.TokenRepository, jwtCfg *config.JWTConfig, tokenCfg config.TokenHashConfig, storage storage.AuthTokenStorage, serviceUserStorage storage.ServiceUserStorage) TokenUsecase {
	return &tokenUsecase{
		repo:               repo,
		jwtCfg:             jwtCfg,
		tokenCfg:           tokenCfg,
		storage:            storage,
		serviceUserStorage: serviceUserStorage,
	}
}

func (r *tokenUsecase) AppTokenValidation(in input.AppTokenValidationInput, ctx context.Context) (bool, error) {

	// 결과 데이터를 미리 정의함.
	var flag bool
	var err error
	var token string
	entity := entity.NewAppTokenValidationEntity(in.Uuid, in.AppToken, in.Token)
	if in.TokenType == "appToken" {
		flag, err = r.repo.GetValidationAppToken(entity)
		token = in.AppToken
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
	err = appTokenValidationCheck(token)

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

func (r *tokenUsecase) GenerateAppToken(ctx context.Context, in input.GenerateAppTokenInput) (*token.GenerateAppTokenResponseDTO, error) {

	// 토큰 발급
	appToken, err := util.GenerateDeviceTokenJWT(r.jwtCfg.AppTokenExp, in.Uuid)

	fmt.Printf("요청한 uuid : %s, 발급된 토큰 : %s \n", in.Uuid, appToken)

	if err != nil {
		return nil, consts.ErrServerError
	}

	refreshToken, err := util.GenerateDeviceTokenJWT(r.jwtCfg.AppRefreshTokenExp, in.Uuid)
	if err != nil {
		return nil, consts.ErrServerError
	}

	// entity 생성
	tokenEntity := &shared.AppTokenEntity{
		Uuid:         in.Uuid,
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

func (r *tokenUsecase) GenerateAuthToken(ctx context.Context, in input.GenerateAuthTokenInput) (output.GenerateAuthTokenOutput, error) {

	entity := entity.MakeGenerateAuthTokenEntity(in.Id, in.Uuid)

	at := r.storage.GetAccessToken(entity.Id, entity.Uuid)
	rt := r.storage.GetRefreshToken(entity.Id, entity.Uuid)
	rtExp := r.storage.GetRefreshTokenExp(entity.Id, entity.Uuid)

	var err error

	// 무조건 재발급인 경우를 우선 체크함.
	if in.Force || at == "" || rt == "" || rtExp == "" {

		log.Printf("[GenerateAuthToken] id : %s at, rt 신규 발급..", entity.Id)

		hash := r.serviceUserStorage.GetUserHash(entity.Id)

		// at 생성
		at, _, err = generateJWT(entity.Id, entity.Uuid, hash, r.storage.GetTokenExpInfo(consts.DEVICE_ACCESSS_TOKEN), []byte(r.tokenCfg.AccessTokenHash), true)
		if err != nil {
			return output.GenerateAuthTokenOutput{}, err
		}
		// rt 생성
		rt, rtExp, err = generateJWT(entity.Id, entity.Uuid, hash, r.storage.GetTokenExpInfo(consts.DEVICE_REFRESH_TOKEN), []byte(r.tokenCfg.RefreshTokenHash), false)
		if err != nil {
			return output.GenerateAuthTokenOutput{}, err
		}
		// DB 저장.
		err = r.repo.PutAuthToken(ctx, entity.Id, entity.Uuid, at, rt, rtExp)
		if err != nil {
			return output.GenerateAuthTokenOutput{}, err
		}
		// 메모리 저장. 추후 redis저장으로 전환
		r.storage.PutAccessToken(entity.Id, entity.Uuid, at)
		r.storage.PutRefreshToken(entity.Id, entity.Uuid, rt)
		r.storage.PutRefreshTokenExp(entity.Id, entity.Uuid, rtExp)

	}

	log.Printf("[GenerateAuthToken] id : %s \n at : %s \n rt : %s \n rtExp : %s", entity.Id, at, rt, rtExp)

	output := output.GenerateAuthTokenOutput{
		RefreshToken:    rt,
		AccessToken:     at,
		RefreshTokenExp: rtExp,
	}

	return output, nil
}

func (r *tokenUsecase) CheckRefreshTokenWithExp(in input.RefreshTokenCheckInput, ctx context.Context) (bool, string, error) {

	entity := entity.MakeRefreshTokenCheckEntity(in.UserId, in.Uuid, in.RefreshToken)

	if in.WithoutId {
		// uuid를 통해 userId를 가져옴.
		temp, err := r.repo.GetUserIdWithRtAndUuid(ctx, entity)
		if err != nil {
			return false, "", err
		}

		if temp == "" {
			return false, "", consts.ErrUserIdDoesNotExist
		}

		entity.UserId = temp
	}

	rtExpDate := r.storage.GetRefreshTokenExp(entity.UserId, entity.Uuid)

	log.Println("[CheckRefreshTokenWithExp] userId : ", entity.UserId, "expDate : ", rtExpDate)

	expTime, err := time.Parse(time.RFC3339, rtExpDate)
	if err != nil {
		fmt.Println("시간 파싱 오류:", err)
		return false, "", consts.ErrTimeDataParsingError
	}

	now := time.Now()

	if now.After(expTime) {
		// 시간 지남
		return false, entity.UserId, nil
	} else {
		return true, entity.UserId, nil
	}
}

func (r *tokenUsecase) ReIssueAccessToken(in input.ReIssueAccessTokenInput, ctx context.Context) (string, error) {

	entity := entity.MakeReIssueAccessTokenEntity(in.UserId, in.Uuid)

	hash := r.serviceUserStorage.GetUserHash(in.UserId)

	// at 생성
	at, _, err := generateJWT(entity.UserId, entity.Uuid, hash, r.storage.GetTokenExpInfo(consts.DEVICE_ACCESSS_TOKEN), []byte(r.tokenCfg.AccessTokenHash), true)
	if err != nil {
		return "", err
	}

	r.storage.PutAccessToken(entity.UserId, entity.Uuid, at)

	return at, nil
}

func (r *tokenUsecase) ReIssueAccessTokenSaved(ctx context.Context, in input.ReIssueAccessTokenSavedInput) error {

	entity := entity.MakeReIssueAccessTokenSavedEntity(in.UserId, in.Uuid, in.Rt, in.At)
	return r.repo.UpdateReIssueAccessTokenInfo(ctx, entity)
}

func (r *tokenUsecase) CheckRefreshToken(in input.RefreshTokenCheckInput, ctx context.Context) (string, error) {

	entity := entity.MakeRefreshTokenCheckEntity(in.UserId, in.Uuid, in.RefreshToken)
	id, err := r.repo.GetUserIdWithRtAndUuid(ctx, entity)

	if err != nil {
		return "", err
	}

	if id == "" {
		log.Println("[CheckRefreshToken] uuid, rt로 사용자 id 조회 실패")
		return "", consts.ErrUserIdDoesNotExist
	}

	latestRt := r.storage.GetRefreshToken(id, entity.Uuid)

	if entity.RefreshToken == latestRt {
		return id, nil
	} else {
		return "", consts.ErrRefreshTokenAuthError
	}

}
