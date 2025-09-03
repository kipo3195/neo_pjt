package adapter

import (
	"core/internal/application/usecase/input"
	"core/internal/delivery/dto/appValidation"
)

func MakeValidateAppInput(req appValidation.AppValidationRequestDTO) input.AppValidationInput {

	return input.AppValidationInput{
		Hash:      req.Header.Hash,
		Device:    req.Header.Device,
		Uuid:      req.Header.Uuid,
		WorksCode: req.Body.WorksCode,
	}

}
