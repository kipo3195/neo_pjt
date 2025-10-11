package handler

import (
	"log"
	"org/internal/application/usecase"
	"org/internal/application/usecase/input"
	"org/internal/consts"
	"org/internal/delivery/dto/user"
	commonConsts "org/pkg/consts"
	"org/pkg/response"

	"github.com/gin-gonic/gin"
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
	id := c.Value(consts.USER_ID)
	myHash, ok := id.(string)
	if !ok {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.ORG_F101, consts.ORG_F101_MSG)
		return
	}

	log.Println("[GetMyInfo] myHash : ", myHash)

	// dto 생성
	myInfoInput := input.MakeMyInfoInput(myHash)
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

}
