package usecase

import (
	"bytes"
	"core/internal/application/usecase/input"
	"core/internal/domain/appValidation/entity"
	"core/internal/domain/appValidation/repository"
	"core/internal/infrastructure/storage"

	"core/internal/delivery/dto/appValidation"
	"core/pkg/dto"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type appValidationUsecase struct {
	repository        repository.AppValidationRepository
	serverInfoStorage storage.ServerInfoStorage
}

type AppValidationUsecase interface {
	CheckValidation(in input.AppValidationInput) (interface{}, error)
	GetWorksInfos(entity entity.ValidationEntity) (*appValidation.AppValidationResponseDTO, error)
}

func NewAppValidationUsecase(repository repository.AppValidationRepository, serverInfoStorage storage.ServerInfoStorage) AppValidationUsecase {

	return &appValidationUsecase{
		repository:        repository,
		serverInfoStorage: serverInfoStorage,
	}
}

func (u *appValidationUsecase) CheckValidation(in input.AppValidationInput) (interface{}, error) {

	entity := entity.NewAppValidationEntity(in.Hash, in.Device, in.Uuid, in.WorksCode)
	err := u.repository.GetValidation(entity)

	if err != nil {
		return nil, err
	}

	result, err := u.GetWorksInfos(entity)

	return result, nil

}

// entitiy로 변경
func (u *appValidationUsecase) GetWorksInfos(entity entity.ValidationEntity) (*appValidation.AppValidationResponseDTO, error) {

	worksCode := entity.WorksCode

	// 메모리를 통한 조회
	worksCommonInfo := u.serverInfoStorage.GetWorksCommonInfo(worksCode)

	if worksCommonInfo == nil {
		fmt.Printf("[GetWorksCommonInfo] %s's worksCommonInfo is empty... check DB  \n", worksCode)
		// DB 조회를 통한 works 서버 정보 조회
		info, err := u.repository.GetWorksCommonInfo(worksCode)

		if err != nil {
			fmt.Printf("[GetWorksCommonInfo] worksCode %s's worksCommonInfo is DB empty... process end. \n", worksCode)
			return nil, err
		}
		u.serverInfoStorage.SaveWorksCommonInfo(worksCode, info)
		worksCommonInfo = info
	}

	// works의 domain/common API 호출 -> auth 호출 해서 jwt 발급, 저장, 결과 response.
	worksInfo, err := getWorksInfoInCommon(entity, worksCommonInfo.ServerUrl)

	log.Println("worksInfo", worksInfo)

	if err != nil {
		log.Println("common service 호출시 에러 발생함.")
		return nil, err
	}

	return &appValidation.AppValidationResponseDTO{
		Body: appValidation.AppValidationResponseBody{
			WorksCommonInfo: worksCommonInfo,
			WorksInfo:       worksInfo,
		},
	}, nil

}

func toWorksInfoRequestDtoBody(uuid string, device string, worksCode string) appValidation.DeviceInitRequestBody {

	return appValidation.DeviceInitRequestBody{
		Uuid:      uuid,
		WorksCode: worksCode,
		Device:    device,
	}
}

func getWorksInfoInCommon(entity entity.ValidationEntity, serverUrl string) (*appValidation.DeviceInitResponseBody, error) {

	// server to server DTO 정의
	header := appValidation.DeviceInitRequestHeader{
		ServerToken: "",
	}

	reqDto := appValidation.DeviceInitRequestDTO{
		Header: header,
		Body:   body,
	}

	// url 정의
	url := "http://" + serverUrl + "/common/sv1/device-init"
	log.Println("common service 호출! url : ", url)

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

	client := &http.Client{}
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

	responseDTO := result.Data

	log.Println("common service 호출 end !")
	return responseDTO.Body, nil
}
