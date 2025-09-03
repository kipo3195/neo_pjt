package repository

import (
	"bytes"
	"context"
	"core/internal/domain/appValidation/entity"
	"core/internal/domain/appValidation/repository"
	"core/internal/infrastructure/dto/appValidation"
	"core/pkg/dto"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type appValidationAPIRepository struct {
	serverUrl  string
	httpClient *http.Client
}

func NewAppValidationAPIRepository(serverUrl string) repository.AppValidationAPIRepository {
	return &appValidationAPIRepository{
		serverUrl:  serverUrl,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}

}

func (c *appValidationAPIRepository) DeviceInit(ctx context.Context, e entity.ValidationEntity) (*entity.DeviceInitResult, error) {

	url := "http://" + c.serverUrl + "/common/sv1/device-init"
	log.Println("common service 호출! url : ", url)

	header := appValidation.DeviceInitRequestHeader{
		ServerToken: "",
	}

	reqDto := appValidation.DeviceInitRequestDTO{
		Header: header,
		Body: appValidation.DeviceInitRequestBody{
			Uuid:      e.Uuid,
			WorksCode: e.WorksCode,
			Device:    e.Device,
		},
	}

	// 직렬화
	bodyByte, err := json.Marshal(reqDto.Body)
	if err != nil {
		// 에러 처리
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyByte))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", reqDto.Header.ServerToken) // works 서버 호출시 필요한 키 작성하기 TODO

	client := c.httpClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// serverResponse로 전달받기 -> dto 뽑아내기 제네릭
	var result dto.ServerResponseDTO[*appValidation.DeviceInitResponseDTO]

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("serverReponse 파싱시 에러")
		return nil, err
	}

	log.Println("common service 호출 end !")

	return &entity.DeviceInitResult{
		IssuedAppToken: (*entity.IssuedAppToken)(result.Data.Body.IssuedAppToken),
		WorksConfig:    result.Data.Body.WorksConfig,
	}, nil
}
