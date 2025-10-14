package usecase

import (
	"auth/internal/application/usecase/input"
	"auth/internal/application/usecase/output"
	"auth/internal/consts"
	"auth/internal/delivery/middleware/claims"
	"auth/internal/domain/device/entity"
	"auth/internal/domain/device/repository"
	"auth/internal/infrastructure/storage"
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type deviceUsecase struct {
	repo             repository.DeviceRepository
	deviceStorage    storage.DeviceStorage
	authTokenStorage storage.AuthTokenStorage
	accessHash       string
	refreshHash      string
}

type DeviceUsecase interface {
	GetDeviceRegistState(ctx context.Context, input input.DeviceRegistCheckInput) (string, error)
	DeviceRegist(ctx context.Context, input input.DeviceRegistInput) (output.DeviceRegistOutput, error)
}

func NewDeviceUsecase(repo repository.DeviceRepository, deviceStorage storage.DeviceStorage, authTokenStorage storage.AuthTokenStorage, accessHash string, refreshHash string) DeviceUsecase {
	return &deviceUsecase{
		repo:             repo,
		deviceStorage:    deviceStorage,
		authTokenStorage: authTokenStorage,
		accessHash:       accessHash,
		refreshHash:      refreshHash,
	}
}

func (r *deviceUsecase) GetDeviceRegistState(ctx context.Context, input input.DeviceRegistCheckInput) (string, error) {

	entity := entity.MakeDeviceRegistStateEntity(input.Id, input.Uuid)

	// 등록된 device 인지 체크
	_, err := r.repo.CheckDeviceRegist(ctx, entity)
	if err != nil {

		// 등록되지 않았을때 challenge 발급
		if err == consts.ErrDeviceNotRegist {
			challenge := generateChallenge(entity.Id, entity.Uuid)
			r.deviceStorage.PutDeviceChallenge(entity.Id, entity.Uuid, challenge)
			return challenge, err
		} else {
			return "", err
		}
	}
	return "", nil
}

func generateJWT(id string, uuid string, exp int, jwtKey []byte, accessFlag bool) (string, string, error) {
	// 현재 기준 시간
	now := time.Now()
	// issuer
	const issuer = "device"

	if accessFlag {
		// Access 토큰 유효기간 설정
		expTime := now.Add(time.Duration(exp) * time.Minute)

		accessClaims := &claims.DeviceJWTClaims{
			Id:   id,
			Uuid: uuid,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expTime),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				Issuer:    issuer,
			},
		}
		// 토큰 생성 (HS256 사용)
		accToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
		// 서명 및 문자열 반환
		accessToken, err := accToken.SignedString(jwtKey)
		if err != nil {
			return "", "", err
		}
		return accessToken, "", nil
	} else {
		// Refresh 토큰 유효기간 설정
		expTime := now.Add(time.Duration(exp) * 24 * time.Hour)

		refreshClaims := &claims.DeviceJWTClaims{
			Id:   id,
			Uuid: uuid,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expTime),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				Issuer:    issuer,
			},
		}
		reToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
		refreshToken, err := reToken.SignedString(jwtKey)
		if err != nil {
			return "", "", err
		}

		return refreshToken, expTime.Format(time.RFC3339), nil
	}
}

func (r *deviceUsecase) DeviceRegist(ctx context.Context, input input.DeviceRegistInput) (output.DeviceRegistOutput, error) {
	entity := entity.MakeDeviceRegistEntity(input.Id, input.Uuid, input.ModelName, input.Version, input.Challenge)

	// challenge 체크
	svChallenge := r.deviceStorage.GetDeviceChallenge(entity.Id, entity.Uuid)

	if svChallenge == "" {
		return output.DeviceRegistOutput{}, consts.ErrDeviceChallengeExpired
	}

	if svChallenge != entity.Challenge {
		return output.DeviceRegistOutput{}, consts.ErrDeviceChallengeMismatch
	}

	err := r.repo.PutDevice(ctx, entity)

	if err != nil {
		return output.DeviceRegistOutput{}, err
	}

	// at 생성 및 저장
	at, _, err := generateJWT(entity.Id, entity.Uuid, r.deviceStorage.GetDeviceTokenExp(consts.DEVICE_ACCESSS_TOKEN), []byte("access"), true)
	if err != nil {
		return output.DeviceRegistOutput{}, err
	}
	r.authTokenStorage.PutAccessToken(entity.Id, entity.Uuid, at)

	// rt 생성 및 저장
	rt, rtExp, err := generateJWT(entity.Id, entity.Uuid, r.deviceStorage.GetDeviceTokenExp(consts.DEVICE_REFRESH_TOKEN), []byte("refresh"), false)
	if err != nil {
		return output.DeviceRegistOutput{}, err
	}

	r.authTokenStorage.PutRefreshToken(entity.Id, entity.Uuid, rt)
	r.authTokenStorage.PutRefreshTokenExp(entity.Id, entity.Uuid, rtExp)

	r.repo.PutAuthToken(ctx, entity.Id, entity.Uuid, at, rt, rtExp)

	output := output.MakeDeviceRegistOutput(at, rt, rtExp)

	// challenge 삭제
	r.deviceStorage.DeleteDeviceChallenge(entity.Id, entity.Uuid)

	return output, nil
}
