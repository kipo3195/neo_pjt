package handlers

import (
	"org/config"
	"org/usecases"
)

type ServerHandler struct {
	usecase usecases.ServerUsecase
	sfg     *config.ServerConfig
}

func NewServerHandler(sfg *config.ServerConfig, uc usecases.ServerUsecase) *ServerHandler {
	return &ServerHandler{usecase: uc, sfg: sfg}
}
