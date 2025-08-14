package server

import (
	requestDTO "common/internal/domains/skin/dto/server/requestDTO"
	serverUsecase "common/internal/domains/skin/usecases/server"
	commonConsts "common/pkg/consts"
	"common/pkg/response"
	"log"

	"github.com/gin-gonic/gin"
)

type SkinHandler struct {
	usecase serverUsecase.SkinUsecase
}

func NewSkinHandler(usecase serverUsecase.SkinUsecase) *SkinHandler {
	return &SkinHandler{
		usecase: usecase,
	}
}

func (h *SkinHandler) PutSkinImg(c *gin.Context) {

	// context 생성
	ctx := c.Request.Context()

	log.Println("111")

	// 파일 데이터 추출
	file, fileInfo, err := c.Request.FormFile("File")
	skinType := c.GetHeader("Skin-Type")

	if err != nil || file == nil || fileInfo.Size == 0 || skinType == "" {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
	}

	var req = requestDTO.CreateSkinImgRequest{
		File:     file,
		FileInfo: fileInfo,
		SkinType: skinType,
	}

	data, err := h.usecase.CreateSkinImg(ctx, req)

	log.Println(data)

	if err == nil {
		// http status code 200
		response.SendSuccess(c, data)
	} else {
		// TODO 모든 에러 세분화 할 것.
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
	}

}
