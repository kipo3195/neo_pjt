package usecases

import (
	"auth/internal/claims"
	"auth/internal/consts"
	clReqDto "auth/internal/domains/certification/dto/client/request"
	clResDto "auth/internal/domains/certification/dto/client/response"
	"auth/internal/domains/certification/entities"
	clRepo "auth/internal/domains/certification/repositories/client"
	"auth/internal/utils"
	"auth/pkg/config"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type certificationUsecase struct {
	repo      clRepo.CerificationRepository
	jwtCfg    *config.JWTConfig
	authUtile *utils.AuthUtil
}

type CertificationUsecase interface {
	GetAuth(requestDTO clReqDto.AuthRequestDTO) (*clResDto.AuthResponseDTO, error)
}

func NewCertificationUsecase(repo clRepo.CerificationRepository, jwtCfg *config.JWTConfig, authUtile *utils.AuthUtil) CertificationUsecase {
	return &certificationUsecase{
		repo:      repo,
		jwtCfg:    jwtCfg,
		authUtile: authUtile}
}

func (u *certificationUsecase) GetAuth(requestDTO clReqDto.AuthRequestDTO) (*clResDto.AuthResponseDTO, error) {

	// app hash 부터 검증
	flag, err := u.repo.GetValidation(toAppTokenValidationEntity(requestDTO.Header.Uuid, requestDTO.Header.Token))
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 등록된 UUID가 아님
			return nil, consts.ErrUnregisteredUuid
		} else {
			// DB error
			return nil, consts.ErrDB
		}
	}
	// 토큰 정보 불일치, 재발급 필요.
	if !flag {
		return nil, consts.ErrTokenMismatch
	}

	// 사용자 정보 검증
	auth, err := u.repo.CheckAuth(toGetAuthEntity(requestDTO.Body))

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 등록된 사용자가 없음.
			return nil, consts.ErrAuthenticationFailed
		} else {
			return nil, consts.ErrDB
		}
	}

	// 등록 사용자 검증 (user_hash 구하기)
	userHash, err := u.repo.GetUserHash(toGetAuthEntity(requestDTO.Body))

	// 등록된 사용자가 아님.
	if userHash == "" || err != nil {
		// 단, repo.GetAuth에서 Scan으로 조회하고 있으므로 이건 ErrRecordNotFound하지 않지만.. err이 발생할 수 있으니 추가함.
		// service_users에 등록된 사용자가 아님.
		return nil, consts.ErrUnregisteredUser
	}

	auth.Userhash = userHash

	var accessToken, refreshToken string

	acc, re, err := GenerateAuthJWT(auth.Userhash, u.jwtCfg.AccessExp, u.jwtCfg.RefressExp, []byte(u.jwtCfg.Key))
	if err != nil {
		println("JWT TOKEN MAKE ERROR ! :", err)
		return nil, consts.ErrServerError
	} else {
		accessToken = acc
		refreshToken = re
	}
	// config 파일을 풀 수 있는 대칭키
	// configKey = getConfigkey()

	return &clResDto.AuthResponseDTO{
		Body: clResDto.AuthResponseBody{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

func toAppTokenValidationEntity(uuid string, appToken string) entities.AppTokenValidationEntity {
	return entities.AppTokenValidationEntity{
		Uuid:     uuid,
		AppToken: appToken,
	}
}

func toGetAuthEntity(body clReqDto.AuthRequestBody) entities.AuthInfoEntity {
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
