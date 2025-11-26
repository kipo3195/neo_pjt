package usecase

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
	"message/internal/application/usecase/input"
	"message/internal/application/usecase/output"
	"message/internal/consts"
	"message/internal/delivery/adapter"
	"message/internal/domain/otp/entity"
	"message/internal/domain/otp/repository"
	"message/internal/infrastructure/storage"
	"strings"
	"time"
)

type otpUsecase struct {
	repository repository.OtpRepository
	svChKey    string
	svNoKey    string
	otpStorage storage.OtpStorage
}

type OtpUsecase interface {
	OtpKeyRegist(ctx context.Context, input input.OtpKeyRegistInput) (output output.OtpKeyregistOutput, err error)
}

func NewOtpUsecase(repo repository.OtpRepository, storage storage.OtpStorage, svChKey string, svNoKey string) OtpUsecase {
	return &otpUsecase{
		repository: repo,
		otpStorage: storage,
		svChKey:    svChKey,
		svNoKey:    svNoKey,
	}
}

func (u *otpUsecase) OtpKeyRegist(ctx context.Context, input input.OtpKeyRegistInput) (out output.OtpKeyregistOutput, err error) {

	entity := entity.MakeOtpKeyRegistEntity(input.Id, input.Uuid, input.ChKey, input.NoKey)

	// 키 생성 + 시간 생성,
	regDate, chatOtpKey, noteOtpKey, err := makeOtpKey(entity.ChKey, entity.NoKey, u.svChKey, u.svNoKey)

	if err != nil {
		return output.OtpKeyregistOutput{}, err
	}

	entity.ChatOtpKey = chatOtpKey
	entity.NoteOtpKey = noteOtpKey
	entity.RegDate = regDate

	log.Printf("[OtpKeyRegist] id:%s, regDate:%s, chatOtpKey:%s, noteOtpKey:%s\n", entity.Id, regDate, chatOtpKey, noteOtpKey)

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

	output := adapter.MakeOtpKeyRegistOutput(regDate)
	return output, nil
}

func makeOtpKey(chKey string, noKey string, svChKey string, svNoKey string) (regDate string, chatOtpKey string, noteOtpKey string, err error) {

	// 현재 여기 로직 점검중 20251127
	//log.Printf("[makeOtpKey] chKey:%s, noKey:%s, svChKey:%s, svNoKey:%s\n", chKey, noKey, svChKey, svNoKey)

	fmt.Printf("chKey >>>%s<<<\n", chKey)
	fmt.Printf("noKey >>>%s<<<\n", noKey)
	// --- PEM 공개키 파싱 함수 ---
	parsePubKey := func(pubKeyPEM string) (*rsa.PublicKey, error) {
		block, _ := pem.Decode([]byte(pubKeyPEM))
		if block == nil {
			return nil, consts.ErrFailedToDecodePEMBlock
		}
		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rsaPub, ok := pub.(*rsa.PublicKey)
		if !ok {
			return nil, consts.ErrFailedToParsePublicKey
		}
		return rsaPub, nil
	}

	// PEM 내부 개행 문자(\n) 문제
	// JSON에서 \n은 문자열 안의 실제 개행으로 변환됩니다.
	// 하지만 만약 Go 코드에서 body를 구조체로 unmarshal 한 후, \n이 그대로 문자열 \ + n으로 들어가는 경우가 있습니다.
	chKeyClean := strings.ReplaceAll(chKey, "\r", "")
	chKeyClean = strings.TrimSpace(chKeyClean)

	noKeyClean := strings.ReplaceAll(noKey, "\r", "")
	noKeyClean = strings.TrimSpace(noKeyClean)

	// --- 공개키 로드 ---
	chPub, err := parsePubKey(chKeyClean)
	if err != nil {
		return "", "", "", err
	}

	noPub, err := parsePubKey(noKeyClean)
	if err != nil {
		return "", "", "", err
	}

	// --- 암호화 ---
	chEncrypted, err := rsa.EncryptPKCS1v15(rand.Reader, chPub, []byte(svChKey))
	if err != nil {
		return "", "", "", err
	}

	noEncrypted, err := rsa.EncryptPKCS1v15(rand.Reader, noPub, []byte(svNoKey))
	if err != nil {
		return "", "", "", err
	}

	// --- Base64 인코딩 (문자열로 반환 용이하게) ---
	chEncB64 := base64.StdEncoding.EncodeToString(chEncrypted)
	noEncB64 := base64.StdEncoding.EncodeToString(noEncrypted)

	// --- 발급시간 UTC ---
	nowUTC := time.Now().UTC().Format(time.RFC3339)

	fmt.Printf("[makeOtpKey] nowUTC:%s, chEncrypted(Base64): %s, noEncrypted(Base64): %s\n", nowUTC, chEncB64, noEncB64)

	return nowUTC, chEncB64, noEncB64, nil
}
