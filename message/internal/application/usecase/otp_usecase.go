package usecase

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"log"
	"message/internal/application/usecase/input"
	"message/internal/application/usecase/output"
	"message/internal/consts"
	"message/internal/delivery/adapter"
	"message/internal/domain/otp/entity"
	"message/internal/domain/otp/repository"
	"message/internal/infrastructure/storage"
	"time"
)

type otpUsecase struct {
	repository   repository.OtpRepository
	svChKey      string
	svNoKey      string
	svKeyVersion string
	otpStorage   storage.OtpStorage
}

type OtpUsecase interface {
	OtpKeyRegist(ctx context.Context, input input.OtpKeyRegistInput) (output output.OtpKeyregistOutput, err error)
}

func NewOtpUsecase(repo repository.OtpRepository, storage storage.OtpStorage, svChKey string, svNoKey string, svKeyVersion string) OtpUsecase {
	return &otpUsecase{
		repository:   repo,
		otpStorage:   storage,
		svChKey:      svChKey,
		svNoKey:      svNoKey,
		svKeyVersion: svKeyVersion,
	}
}

func (u *otpUsecase) OtpKeyRegist(ctx context.Context, input input.OtpKeyRegistInput) (out output.OtpKeyregistOutput, err error) {

	entity := entity.MakeOtpKeyRegistEntity(input.Id, input.Uuid, input.ChKey, input.NoKey)

	// 키 생성 + 시간 생성,
	chatOtpKey, err := makeOtpKey(entity.ChKey, u.svChKey, consts.CHAT)
	if err != nil {
		return output.OtpKeyregistOutput{}, err
	}

	noteOtpKey, err := makeOtpKey(entity.NoKey, u.svNoKey, consts.NOTE)
	if err != nil {
		return output.OtpKeyregistOutput{}, err
	}

	// --- 발급시간 UTC ---
	nowUTC := time.Now().UTC().Format(time.RFC3339)

	// DB 저장 데이터
	entity.ChatOtpKey = chatOtpKey
	entity.NoteOtpKey = noteOtpKey
	entity.OtpRegDate = nowUTC
	entity.SvKeyVersion = u.svKeyVersion

	// DB 저장 (키 + 시간)
	err = u.repository.SaveOtpKey(ctx, entity)
	if err != nil {
		return output.OtpKeyregistOutput{}, err
	}

	// 메모리 저장 (storage)
	err = u.otpStorage.SaveOtpKeyStorage(ctx, entity)
	if err != nil {
		return output.OtpKeyregistOutput{}, err
	}

	log.Printf("[OtpKeyRegist] id:%s, regDate:%s, version:%s success.\n", entity.Id, entity.OtpRegDate, u.svKeyVersion)

	output := adapter.MakeOtpKeyRegistOutput(entity.OtpRegDate, entity.SvKeyVersion)
	return output, nil
}

func makeOtpKey(clientKey string, serverKey string, t string) (string, error) {

	log.Printf("[makeOtpKey] type:%s, clientKey:%s, serverKey:%s\n", t, clientKey, serverKey)

	// 1) Base64 decode
	decoded, err := base64.StdEncoding.DecodeString(clientKey)
	if err != nil {
		log.Println("[parsePubKeyFromBase64] Base64 Decode Error:", err)
		return "", err
	}

	// 2) decoded = PEM 전체 문서여야 한다
	block, _ := pem.Decode(decoded)
	if block == nil {
		log.Println("[parsePubKeyFromBase64] PEM Decode Error")
		return "", consts.ErrFailedToDecodePEMBlock
	}

	// 3) parse PKIX
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Println("[parsePubKeyFromBase64] ParsePKIXPublicKey Error:", err)
		return "", consts.ErrFailedToParsePublicKey
	}

	// 4) 타입 체크
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		log.Println("[parsePubKeyFromBase64] Public Key is not RSA Type")
		return "", consts.ErrFailedToParsePublicKey
	}

	// --- 암호화 ---
	// RSA PKCS1v15 방식으로 서비스키 암호화하여 OTP 키 생성
	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPub, []byte(serverKey))
	if err != nil {
		return "", consts.ErrFailedToEncryptOtpKey
	}

	// --- Base64 인코딩 (문자열로 반환 용이하게) ---
	EncB64 := base64.StdEncoding.EncodeToString(encrypted)

	return EncB64, nil
}
