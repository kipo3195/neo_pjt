package client

import "org/usecases"

type UserHandler struct {
	usecase usecases.UserUsecase
}

func NewUserHandler(usecase usecases.UserUsecase) *UserHandler {

	return &UserHandler{
		usecase: usecase,
	}
}
