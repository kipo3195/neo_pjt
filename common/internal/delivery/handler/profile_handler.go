package handler

import (
	"common/internal/application/usecase"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	usecase usecase.ProfileUsecase
}

func NewProfileHandler(usecase usecase.ProfileUsecase) *ProfileHandler {
	return &ProfileHandler{
		usecase: usecase,
	}
}

func (h *ProfileHandler) UploadProfileImg(c *gin.Context) {

}

func (h *ProfileHandler) DeleteProfileImg(c *gin.Context) {

}

func (h *ProfileHandler) RegistProfileMsg(c *gin.Context) {

}
