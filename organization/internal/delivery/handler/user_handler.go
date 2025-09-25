package handler

import (
	"context"
	"encoding/json"
	"log"
	"org/internal/application/usecase"
	"org/internal/application/usecase/input"
	"org/internal/consts"
	"org/internal/delivery/dto/user"
	"org/internal/delivery/middleware/contextkey"
	commonConsts "org/pkg/consts"
	"org/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type UserHandler struct {
	usecase usecase.UserUsecase
}

func NewUserHandler(usecase usecase.UserUsecase) *UserHandler {

	return &UserHandler{
		usecase: usecase,
	}
}

func (h *UserHandler) GetMyInfo(c *gin.Context) {

	// context 생성
	ctx := c.Request.Context()

	// 인증 토큰에서 요청 사용자의 hash 정보 추출
	val := c.Value(contextkey.UserHashKey)
	myHash, ok := val.(string)
	if !ok {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, consts.ORG_F101, consts.ORG_F101_MSG)
		log.Println("인증 토큰의 userHash 데이터 에러 ")
		return
	}

	var req user.GetMyInfoRequest

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

	log.Println("내 정보 요청시 myHash : ", myHash)
	// dto 생성
	myInfoInput := input.MakeMyInfoInput(req.MyHash)
	output, err := h.usecase.GetMyInfo(ctx, myInfoInput)

	username := user.UsernameDto{
		Def: output.Username.Ko, // 수정 필요
		Ko:  output.Username.Ko,
		En:  output.Username.En,
		Jp:  output.Username.Jp,
		Zh:  output.Username.Zh,
		Ru:  output.Username.Ru,
		Vi:  output.Username.Vi,
	}

	var deptInfo []user.DeptInfoDto

	for _, temp := range output.DeptInfo {
		deptInfo = append(deptInfo, user.DeptInfoDto{
			DeptOrg:  temp.DeptOrg,
			DeptCode: temp.DeptCode,
			DefLang:  temp.DefLang,
			KoLang:   temp.KoLang,
			EnLang:   temp.EnLang,
			JpLang:   temp.JpLang,
			ZhLang:   temp.ZhLang,
			ViLang:   temp.ViLang,
			RuLang:   temp.RuLang,
			Header:   temp.Header,
		})
	}

	res := user.GetMyInfoResponse{
		UserHash:     output.UserHash,
		UserPhoneNum: output.UserPhoneNum,
		Username:     username,
		OrgCodes:     nil,
		ProfileUrl:   output.ProfileUrl,
		ProfileMsg:   output.ProfileMsg,
		DeptInfo:     deptInfo,
	}

	// response.
	if err == nil {
		// http status code 200
		response.SendSuccess(c, res)
	} else {
		// http status code 400
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
	}

}

func (h *UserHandler) GetUserInfo(c *gin.Context) {

	// context 생성
	ctx := c.Request.Context()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

}
