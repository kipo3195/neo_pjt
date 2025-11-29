package handler

import (
	"encoding/json"
	"io"
	"log"
	"user/internal/application/usecase"
	"user/internal/consts"
	"user/internal/delivery/adapter"
	"user/internal/delivery/dto/profile"
	"user/internal/delivery/util"
	commonConsts "user/pkg/consts"
	"user/pkg/response"

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

	userHash := util.GetUserIdByAccessToken(c)
	userId := util.GetUserHashByAccessToken(c)

	// 테스트 용 -> 다른 사람꺼 등록가능
	headerUserHash := c.GetHeader("User-Hash")
	if headerUserHash != "" {
		//response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_110, commonConsts.E_110_MSG)
		userHash = headerUserHash
		userId = "TEST_USER_ID"
	}
	// 테스트용 끝 -> 이후 제거될 코드

	log.Print("userHash : ", userHash)
	log.Print("userId : ", userId)

	// 파일 데이터 추출
	fileInfo, err := c.FormFile("profile_img")
	if err != nil {
		// 파일
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
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
	req := profile.ProfileUploadRequest{
		ProfileImg:     &fileBytes,
		ProfileImgSize: fileInfo.Size,
		ProfileImgName: fileInfo.Filename,
	}

	input := adapter.MakeProfileUploadInput(req.ProfileImg, req.ProfileImgSize, req.ProfileImgName, userId, userHash)
	err = h.usecase.ProfileImgUpload(ctx, input)

	if err != nil {
		if err == consts.ErrFileSizeExceeded {
			// 사이즈 에러
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.USER_PROFILE_F001, consts.USER_PROFILE_F001_MSG)
		} else if err == consts.ErrFileExtentionDetect {
			// 확장자 에러
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.USER_PROFILE_F002, consts.USER_PROFILE_F002_MSG)
		} else if err == consts.ErrProfileImgSaveError {
			// 서버 저장에러
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.USER_PROFILE_F003, consts.USER_PROFILE_F003_MSG)
		} else if err == consts.ErrProfileImgDBSaveError {
			// DB 저장 에러
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.USER_PROFILE_F004, consts.USER_PROFILE_F004_MSG)
		} else {
			response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		}
	} else {
		response.SendSuccess(c, "")
	}
}

func (h *ProfileHandler) GetProfileImg(c *gin.Context) {

	ctx := c.Request.Context()

	var req profile.GetProfileImgRequest

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	// 필수 데이터 검증
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_108, commonConsts.E_108_MSG)
		return
	}

	getProfileImgInput := adapter.MakeGetProfileImgInput(req.UserHash)
	output, err := h.usecase.GetProfileImg(ctx, getProfileImgInput)
	if err != nil {
		if err == consts.ErrProfileImgNotRegist {
			response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, consts.USER_PROFILE_F005, consts.USER_PROFILE_F005_MSG)
		} else {
			response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		}
	} else {
		response.SendFileStream(c, output.ProfileImg, output.ProfileImgName, "")
	}

}

func (h *ProfileHandler) DeleteProfileImg(c *gin.Context) {

	ctx := c.Request.Context()

	userId := util.GetUserIdByAccessToken(c)
	if userId == "" {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_108, commonConsts.E_108_MSG)
		return
	}

	req := profile.DeleteProfileImgRequest{
		UserId: userId,
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_108, commonConsts.E_108_MSG)
		return
	}

	deleteProfileImgInput := adapter.MakeDeleteProfileImgInput(req.UserId)
	err := h.usecase.DeleteProfileImg(ctx, deleteProfileImgInput)

	if err != nil {
		if err == consts.ErrProfileImgNotRegist || err == consts.ErrProfileImgDBDeleteError {
			// 메모리, 서버에 없음
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.USER_PROFILE_F005, consts.USER_PROFILE_F005_MSG)
		} else {
			// server error
			response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		}
		return
	}

	response.SendSuccess(c, "")
}

func (h *ProfileHandler) RegistProfileMsg(c *gin.Context) {

	ctx := c.Request.Context()

	userHash := util.GetUserHashByAccessToken(c)
	if userHash == "" {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_108, commonConsts.E_108_MSG)
		return
	}

	var req profile.RegistProfileMsgRequest

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_108, commonConsts.E_108_MSG)
		return
	}

	input := adapter.MakeRegistProfileMsgInput(userHash, req.ProfileMsg)
	err := h.usecase.RegistProfileMsg(ctx, input)

	if err != nil {
		// server error
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)

	} else {
		response.SendSuccess(c, "")
	}

}

func (h *ProfileHandler) GetProfileMsg(c *gin.Context) {

	ctx := c.Request.Context()

	var req profile.GetProfileMsgRequest

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_108, commonConsts.E_108_MSG)
		return
	}

	input := adapter.MakeGetProfileMsgInput(req.UserHashs)
	output, err := h.usecase.GetProfileMsg(ctx, input)

	res := profile.GetProfileMsgResponse{
		ProfileMsgInfos: []profile.ProfileMsgInfo{}, // 여기서 빈 배열 초기화
	}

	if err != nil {
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		return
	} else {

		for i := 0; i < len(output.ProfileMsg); i++ {
			profileMsgInfo := profile.ProfileMsgInfo{
				UserHash:   output.ProfileMsg[i].UserHash,
				ProfileMsg: output.ProfileMsg[i].ProfileMsg,
			}
			res.ProfileMsgInfos = append(res.ProfileMsgInfos, profileMsgInfo)
		}
	}

	response.SendSuccess(c, res)

}
