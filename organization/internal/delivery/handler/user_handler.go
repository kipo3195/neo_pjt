package handler

import (
	"org/internal/application/usecase"
)

type UserHandler struct {
	usecase usecase.UserUsecase
}

func NewUserHandler(usecase usecase.UserUsecase) *UserHandler {

	return &UserHandler{
		usecase: usecase,
	}
}

// func (h *UserHandler) GetMyInfo(c *gin.Context) {

// 	// context 생성
// 	ctx := c.Request.Context()

// 	// 인증 토큰에서 요청 사용자의 hash 정보 추출
// 	myHash := util.GetUserHashByAccessToken(c)
// 	if myHash == "" {
// 		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_110, commonConsts.E_110_MSG)
// 		return
// 	}

// 	log.Println("[GetMyInfo] myHash : ", myHash)

// 	// dto 생성
// 	myInfoInput := adapter.MakeMyInfoInput(myHash)
// 	output, err := h.usecase.GetMyInfo(ctx, myInfoInput)

// 	userName := user.UserNameDto{
// 		Def: output.Username.Ko, // 수정 필요
// 		Ko:  output.Username.Ko,
// 		En:  output.Username.En,
// 		Jp:  output.Username.Jp,
// 		Zh:  output.Username.Zh,
// 		Ru:  output.Username.Ru,
// 		Vi:  output.Username.Vi,
// 	}

// 	var deptInfo []user.DeptInfoDto

// 	for _, temp := range output.DeptInfo {

// 		positionName := user.PositionNameDto{
// 			Ko: temp.PositionName.KoLang,
// 			En: temp.PositionName.EnLang,
// 			Zh: temp.PositionName.ZhLang,
// 			Jp: temp.PositionName.JpLang,
// 		}
// 		roleName := user.RoleNameDto{
// 			Ko: temp.RoleName.KoLang,
// 			En: temp.RoleName.EnLang,
// 			Zh: temp.RoleName.ZhLang,
// 			Jp: temp.RoleName.JpLang,
// 		}

// 		deptName := user.DeptNameDto{
// 			Def: temp.DeptName.DefLang,
// 			Ko:  temp.DeptName.KoLang,
// 			En:  temp.DeptName.EnLang,
// 			Jp:  temp.DeptName.JpLang,
// 			Zh:  temp.DeptName.ZhLang,
// 			Vi:  temp.DeptName.ViLang,
// 			Ru:  temp.DeptName.RuLang,
// 		}

// 		deptInfo = append(deptInfo, user.DeptInfoDto{
// 			DeptOrg:      temp.DeptOrg,
// 			DeptCode:     temp.DeptCode,
// 			DeptName:     deptName,
// 			Header:       temp.Header,
// 			Description:  temp.Description,
// 			PositionName: positionName,
// 			RoleName:     roleName,
// 		})
// 	}

// 	myInfo := user.DetailInfo{
// 		UserHash: output.UserHash,
// 		UserName: userName,
// 		OrgCode:  output.OrgCodes, // 어느 부서에도 속하지 않았다면 org code는 알 수 없는 구조
// 		DeptInfo: deptInfo,
// 	}

// 	res := user.GetMyInfoResponse{
// 		MyInfo: myInfo,
// 	}

// 	// response.
// 	if err == nil {
// 		// http status code 200
// 		response.SendSuccess(c, res)
// 	} else {
// 		// http status code 400
// 		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
// 	}

// }

// func (h *UserHandler) GetUserInfo(c *gin.Context) {

// 	//context 생성
// 	ctx := c.Request.Context()

// 	// request 생성
// 	var req user.GetUserInfoRequest

// 	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
// 		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
// 		return
// 	}

// 	// 필수 데이터 검증
// 	validate := validator.New()
// 	if err := validate.Struct(req); err != nil {
// 		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_108, commonConsts.E_108_MSG)
// 		return
// 	}

// 	input := adapter.MakeGetUserInfoInput(req.UserHashs)
// 	output, err := h.usecase.GetUserInfo(ctx, input)

// 	if err != nil {
// 		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
// 		return
// 	}

// 	var res user.GetUserInfoResponse
// 	res.DetailInfos = make([]user.DetailInfo, 0)
// 	for i := 0; i < len(output); i++ {

// 		userName := user.UserNameDto{
// 			Def: output[i].Username.Ko,
// 			Ko:  output[i].Username.Ko,
// 			En:  output[i].Username.En,
// 			Zh:  output[i].Username.Zh,
// 			Jp:  output[i].Username.Jp,
// 			Ru:  output[i].Username.Ru,
// 			Vi:  output[i].Username.Vi,
// 		}

// 		var deptInfo []user.DeptInfoDto

// 		for _, temp := range output[i].DeptInfo {

// 			positionName := user.PositionNameDto{
// 				Ko: temp.PositionName.KoLang,
// 				En: temp.PositionName.EnLang,
// 				Zh: temp.PositionName.ZhLang,
// 				Jp: temp.PositionName.JpLang,
// 			}
// 			roleName := user.RoleNameDto{
// 				Ko: temp.RoleName.KoLang,
// 				En: temp.RoleName.EnLang,
// 				Zh: temp.RoleName.ZhLang,
// 				Jp: temp.RoleName.JpLang,
// 			}

// 			deptName := user.DeptNameDto{
// 				Def: temp.DeptName.DefLang,
// 				Ko:  temp.DeptName.KoLang,
// 				En:  temp.DeptName.EnLang,
// 				Jp:  temp.DeptName.JpLang,
// 				Zh:  temp.DeptName.ZhLang,
// 				Vi:  temp.DeptName.ViLang,
// 				Ru:  temp.DeptName.RuLang,
// 			}

// 			deptInfo = append(deptInfo, user.DeptInfoDto{
// 				DeptOrg:      temp.DeptOrg,
// 				DeptCode:     temp.DeptCode,
// 				DeptName:     deptName,
// 				Header:       temp.Header,
// 				Description:  temp.Description,
// 				PositionName: positionName,
// 				RoleName:     roleName,
// 			})
// 		}

// 		temp := user.DetailInfo{
// 			UserHash: output[i].UserHash,
// 			UserName: userName,
// 			OrgCode:  output[i].OrgCodes,
// 			DeptInfo: deptInfo,
// 		}
// 		res.DetailInfos = append(res.DetailInfos, temp)
// 	}

// 	response.SendSuccess(c, res)

// }
