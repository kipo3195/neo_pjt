package handler

import (
	"encoding/json"
	"log"
	"org/internal/application/usecase"
	"org/internal/application/usecase/input"
	"org/internal/consts"
	"org/internal/delivery/adapter"
	"org/internal/delivery/dto/user"
	commonConsts "org/pkg/consts"
	"org/pkg/response"

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

	userName := user.UsernameDto{
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

	profile := user.UserProfile{
		ProfileUrl: output.ProfileUrl,
		ProfileMsg: output.ProfileMsg,
	}

	userDetail := user.UserDetail{
		UserEmail:    output.UserEmail,
		UserPhoneNum: output.UserPhoneNum,
	}

	res := user.GetMyInfoResponse{
		UserHash:   output.UserHash,
		UserDetail: userDetail,
		Username:   userName,
		OrgCode:    nil,
		Profile:    profile,
		DeptInfo:   deptInfo,
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

	//context 생성
	ctx := c.Request.Context()

	// request 생성
	var req user.GetUserInfoRequest

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

	input := adapter.MakeGetUserInfoInput(req.UserIds)
	output, err := h.usecase.GetUserInfo(ctx, input)

	if err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		return
	}

	var res []user.GetUserInfoResponse
	for i := 0; i < len(output); i++ {

		userDetail := user.UserDetail{
			UserEmail:    output[i].UserEmail,
			UserPhoneNum: output[i].UserPhoneNum,
		}

		userName := user.UsernameDto{
			Def: output[i].Username.Ko,
			Ko:  output[i].Username.Ko,
			En:  output[i].Username.En,
			Zh:  output[i].Username.Zh,
			Jp:  output[i].Username.Jp,
			Ru:  output[i].Username.Ru,
			Vi:  output[i].Username.Vi,
		}

		profile := user.UserProfile{
			ProfileUrl: output[i].ProfileUrl,
			ProfileMsg: output[i].ProfileMsg,
		}

		var deptInfo []user.DeptInfoDto

		for _, temp := range output[i].DeptInfo {
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

		temp := user.GetUserInfoResponse{
			UserHash:   output[i].UserHash,
			UserDetail: userDetail,
			Username:   userName,
			Profile:    profile,
			OrgCode:    output[i].OrgCodes,
			DeptInfo:   deptInfo,
		}
		res = append(res, temp)
	}

	response.SendSuccess(c, res)

}
