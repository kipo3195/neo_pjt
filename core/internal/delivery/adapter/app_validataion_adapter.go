package adapter

import (
	"core/internal/application/usecase/input"
	"core/internal/application/usecase/output"
	"core/internal/delivery/dto/appValidation"
	"core/internal/domain/appValidation/entity"
)

func MakeValidateAppInput(req appValidation.AppValidationRequestDTO) input.AppValidationInput {

	return input.AppValidationInput{
		Hash:      req.Header.Hash,
		Device:    req.Header.Device,
		Uuid:      req.Header.Uuid,
		WorksCode: req.Body.WorksCode,
	}

}

func MakeValidateAppOutput(en *entity.DeviceInitResult) output.AppValidationOutput {
	return output.AppValidationOutput{
		// 데이터 정의 필요
	}
}
