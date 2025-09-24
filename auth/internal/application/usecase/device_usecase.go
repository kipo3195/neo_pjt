package usecase

import (
	"auth/internal/application/usecase/input"
	"auth/internal/application/usecase/output"
	"auth/internal/claims"
	"auth/internal/consts"
	"auth/internal/domain/device/entity"
	"auth/internal/domain/device/repository"
	"auth/internal/infrastructure/storage"
	"context"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type deviceUsecase struct {
	repo    repository.DeviceRepository
	storage storage.DeviceStorage
}

type DeviceUsecase interface {
	GetDeviceRegistState(ctx context.Context, input input.DeviceRegistCheckInput) (output.DeviceRegistStateOutput, error)
	DeviceRegist(ctx context.Context, input input.DeviceRegistInput) (output.DeviceRegistOutput, error)
}

func NewDeviceUsecase(repo repository.DeviceRepository, storage storage.DeviceStorage) DeviceUsecase {
	return &deviceUsecase{
		repo:    repo,
		storage: storage,
	}
}

func (r *deviceUsecase) GetDeviceRegistState(ctx context.Context, input input.DeviceRegistCheckInput) (output.DeviceRegistStateOutput, error) {

	entity := entity.MakeDeviceRegistStateEntity(input.Id, input.Uuid)
	// 등록된 device 인지 체크

	_, err := r.repo.CheckDeviceRegist(ctx, entity)
	if err != nil {

		// 등록되지 않았을때
		if err == consts.ErrDeviceNotRegist {
			challenge := generateChallenge(entity.Id, entity.Uuid)
			r.storage.PutDeviceChallenge(entity.Id, entity.Uuid, challenge)
			return output.DeviceRegistStateOutput{
				DeviceRegistChallenge: challenge,
			}, nil
		} else {
			return output.DeviceRegistStateOutput{}, err
		}

	}

	at := r.storage.GetAccessToken(entity.Id, entity.Uuid)
	rt := r.storage.GetRefreshToken(entity.Id, entity.Uuid)

	// 없으면 신규 생성함.
	if at == "" || rt == "" {
		log.Printf("[GetDeviceRegistState] id : %s at, rt 신규 발급.", entity.Id)
		at, err = generateJWT(entity.Id, entity.Uuid, 60, []byte("access"), true)
		if err != nil {
			return output.DeviceRegistStateOutput{}, err
		}
		r.storage.PutAccessToken(entity.Id, entity.Uuid, at)
		rt, err := generateJWT(entity.Id, entity.Uuid, 30, []byte("refresh"), false)
		if err != nil {
			return output.DeviceRegistStateOutput{}, err
		}
		r.storage.PutRefreshToken(entity.Id, entity.Uuid, rt)
	}

	log.Printf("[GetDeviceRegistState] id : %s at : %s rt : %s", entity.Id, at, rt)

	return output.DeviceRegistStateOutput{
		RefreshToken: rt,
		AccessToken:  at,
	}, nil
}

func generateJWT(id string, uuid string, exp int, jwtKey []byte, accessFlag bool) (string, error) {
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
			return "", err
		}
		return accessToken, nil
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
			return "", err
		}

		return refreshToken, nil
	}

}

func (r *deviceUsecase) DeviceRegist(ctx context.Context, input input.DeviceRegistInput) (output.DeviceRegistOutput, error) {
	entity := entity.MakeDeviceRegistEntity(input.Id, input.Uuid, input.ModelName, input.Version, input.Challenge)

	// challenge 체크
	svChallenge := r.storage.GetDeviceChallenge(entity.Id, entity.Uuid)

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
	at, err := generateJWT(entity.Id, entity.Uuid, 60, []byte("access"), true)
	if err != nil {
		return output.DeviceRegistOutput{}, err
	}
	r.storage.PutAccessToken(entity.Id, entity.Uuid, at)

	// rt 생성 및 저장
	rt, err := generateJWT(entity.Id, entity.Uuid, 30, []byte("refresh"), false)
	if err != nil {
		return output.DeviceRegistOutput{}, err
	}
	r.storage.PutRefreshToken(entity.Id, entity.Uuid, rt)

	output := output.MakeDeviceRegistOutput(at, rt)

	// challenge 삭제
	r.storage.DeleteDeviceChallenge(entity.Id, entity.Uuid)

	return output, nil
}
