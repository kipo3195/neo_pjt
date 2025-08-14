package client

import (
	"common/internal/domains/skin/dto/client/requestDTO"
	clientUsecase "common/internal/domains/skin/usecases/client"
	commonConsts "common/pkg/consts"
	"common/pkg/response"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type SkinHandler struct {
	usecase clientUsecase.SkinUsecase
}

func NewSkinHandler(usecase clientUsecase.SkinUsecase) *SkinHandler {
	return &SkinHandler{
		usecase: usecase,
	}
}

func (h *SkinHandler) GetSkinImage(c *gin.Context) {
	// context 생성
	ctx := c.Request.Context()

	// 데이터 -> dto
	var req = requestDTO.GetSkinImgRequest{
		SkinHash: c.Query("skinHash"),
		SkinType: c.Query("skinType"),
	}

	log.Println("2")
	// 유효성 검증 로직
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_108, commonConsts.E_108_MSG)
		return
	}

	log.Println("3")
	// 검증
	file, err := h.usecase.GetSkinImg(ctx, req)

	// TODO send file

	defer file.Close()
	if err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		return
	}

	response.SendFileStream(c, file, "", "")
}
