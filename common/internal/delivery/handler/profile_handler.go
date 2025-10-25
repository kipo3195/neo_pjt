package handler

import (
	"common/internal/application/usecase"
	"common/internal/application/usecase/input"
	"common/internal/consts"
	"common/internal/delivery/dto/profile"
	"common/internal/delivery/util"
	commonConsts "common/pkg/consts"
	"common/pkg/response"
	"io"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
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

	ctx := c.Request.Context()
	log.Println("111")
	userId := util.GetUserIdByAccessToken(c)
	if userId == "" {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_110, commonConsts.E_110_MSG)
		return
	}

	log.Println("2222")
	// 파일 데이터 추출
	fileInfo, err := c.FormFile("profile_img")
	if err != nil {
		// 파일
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_104, commonConsts.E_104_MSG)
		return
	}

	log.Println("333")
	// 파일 처리 []byte -> dto
	file, err := fileInfo.Open()
	defer file.Close()
	if err != nil {
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		return
	}
	log.Println("4444")

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		return
	}
	log.Println("5555")
	req := profile.ProfileUploadRequest{
		ProfileImg:     &fileBytes,
		ProfileImgSize: fileInfo.Size,
		ProfileImgName: fileInfo.Filename,
	}
	log.Println("6666")
	input := input.MakeProfileUploadInput(req.ProfileImg, req.ProfileImgSize, req.ProfileImgName, userId)
	err = h.usecase.ProfileImgUpload(ctx, input)
	log.Println("7777")
	if err != nil {
		if err == consts.ErrFileSizeExceeded {
			// 사이즈 에러
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.COMMON_PROFILE_F001, consts.COMMON_PROFILE_F001_MSG)
		} else if err == consts.ErrFileExtentionDetect {
			// 확장자 에러
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.COMMON_PROFILE_F002, consts.COMMON_PROFILE_F002_MSG)
		} else if err == consts.ErrProfileImgSaveError {
			// 서버 저장에러
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.COMMON_PROFILE_F003, consts.COMMON_PROFILE_F003_MSG)
		} else if err == consts.ErrProfileImgDBSaveError {
			// DB 저장 에러
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.COMMON_PROFILE_F004, consts.COMMON_PROFILE_F004_MSG)
		} else {
			response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		}
	} else {
		response.SendSuccess(c, "")
	}

}

func (h *ProfileHandler) GetProfileImg(c *gin.Context) {

	ctx := c.Request.Context()

	userId := c.Query("userId")
	if userId == "" {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_108, commonConsts.E_108_MSG)
		return
	}

	req := profile.GetProfileImgRequest{
		UserId: userId,
	}

	// 필수 데이터 검증
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_108, commonConsts.E_108_MSG)
		return
	}

	getProfileImgInput := input.MakeGetProfileImgInput(req.UserId)
	output, err := h.usecase.GetProfileImg(ctx, getProfileImgInput)
	if err != nil {
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
	} else {
		response.SendFileStream(c, output.ProfileImg, output.ProfileImgName, "")
	}

}

func (h *ProfileHandler) DeleteProfileImg(c *gin.Context) {

}

func (h *ProfileHandler) RegistProfileMsg(c *gin.Context) {

}
