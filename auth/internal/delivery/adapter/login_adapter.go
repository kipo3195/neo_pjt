package adapter

import (
	"auth/internal/application/usecase/input"
	"auth/internal/delivery/dto/certification"
)

func MakeLoginInput(req certification.AuthRequestDTO) input.LoginInput {

	return input.LoginInput{
		Id:       req.Body.Id,
		Password: req.Body.Password,
		Token:    req.Header.Token,
		Uuid:     req.Header.Uuid,
	}

}
