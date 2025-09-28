package usecase

import (
	"auth/internal/application/usecase/input"
	"auth/internal/delivery/dto/certification"
	"auth/internal/delivery/middleware/claims"
	"auth/internal/domain/certification/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type certificationUsecase struct {
	repository repository.CerificationRepository
}

type CertificationUsecase interface {
	GetAuth(in input.LoginInput) (*certification.AuthResponseDTO, error)
}

func NewCertificationUsecase(repository repository.CerificationRepository) CertificationUsecase {
	return &certificationUsecase{
		repository: repository,
	}
}

func (u *certificationUsecase) GetAuth(n input.LoginInput) (*certification.AuthResponseDTO, error) {

	return nil, nil
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
