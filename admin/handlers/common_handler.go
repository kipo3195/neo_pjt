package handlers

import (
	"admin/consts"
	commonReqDto "admin/dto/client/common/request"
	"admin/usecases"
	"io"
	"log"

	"github.com/gin-gonic/gin"
)

// 관리자 - common service 연계 handler
type CommonHandler struct {
	usecase usecases.CommonUsecase
}

func NewCommonHandler(uc usecases.CommonUsecase) *CommonHandler {
	return &CommonHandler{usecase: uc}
}

func (h *CommonHandler) CreateSkinImg(c *gin.Context) {

	// context 생성 - admin_route에 정의된 middleware에서 context에 관여함.
	ctx := c.Request.Context()

	// 역할 분리:
	// 계층	역할
	// Handler	파라미터 파싱, 최소 유효성 검사 (파일 존재 여부, 너무 큰 파일 여부)
	// Usecase	비즈니스 로직 관점의 검증 (파일 확장자, 내용 유형 체크 등)
	// Repository/Infra	외부 전송, 저장 등 처리 수행

	// 파일 데이터 추출
	fileInfo, err := c.FormFile(consts.FILE)
	if err != nil {
		// 파일
		sendErrorResponse(c, consts.BAD_REQUEST, consts.FAIL, consts.ADMIN_F001, consts.ADMIN_F001_MSG)
	}

	skinType := c.GetHeader(consts.SKIN_TYPE)

	if skinType == "" {
		sendErrorResponse(c, consts.BAD_REQUEST, consts.ERROR, consts.E_104, consts.E_104_MSG)
		return
	}

	// 파일 처리 []byte -> dto
	file, err := fileInfo.Open()
	defer file.Close()
	if err != nil {
		sendErrorResponse(c, consts.SERVER_ERROR, consts.ERROR, consts.E_500, consts.E_500_MSG)
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		sendErrorResponse(c, consts.SERVER_ERROR, consts.ERROR, consts.E_500, consts.E_500_MSG)
	}

	body := commonReqDto.CreateSkinImgRequestBody{
		SkinType: skinType,
		File:     fileBytes,
		FileSize: fileInfo.Size,
		FileName: fileInfo.Filename,
	}

	requestDTO := commonReqDto.CreateSkinImgRequestDTO{
		Body: body,
	}

	err = h.usecase.CreateSkinImg(ctx, requestDTO.Body)

	log.Println("admin 서비스 스킨 이미지 업로드에 대한 response : ", err)

	if err != nil {

		if err == consts.ErrFileSizeExceeded {
			sendErrorResponse(c, consts.BAD_REQUEST, consts.FAIL, consts.ADMIN_F002, consts.ADMIN_F002_MSG)
		} else if err == consts.ErrFileExtentionDetect {
			sendErrorResponse(c, consts.BAD_REQUEST, consts.FAIL, consts.ADMIN_F003, consts.ADMIN_F003_MSG)
		} else {
			sendErrorResponse(c, consts.SERVER_ERROR, consts.ERROR, consts.E_500, consts.E_500_MSG)
		}
	} else {
		sendSuccessResponse(c, "")
	}
}
