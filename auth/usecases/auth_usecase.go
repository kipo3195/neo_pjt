package usecases

import (
	"auth/claims"
	"auth/config"
	consts "auth/consts"
	clDto "auth/dto/client"
	dto "auth/dto/common"
	commonDto "auth/dto/server/common"
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
	GetAuth(*clDto.LoginRequestHeader, *clDto.AuthRequest) (*entities.AuthEntity, *dto.ErrorResponse, bool)
	GenerateAppToken(body commonDto.GenerateAppTokenRequest) (*entities.AppTokenEntity, error)
}

func NewAuthUsecase(repo repositories.AuthRepository, jwtCfg *config.JWTConfig) AuthUsecase {
	return &authUsecase{repo: repo, jwtCfg: jwtCfg}
}

func (u *authUsecase) GetAuth(header *clDto.LoginRequestHeader, body *clDto.AuthRequest) (*entities.AuthEntity, *dto.ErrorResponse, bool) {

	// app hash 부터 검증
	flag, err := u.repo.GetValidation(toAppTokenValidationEntity(header.Uuid, header.Token))
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			// 매핑된 hash 정보가 없음
			return nil, &dto.ErrorResponse{
				Code:    consts.AUTH_F001,
				Message: consts.AUTH_F001_MSG,
			}, true
		default:
			// 기타 DB 에러
			return nil, &dto.ErrorResponse{
				Code:    consts.E_102,
				Message: consts.E_102_MSG,
			}, false
		}
	}
	// 토큰 정보 불일치, 재발급 필요.
	if !flag {
		return nil, &dto.ErrorResponse{
			Code:    consts.AUTH_F002,
			Message: consts.AUTH_F002_MSG,
		}, true
	}

	// 사용자 정보 검증
	auth, err := u.repo.CheckAuth(toGetAuthEntity(body))

	// ID, PW 일치하는 사용자가 없음.
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			// ID, PW와 일치하는 사용자가 없음
			return nil, &dto.ErrorResponse{
				Code:    consts.AUTH_F003,
				Message: consts.AUTH_F003_MSG,
			}, true
		default:
			// 기타 DB 에러
			return nil, &dto.ErrorResponse{
				Code:    consts.E_102,
				Message: consts.E_102_MSG,
			}, false
		}
	}

	// 등록 사용자 검증 (user_hash 구하기)
	auth.Userhash, err = u.repo.GetUserHash(toGetAuthEntity(body))

	// 등록된 사용자가 아님.
	if auth.Userhash == "" || err != nil {
		// 단, repo.GetAuth에서 Scan으로 조회하고 있으므로 이건 ErrRecordNotFound하지 않지만.. err이 발생할 수 있으니 추가함.
		// service_users에 등록된 사용자가 아님.
		return nil, &dto.ErrorResponse{
			Code:    consts.AUTH_F004,
			Message: consts.AUTH_F004_MSG,
		}, true
	}

	var accessToken, refreshToken string

	acc, re, err := GenerateAuthJWT(auth.Userhash, u.jwtCfg.AccessExp, u.jwtCfg.RefressExp, []byte(u.jwtCfg.Key))
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
	// configKey = getConfigkey()

	return &entities.AuthEntity{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil, false
}

// func toAppTokenValidationEntity(uuid string, appToken string) entities.AppTokenValidationEntity {
// 	return entities.AppTokenValidationEntity{
// 		Uuid:     uuid,
// 		AppToken: appToken,
// 	}
// }

func toGetAuthEntity(body *clDto.AuthRequest) entities.AuthInfoEntity {
	return entities.AuthInfoEntity{
		Id:       body.Id,
		Password: body.Password,
	}
}

// 사인을 위한 key는 byte[]여야 함
func GenerateAuthJWT(userHash string, accessExp int, refreshExp int, jwtKey []byte) (string, string, error) {
	// 현재 기준 시간
	now := time.Now()
	// issuer
	const issuer = "neo"

	// Access 토큰 유효기간 설정
	accExpTime := now.Add(time.Duration(accessExp) * time.Minute)

	// 사용자 정보 포함한 Claims 생성
	accessClaims := &claims.AuthJWTClaims{
		UserHash: userHash,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accExpTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    issuer,
		},
	}

	// Refresh 토큰 유효기간 설정
	reExpTime := now.Add(time.Duration(refreshExp) * 24 * time.Hour)

	refreshClaims := &claims.AuthJWTClaims{
		UserHash: userHash,
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

func (r *authUsecase) GenerateAppToken(body commonDto.GenerateAppTokenRequest) (*entities.AppTokenEntity, error) {

	// 토큰 발급
	appToken, err := generateDeviceTokenJWT(r.jwtCfg.AppTokenExp, body.Uuid)

	fmt.Printf("요청한 uuid : %s, 발급된 토큰 : %s \n", body.Uuid, appToken)

	if err != nil {
		return nil, err
	}

	refreshToken, err := generateDeviceTokenJWT(r.jwtCfg.AppRefreshTokenExp, body.Uuid)
	if err != nil {
		return nil, err
	}

	// entity 생성
	tokenEntity := &entities.AppTokenEntity{
		Uuid:         body.Uuid,
		AppToken:     appToken,
		RefreshToken: refreshToken,
	}

	// DB 저장 실패시
	result, err := r.repo.PutIssuedAppToken(tokenEntity)

	if !result || err != nil {
		return nil, err
	}

	return tokenEntity, nil
}

func generateDeviceTokenJWT(appTokenExp int, uuid string) (string, error) { //
	now := time.Now()
	// issuer
	const issuer = "auth"

	// Access 토큰 유효기간 설정
	accExpTime := now.Add(time.Duration(appTokenExp) * 24 * time.Hour)

	fmt.Println("jwt 토큰 생성  1 :", accExpTime)

	// 사용자 정보 포함한 Claims 생성
	uuidClaim := &claims.DeviceJWTClaims{
		Uuid: uuid,
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
