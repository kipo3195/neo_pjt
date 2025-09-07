package handler

import (
	"auth/internal/application/usecase"

	"github.com/gin-gonic/gin"
)

type UserAuthHandler struct {
	usecase usecase.UserAuthUsecase
}

func NewUserAuthHandler(uc usecase.UserAuthUsecase) UserAuthHandler {
	return UserAuthHandler{
		usecase: uc,
	}
}

func (h UserAuthHandler) GenerateAuthChallenge(c *gin.Context) {

}

func (h UserAuthHandler) GetAuthStatus(c *gin.Context) {

}
