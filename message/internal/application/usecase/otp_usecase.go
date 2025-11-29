package usecase

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"log"
	"message/internal/application/usecase/input"
	"message/internal/application/usecase/output"
	"message/internal/consts"
	"message/internal/delivery/adapter"
	"message/internal/domain/otp/entity"
	"message/internal/domain/otp/repository"
	"message/internal/infrastructure/storage"
	"time"

	"gorm.io/gorm"
)

type otpUsecase struct {
	repository     repository.OtpRepository
	svChKey        string
	svChKeyVersion string
	svNoKey        string
	svNoKeyVersion string
	otpStorage     storage.OtpStorage
}

type OtpUsecase interface {
	OtpKeyRegist(ctx context.Context, input input.OtpKeyRegistInput) (output output.OtpKeyregistOutput, err error)
	GetMyOtpInfo(ctx context.Context, input input.MyOtpInfoInput) (output []output.MyOtpInfoOutput, err error)
}

func NewOtpUsecase(repo repository.OtpRepository, storage storage.OtpStorage, svChKey string, svChKeyVersion string, svNoKey string, svNoKeyVersion string) OtpUsecase {

	// 서버 대칭키 종류에 따른 저장
	storage.PutServerKey(svChKey, "chat", svChKeyVersion)
	storage.PutServerKey(svNoKey, "note", svNoKeyVersion)

	return &otpUsecase{
		repository:     repo,
		otpStorage:     storage,
		svChKey:        svChKey,
		svChKeyVersion: svChKeyVersion,
		svNoKeyVersion: svNoKeyVersion,
		svNoKey:        svNoKey,
	}
}

func (u *otpUsecase) OtpKeyRegist(ctx context.Context, input input.OtpKeyRegistInput) (output.OtpKeyregistOutput, error) {

	otpKeyInfoEntity := make([]entity.OtpKeyInfoEntity, 0)

	for i := 0; i < len(input.DevicePubKey); i++ {
		e := entity.OtpKeyInfoEntity{
			Kind: input.DevicePubKey[i].Kind,
			Key:  input.DevicePubKey[i].Key,
		}
		otpKeyInfoEntity = append(otpKeyInfoEntity, e)
	}
	entity := entity.MakeOtpKeyRegistEntity(input.Id, input.Uuid, otpKeyInfoEntity)

	// 키 생성 + 시간 생성,
	// --- 발급시간 UTC ---
	nowUTC := time.Now().UTC().Format(time.RFC3339)
	// 키 생성
	for i := 0; i < len(entity.OtpKeyInfoEntity); i++ {

		svKey, svKeyVersion := u.otpStorage.GetServerKey(entity.OtpKeyInfoEntity[i].Kind)

		key, err := makeOtpKey(entity.OtpKeyInfoEntity[i].Key, svKey, entity.OtpKeyInfoEntity[i].Kind)
		if err != nil {
			log.Printf("[OtpKeyRegist] key :%s is invalid.\n", entity.OtpKeyInfoEntity[i].Key)
			continue
		}
		entity.OtpKeyInfoEntity[i].OtpKey = key
		entity.OtpKeyInfoEntity[i].OtpRegDate = nowUTC
		entity.OtpKeyInfoEntity[i].SvKeyVersion = svKeyVersion
	}
	// DB 저장 (키 + 시간)
	err := u.repository.SaveOtpKey(ctx, entity)
	if err != nil {
		return output.OtpKeyregistOutput{}, err
	}

	// 메모리 저장 (storage)
	id := entity.Id
	uuid := entity.Uuid
	for _, keys := range entity.OtpKeyInfoEntity {
		u.otpStorage.SaveOtpKeyStorage(ctx, id, uuid, keys)
	}

	log.Printf("[OtpKeyRegist] id:%s, uuid:%s, regDate:%s, success.\n", id, uuid, nowUTC)

	output := adapter.MakeOtpKeyRegistOutput(nowUTC, u.svChKeyVersion, u.svNoKeyVersion)
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

func (u *otpUsecase) GetMyOtpInfo(ctx context.Context, input input.MyOtpInfoInput) ([]output.MyOtpInfoOutput, error) {

	en := entity.MakeMyOtpInfoEntity(input.UserId, input.VersionType, input.VersionInfo, input.Uuid)

	result := make([]output.MyOtpInfoOutput, 0)

	if en.VersionType == consts.VERSION_TYPE_LATEST {
		// 최신 버전 - 메모리 부터 조회, 없으면 DB 조회 후 캐싱 처리
		chatTemp, err := u.otpStorage.GetMyOtpInfoStorage(ctx, en, u.svChKeyVersion, consts.CHAT)
		if err != nil {
			if err == consts.ErrOtpNotFound {
				log.Printf("[GetMyOtpInfo] %s GetMyOtpInfoStorage chat Key is empty. check DB\n", en.VersionType)
				temp, err := u.repository.GetMyOtpInfo(ctx, en, consts.CHAT, u.svChKeyVersion)
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						log.Printf("[GetMyOtpInfo] %s GetMyOtpInfoStorage chat key is empty at DB\n", en.VersionType)
					} else {
						return nil, err
					}
				} else {
					u.otpStorage.SaveOtpKeyStorage(ctx, en.UserId, en.Uuid, temp)
					chatTemp.SvKeyVersion = temp.SvKeyVersion
					chatTemp.Kind = temp.Kind
					chatTemp.OtpKey = temp.OtpKey
					chatTemp.OtpRegDate = temp.OtpRegDate
				}
			} else {
				log.Printf("[GetMyOtpInfo] %s GetMyOtpInfoStorage check server key regist.\n", en.VersionType)
			}
		}

		noteTemp, err := u.otpStorage.GetMyOtpInfoStorage(ctx, en, u.svNoKeyVersion, consts.NOTE)
		if err != nil {
			if err == consts.ErrOtpNotFound {
				log.Printf("[GetMyOtpInfo] %s GetMyOtpInfoStorage note key is empty. check DB\n", en.VersionType)
				temp, err := u.repository.GetMyOtpInfo(ctx, en, consts.NOTE, u.svNoKeyVersion)
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						log.Printf("[GetMyOtpInfo] %s GetMyOtpInfoStorage note key is empty at DB\n", en.VersionType)
					} else {
						return nil, err
					}
				} else {
					u.otpStorage.SaveOtpKeyStorage(ctx, en.UserId, en.Uuid, temp)
					noteTemp.SvKeyVersion = temp.SvKeyVersion
					noteTemp.Kind = temp.Kind
					noteTemp.OtpKey = temp.OtpKey
					noteTemp.OtpRegDate = temp.OtpRegDate
				}
			} else {
				log.Printf("[GetMyOtpInfo] %s GetMyOtpInfoStorage check server key regist.\n", en.VersionType)
			}
		}

		noteKey := output.MyOtpInfoOutput{
			Version:    noteTemp.SvKeyVersion,
			KeyType:    noteTemp.Kind,
			Key:        noteTemp.OtpKey,
			OtpRegDate: noteTemp.OtpRegDate,
		}

		chatKey := output.MyOtpInfoOutput{
			Version:    chatTemp.SvKeyVersion,
			KeyType:    chatTemp.Kind,
			Key:        chatTemp.OtpKey,
			OtpRegDate: chatTemp.OtpRegDate,
		}

		result = append(result, noteKey, chatKey)

	} else if en.VersionType == consts.VERSION_TYPE_SPECIFIC {
		// 특정 버전 - DB 조회 TODO

		for i := 0; i < len(en.VersionInfo); i++ {

			temp, err := u.repository.GetMyOtpInfo(ctx, en, consts.NOTE, en.VersionInfo[i])
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					log.Printf("[GetMyOtpInfo] %s GetMyOtpInfoStorage key is empty at DB\n", en.VersionType)
				} else {
					return nil, err
				}
			}

			noteKey := output.MyOtpInfoOutput{
				Version:    temp.SvKeyVersion,
				KeyType:    temp.Kind,
				Key:        temp.OtpKey,
				OtpRegDate: temp.OtpRegDate,
			}

			result = append(result, noteKey)

			temp, err = u.repository.GetMyOtpInfo(ctx, en, consts.CHAT, en.VersionInfo[i])
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					log.Printf("[GetMyOtpInfo] %s GetMyOtpInfoStorage key is empty at DB\n", en.VersionType)
				} else {
					return nil, err
				}
			}
			chatKey := output.MyOtpInfoOutput{
				Version:    temp.SvKeyVersion,
				KeyType:    temp.Kind,
				Key:        temp.OtpKey,
				OtpRegDate: temp.OtpRegDate,
			}
			result = append(result, chatKey)
		}

	} else if en.VersionType == consts.VERSION_TYPE_ALL {
		// 전체 버전 - DB 조회 TODO

	}
	return result, nil
}
