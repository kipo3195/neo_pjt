package client

import (
	"io"
	"log"

	"admin/internal/consts"
	"admin/internal/domains/skinImg/dto/client/requestDTO"
	clientUsecases "admin/internal/domains/skinImg/usecases/client"
	commonConsts "admin/pkg/consts"
	response "admin/pkg/response"

	"github.com/gin-gonic/gin"
)

type SkinImgHandler struct {
	usecase clientUsecases.SkinImgUsecase
}

func NewSkinImgHandler(usecase clientUsecases.SkinImgUsecase) *SkinImgHandler {
	return &SkinImgHandler{usecase: usecase}
}

func (h *SkinImgHandler) CreateSkinImg(c *gin.Context) {

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
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.ADMIN_F001, consts.ADMIN_F001_MSG)
		return
	}

	skinType := c.GetHeader(consts.SKIN_TYPE)

	if skinType == "" {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_104, commonConsts.E_104_MSG)
		return
	}

	// 파일 처리 []byte -> dto
	file, err := fileInfo.Open()
	defer file.Close()
	if err != nil {
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		return
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		return
	}

	body := requestDTO.CreateSkinImgRequestBody{
		SkinType: skinType,
		File:     fileBytes,
		FileSize: fileInfo.Size,
		FileName: fileInfo.Filename,
	}

	requestDTO := requestDTO.CreateSkinImgRequestDTO{
		Body: body,
	}

	err = h.usecase.CreateSkinImg(ctx, requestDTO.Body)

	log.Println("admin 서비스 스킨 이미지 업로드에 대한 response : ", err)

	if err != nil {

		if err == consts.ErrFileSizeExceeded {
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.ADMIN_F002, consts.ADMIN_F002_MSG)
		} else if err == consts.ErrFileExtentionDetect {
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.ADMIN_F003, consts.ADMIN_F003_MSG)
		} else {
			response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		}
	} else {
		response.SendSuccess(c, "")
	}
}
