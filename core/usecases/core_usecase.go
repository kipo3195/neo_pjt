package usecases

import (
	"bytes"
	consts "core/consts"
	clDto "core/dto/client"
	dto "core/dto/common"
	commonDto "core/dto/server/common"
	"core/entities"
	"core/repositories"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type coreUsecase struct {
	repo repositories.CoreRepository
}

type CoreUsecase interface {
	CheckValidation(header *clDto.AppValidationRequestHeader) bool

	GetWorksInfos(body clDto.AppValidationRequest, uuid string, device string) (*clDto.WorksInfos, *dto.ErrorResponse, bool)

	// 변환
	ToValidationWhereEntity(header *clDto.AppValidationRequestHeader) entities.ValidationWhere
}

func NewCoreUsecase(repo repositories.CoreRepository) CoreUsecase {
	return &coreUsecase{repo: repo}
}

func (u *coreUsecase) CheckValidation(header *clDto.AppValidationRequestHeader) bool {

	validationWhere := u.ToValidationWhereEntity(header)

	flag, err := u.repo.GetValidation(validationWhere)

	if !flag || err != nil {
		return false
	}

	return true
}

func (u *coreUsecase) ToValidationWhereEntity(header *clDto.AppValidationRequestHeader) entities.ValidationWhere {
	return entities.ValidationWhere{
		Hash:   header.Hash,
		Device: header.Device,
	}
}

func (u *coreUsecase) GetWorksInfos(body clDto.AppValidationRequest, uuid string, device string) (*clDto.WorksInfos, *dto.ErrorResponse, bool) {

	// 에러타입이 뭐냐에따라 처리..
	worksCommonInfo, err := u.repo.GetWorksCommonInfo(body)
	fmt.Println("[GetWorksCommonInfo] worksCommonInfo : ", worksCommonInfo)

	if err != nil {
		switch {
		case errors.Is(err, consts.ErrInvalidType):
			// 타입에러
			return nil, &dto.ErrorResponse{
				Code:    consts.E_101,
				Message: consts.E_101_MSG,
			}, false
		case errors.Is(err, consts.ErrInvalidMappingServer):
			// 매핑된 서버 정보 없음
			return nil, &dto.ErrorResponse{
				Code:    consts.CORE_F102,
				Message: consts.CORE_F102_MSG,
			}, true
		default:
			// 기타 DB 에러
			return nil, &dto.ErrorResponse{
				Code:    consts.E_102,
				Message: consts.E_102_MSG,
			}, false
		}

	} else {

		// works의 domain/common API 호출 -> auth 호출 해서 jwt 발급, 저장, 결과 response.
		worksInfo, err := getWorksInfo(uuid, device, body.WorksCode, worksCommonInfo.ServerUrl)

		fmt.Println("worksInfo", worksInfo)

		if err != nil {
			fmt.Println("common service 호출시 에러 발생함.")
			return nil, &dto.ErrorResponse{
				Code:    consts.E_500,
				Message: consts.E_500_MSG,
			}, false
		}

		return &clDto.WorksInfos{
			WorksCommonInfo: worksCommonInfo,
			WorksInfo:       worksInfo,
		}, nil, false
	}

}

func getWorksInfo(uuid string, device string, worksCode string, serverUrl string) (*commonDto.DeviceInitResponse, error) {
	// 소스 모듈화 처리하기
	data := map[string]string{
		"uuid":      uuid,
		"worksCode": worksCode,
		"device":    device,
	}

	// JSON 변환
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// POST 요청 보내기
	url := "http://" + serverUrl + "/common/sv1/device-init"

	fmt.Println("common service 호출! url : ", url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer serverToken") // works 서버 호출시 필요한 키 작성하기 TODO

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// serverResponse로 전달받기 -> dto 뽑아내기 제네릭
	var result dto.ServerResponse[*commonDto.DeviceInitResponse]

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("serverReponse 파싱시 에러")
		return nil, err
	}

	res := result.Data

	fmt.Println("common service 호출 end !")
	return res, nil
}
