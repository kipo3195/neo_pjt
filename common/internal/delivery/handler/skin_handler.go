package handler

import (
	"common/internal/application/usecase"
	"common/internal/delivery/dto/skin"
	commonConsts "common/pkg/consts"
	"common/pkg/response"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type SkinHandler struct {
	usecase usecase.SkinUsecase
}

func NewSkinHandler(usecase usecase.SkinUsecase) *SkinHandler {
	return &SkinHandler{
		usecase: usecase,
	}
}

func (h *SkinHandler) GetSkinImage(c *gin.Context) {
	// context 생성
	ctx := c.Request.Context()

	// 데이터 -> dto
	var req = skin.GetSkinImgRequest{
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

func (h *SkinHandler) PutSkinImg(c *gin.Context) {

	// context 생성
	ctx := c.Request.Context()

	log.Println("111")

	// 파일 데이터 추출
	file, fileInfo, err := c.Request.FormFile("File")
	skinType := c.GetHeader("Skin-Type")

	if err != nil || file == nil || fileInfo.Size == 0 || skinType == "" {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	var req = skin.CreateSkinImgRequest{
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
