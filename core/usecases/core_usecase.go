package usecases

import (
	"bytes"
	clReqDto "core/dto/client/request"
	clResDto "core/dto/client/response"
	dto "core/dto/common"
	svCommonReqDto "core/dto/server/common/request"
	svCommonResDto "core/dto/server/common/response"
	"core/entities"
	storage "core/infra/storage"
	"core/repositories"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type coreUsecase struct {
	repo              repositories.CoreRepository
	serverInfoStorage storage.ServerInfoStorage
}

type CoreUsecase interface {
	CheckValidation(header clReqDto.AppValidationRequestHeader) (bool, error)
	GetWorksInfos(reqDto clReqDto.AppValidationRequestDTO) (*clResDto.AppValidationResponseDTO, error)
}

func NewCoreUsecase(repo repositories.CoreRepository, serverInfoStorage storage.ServerInfoStorage) CoreUsecase {
	return &coreUsecase{
		repo:              repo,
		serverInfoStorage: serverInfoStorage,
	}
}

func (u *coreUsecase) CheckValidation(header clReqDto.AppValidationRequestHeader) (bool, error) {

	return u.repo.GetValidation(toValidationEntity(header))

}
func toValidationEntity(header clReqDto.AppValidationRequestHeader) entities.ValidationEntity {
	return entities.ValidationEntity{
		Hash:   header.Hash,
		Device: header.Device,
	}
}

func (u *coreUsecase) GetWorksInfos(reqDto clReqDto.AppValidationRequestDTO) (*clResDto.AppValidationResponseDTO, error) {

	var worksCommonInfo *entities.WorksCommonInfo
	worksCode := reqDto.Body.WorksCode

	// 메모리를 통한 조회
	worksCommonInfo = u.serverInfoStorage.GetWorksCommonInfo(worksCode)

	if worksCommonInfo == nil {
		fmt.Printf("[GetWorksCommonInfo] %s's worksCommonInfo is empty... check DB  \n", worksCode)
		// DB 조회를 통한 works 서버 정보 조회
		info, err := u.repo.GetWorksCommonInfo(worksCode)

		if err != nil {
			fmt.Printf("[GetWorksCommonInfo] worksCode %s's worksCommonInfo is DB empty... process end. \n", worksCode)
			return nil, err
		}
		u.serverInfoStorage.SaveWorksCommonInfo(worksCode, info)
		worksCommonInfo = info
	}

	// works의 domain/common API 호출 -> auth 호출 해서 jwt 발급, 저장, 결과 response.
	worksInfo, err := getWorksInfoInCommon(toWorksInfoRequestDtoBody(reqDto.Header.Uuid, reqDto.Header.Device, reqDto.Body.WorksCode), worksCommonInfo.ServerUrl)

	log.Println("worksInfo", worksInfo)

	if err != nil {
		log.Println("common service 호출시 에러 발생함.")
		return nil, err
	}

	return &clResDto.AppValidationResponseDTO{
		Body: clResDto.AppValidationResponseBody{
			WorksCommonInfo: worksCommonInfo,
			WorksInfo:       worksInfo,
		},
	}, nil

}

func toWorksInfoRequestDtoBody(uuid string, device string, worksCode string) svCommonReqDto.DeviceInitRequestBody {

	return svCommonReqDto.DeviceInitRequestBody{
		Uuid:      uuid,
		WorksCode: worksCode,
		Device:    device,
	}
}

func getWorksInfoInCommon(body svCommonReqDto.DeviceInitRequestBody, serverUrl string) (*svCommonResDto.DeviceInitResponseBody, error) {

	// server to server DTO 정의
	header := svCommonReqDto.DeviceInitRequestHeader{
		ServerToken: "",
	}

	reqDto := svCommonReqDto.DeviceInitRequestDTO{
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
	var result dto.ServerResponseDTO[*svCommonResDto.DeviceInitResponseBody]

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("serverReponse 파싱시 에러")
		return nil, err
	}

	res := result.Data

	log.Println("common service 호출 end !")
	return res, nil
}
