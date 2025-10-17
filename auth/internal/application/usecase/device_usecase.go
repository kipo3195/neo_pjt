package usecase

import (
	"auth/internal/application/usecase/input"
	"auth/internal/consts"
	"auth/internal/delivery/middleware/claims"
	"auth/internal/domain/device/entity"
	"auth/internal/domain/device/repository"
	"auth/internal/infrastructure/storage"
	"context"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type deviceUsecase struct {
	repo          repository.DeviceRepository
	deviceStorage storage.DeviceStorage
	accessHash    string
	refreshHash   string
}

type DeviceUsecase interface {
	GetDeviceRegistState(ctx context.Context, input input.DeviceRegistStateInput) (string, error)
	DeviceRegistCheck(ctx context.Context, input input.DeviceRegistInput) (bool, error)
	RemoveDeviceChallenge(ctx context.Context, input input.RemoveDeviceChallengeInput)
	UpdateDeviceInfo(ctx context.Context, input input.UpdateDeviceInfoInput) error
}

func NewDeviceUsecase(repo repository.DeviceRepository, deviceStorage storage.DeviceStorage, accessHash string, refreshHash string) DeviceUsecase {
	return &deviceUsecase{
		repo:          repo,
		deviceStorage: deviceStorage,
		accessHash:    accessHash,
		refreshHash:   refreshHash,
	}
}

func (r *deviceUsecase) GetDeviceRegistState(ctx context.Context, input input.DeviceRegistStateInput) (string, error) {

	// 디바이스 등록 체크용 entity 생성함.
	entity := entity.MakeDeviceRegistStateEntity(input.Id, input.Uuid)

	// 등록된 device 인지 체크 - 신규 인증이라면 값이 없이 ErrDeviceNotRegist, 별도 DB 에러라면 그대로 리턴
	refresh, err := r.repo.CheckDeviceRegist(ctx, entity)
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

	if refresh {
		// 재등록 API에서만 호출 되어야함.
		log.Println("[GetDeviceRegistState] device refresh check OK")
		// 만료 되었을 수 있으므로 재발급함.
		challenge := generateChallenge(entity.Id, entity.Uuid)
		r.deviceStorage.PutDeviceChallenge(entity.Id, entity.Uuid, challenge)
		return challenge, err
	} else {
		return "", consts.ErrServerError
	}
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

func (r *deviceUsecase) DeviceRegistCheck(ctx context.Context, input input.DeviceRegistInput) (bool, error) {
	entity := entity.MakeDeviceRegistEntity(input.Id, input.Uuid, input.ModelName, input.Version, input.Challenge)

	// challenge 체크
	svChallenge := r.deviceStorage.GetDeviceChallenge(entity.Id, entity.Uuid)
	log.Println("svChallenge : ", svChallenge)
	if svChallenge == "" {
		return false, consts.ErrDeviceChallengeExpired
	}

	log.Println("clientChallenge : ", entity.Challenge)
	if svChallenge != entity.Challenge {
		return false, consts.ErrDeviceChallengeMismatch
	}

	err := r.repo.PutDevice(ctx, entity)
	log.Println("err : ", err)

	if err != nil {
		return false, err
	} else {
		return true, nil
	}

}

func (r *deviceUsecase) RemoveDeviceChallenge(ctx context.Context, input input.RemoveDeviceChallengeInput) {
	r.deviceStorage.DeleteDeviceChallenge(input.Id, input.Uuid)
}

func (r *deviceUsecase) UpdateDeviceInfo(ctx context.Context, input input.UpdateDeviceInfoInput) error {

	entity := entity.MakeDeviceRegistEntity(input.Id, input.Uuid, input.ModelName, input.Version, "")
	return r.repo.UpdateDeviceInfo(ctx, entity)

}
