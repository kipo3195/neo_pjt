package usecases

import (
	"auth/config"
	consts "auth/consts"
	"auth/dto"
	"auth/entities"
	"auth/repositories"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type authUsecase struct {
	repo   repositories.AuthRepository
	jwtCfg *config.JWTConfig
}

type AuthUsecase interface {
	GetAuth(*dto.LoginRequestHeader, *dto.AuthRequest) (*entities.Auth, *dto.ErrorResponse, bool)
	GenerateDeviceToken(body dto.GenerateDeviceTokenRequest) (*entities.DeviceToken, error)
	GenerateJWT(uuid string) (string, error)
}

func NewAuthUsecase(repo repositories.AuthRepository, jwtCfg *config.JWTConfig) AuthUsecase {
	return &authUsecase{repo: repo, jwtCfg: jwtCfg}
}

func (u *authUsecase) GetAuth(header *dto.LoginRequestHeader, body *dto.AuthRequest) (*entities.Auth, *dto.ErrorResponse, bool) {

	// app hash 부터 검증
	flag, err := u.repo.GetValidation(header)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			// 매핑된 hash 정보가 없음
			return nil, &dto.ErrorResponse{
				Code:    consts.F_103,
				Message: consts.F_103_MSG,
			}, true
		default:
			// 기타 DB 에러
			return nil, &dto.ErrorResponse{
				Code:    consts.E_102,
				Message: consts.E_102_MSG,
			}, false
		}
	}
	// 토큰 정보 불일치
	if !flag {
		return nil, &dto.ErrorResponse{
			Code:    consts.F_104,
			Message: consts.F_104_MSG,
		}, true
	}

	// 사용자 정보 검증
	auth, err := u.repo.GetAuth(body)

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			// 매핑된 hash 정보가 없음
			return nil, &dto.ErrorResponse{
				Code:    consts.F_105,
				Message: consts.F_105_MSG,
			}, true
		default:
			// 기타 DB 에러
			return nil, &dto.ErrorResponse{
				Code:    consts.E_102,
				Message: consts.E_102_MSG,
			}, false
		}
	}

	var result, accessToken, refreshToken, configKey string

	result = "success"

	acc, re, err := GenerateJWT(auth.Id, u.jwtCfg.AccessExp, u.jwtCfg.RefressExp, []byte(u.jwtCfg.Key))
	if err != nil {
		println("JWT TOKEN MAKE ERROR ! :", err)
		return nil, &dto.ErrorResponse{
			Code:    consts.E_500,
			Message: consts.E_500_MSG,
		}, false
	} else {
		accessToken = acc
		refreshToken = re
	}
	// config 파일을 풀 수 있는 대칭키
	configKey = getConfigkey()

	return &entities.Auth{Result: result, AccessToken: accessToken, RefreshToken: refreshToken, ConfigKey: configKey}, nil, false
}

func getConfigkey() string {
	return ""
}

type JWTClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// 사인을 위한 key는 byte[]여야 함
func GenerateJWT(username string, accessExp int, refreshExp int, jwtKey []byte) (string, string, error) {
	// 현재 기준 시간
	now := time.Now()
	// issuer
	const issuer = "neo"

	// Access 토큰 유효기간 설정
	accExpTime := now.Add(time.Duration(accessExp) * time.Minute)

	// 사용자 정보 포함한 Claims 생성
	accessClaims := JWTClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accExpTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    issuer,
		},
	}

	// Refresh 토큰 유효기간 설정
	reExpTime := now.Add(time.Duration(refreshExp) * 24 * time.Hour)

	refreshClaims := JWTClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(reExpTime),
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

	reToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err := reToken.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (r *authUsecase) GenerateDeviceToken(body dto.GenerateDeviceTokenRequest) (*entities.DeviceToken, error) {

	// 토큰 발급
	token, err := r.GenerateJWT(body.Uuid)

	fmt.Printf("요청한 uuid : %s, 발급된 토큰 : %s \n", body.Uuid, token)

	if err != nil {
		return &entities.DeviceToken{}, err
	}

	// entity 생성
	tokenEntity := &entities.DeviceToken{
		Uuid:  body.Uuid,
		Token: token,
	}

	// DB 저장 실패시
	result, err := r.repo.PutDeviceToken(tokenEntity)

	if !result || err != nil {
		return &entities.DeviceToken{}, err
	}

	return tokenEntity, nil
}

func (r *authUsecase) GenerateJWT(uuid string) (string, error) {
	now := time.Now()
	// issuer
	const issuer = "auth"

	// Access 토큰 유효기간 설정
	accExpTime := now.Add(time.Duration(3) * time.Minute)

	fmt.Println("jwt 토큰 생성  1 :", accExpTime)

	// 사용자 정보 포함한 Claims 생성
	uuidClaim := JWTClaims{
		Username: uuid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accExpTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    issuer,
		},
	}

	// 토큰 생성 (HS256 사용)
	accToken := jwt.NewWithClaims(jwt.SigningMethodHS256, uuidClaim)

	secret := []byte("neo-test-secret-key")

	// 서명 및 문자열 반환
	accessToken, err := accToken.SignedString(secret)

	fmt.Println("jwt 토큰 생성  2 :", accessToken)
	return accessToken, err
}
