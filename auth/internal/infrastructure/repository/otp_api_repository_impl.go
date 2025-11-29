package repository

import (
	"auth/internal/consts"
	"auth/internal/domain/otp/entity"
	"auth/internal/domain/otp/repository"
	"auth/internal/infrastructure/dto/otp"
	"auth/pkg/dto"
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type otpApiRepositoryImpl struct {
	serverUrl  string
	httpClient *http.Client
}

func NewOtpApiRepository(serverUrl string) repository.OtpApiRepository {
	return &otpApiRepositoryImpl{
		serverUrl:  serverUrl,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (r *otpApiRepositoryImpl) OtpKeyRegistInMessage(ctx context.Context, en entity.OtpKeyRegistEntity) (entity.OtpKeyRegistResultEntity, error) {

	// message API 호출
	url := "http://" + r.serverUrl + "/message/server/v1/otp/regist"
	log.Println("message service 호출! url : ", url)

	devicePubKeyDto := make([]otp.DevicePubKeyDto, 0)

	for i := 0; i < len(en.DevicePubKeyEntity); i++ {
		d := otp.DevicePubKeyDto{
			Kind: en.DevicePubKeyEntity[i].Kind,
			Key:  en.DevicePubKeyEntity[i].Key,
		}
		devicePubKeyDto = append(devicePubKeyDto, d)
	}

	// request dto
	// adapter에서 처리하지 하는 것이 오버엔지니어링이라고 판단함.
	body := otp.OtpKeySvRegistRequest{
		Id:              en.Id,
		Uuid:            en.Uuid,
		DevicePubKeyDto: devicePubKeyDto,
	}

	// 직렬화
	bodyByte, err := json.Marshal(body)
	if err != nil {
		// 에러 처리
		return entity.OtpKeyRegistResultEntity{}, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyByte))
	if err != nil {
		return entity.OtpKeyRegistResultEntity{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "") // works 서버 호출시 필요한 키 작성하기 TODO

	client := r.httpClient
	resp, err := client.Do(req)
	if err != nil {
		return entity.OtpKeyRegistResultEntity{}, err
	}

	defer resp.Body.Close()

	resultCode := resp.StatusCode

	if resultCode != http.StatusOK {
		log.Println("message service OtpKeyRegistInMessage error : ", resultCode)
		return entity.OtpKeyRegistResultEntity{}, consts.ErrHttpStatusError
	}

	// response dto
	var result dto.ServerResponseDTO[*otp.OtpKeySvRegistResponse]

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("serverReponse 파싱시 에러")
		return entity.OtpKeyRegistResultEntity{}, err
	}

	// dto -> entity
	// adapter에서 처리하지 하는 것이 오버엔지니어링이라고 판단함.
	return entity.OtpKeyRegistResultEntity{
		OtpRegDate:       result.Data.OtpRegDate,
		SvChatKeyVersion: result.Data.SvChatKeyVersion,
		SvNoteKeyVersion: result.Data.SvNoteKeyVersion,
	}, nil
}
