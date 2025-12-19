package handler

import (
	"encoding/json"
	"log"
	"user/internal/application/orchestrator"
	"user/internal/delivery/adapter"
	"user/internal/delivery/dto/userInfoService"
	"user/internal/delivery/util"
	commonConsts "user/pkg/consts"
	"user/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type UserInfoServiceHandler struct {
	svc *orchestrator.UserInfoService
}

func NewUserInfoServiceHandler(svc *orchestrator.UserInfoService) *UserInfoServiceHandler {
	return &UserInfoServiceHandler{
		svc: svc,
	}
}

func (h *UserInfoServiceHandler) GetMyDetailInfo(c *gin.Context) {

	ctx := c.Request.Context()

	// AT에 있는 정보로 요청
	myHash := util.GetUserHashByAccessToken(c)

	req := userInfoService.GetMyInfoRequest{
		MyHash: myHash,
	}

	log.Println(req)

	// 어차피 adapter layer는 원시 값을 usecase에 맞는 input으로 변환하는 처리를 하므로
	myInfoInput := adapter.MakeGetMyInfoInput(req.MyHash)
	myInfoOutput, err := h.svc.GetMyInfo(ctx, myInfoInput)

	if err != nil {
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		return
	}

	detail := make([]userInfoService.UserDetail, 0)
	for _, d := range myInfoOutput.UserDetail {

		temp := userInfoService.UserDetail{
			UserHash:     d.UserHash,
			UserEmail:    d.UserEmail,
			UserPhoneNum: d.UserPhoneNum,
		}

		detail = append(detail, temp)
	}

	profile := make([]userInfoService.UserProfile, 0)
	for _, d := range myInfoOutput.UserProfile {

		temp := userInfoService.UserProfile{
			UserHash:     d.UserHash,
			ProlfileHash: d.ProfileHash,
			ProfileMsg:   d.ProfileMsg,
		}

		profile = append(profile, temp)
	}

	userInfo := userInfoService.UserInfo{
		UserDetail:  detail,
		UserProfile: profile,
	}

	res := userInfoService.GetMyInfoServiceResponse{
		UserInfo: userInfo,
	}

	// 리스트 응답을 감싸는 구조체를 유지하는 게 더 낫습니다.
	// 즉, 지금처럼 GetUserDetailInfoResponse 안에 UserDetails []UserDetail 필드를 두는 구조가 더 바람직합니다.
	response.SendSuccess(c, res)

}

func (h *UserInfoServiceHandler) GetUserInfo(c *gin.Context) {

	ctx := c.Request.Context()

	var req userInfoService.GetUserInfoServiceRequest

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

	detailInput, profileInput := adapter.MakeGetUserInfoInput(req.ReqUsers)
	userInfoOutput, err := h.svc.GetUserInfo(ctx, detailInput, profileInput)

	if err != nil {
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		return
	}

	detail := make([]userInfoService.UserDetail, 0)
	for _, d := range userInfoOutput.UserDetail {

		temp := userInfoService.UserDetail{
			UserHash:     d.UserHash,
			UserEmail:    d.UserEmail,
			UserPhoneNum: d.UserPhoneNum,
		}

		detail = append(detail, temp)
	}

	profile := make([]userInfoService.UserProfile, 0)
	for _, d := range userInfoOutput.UserProfile {

		temp := userInfoService.UserProfile{
			UserHash:     d.UserHash,
			ProlfileHash: d.ProfileHash,
			ProfileMsg:   d.ProfileMsg,
		}

		profile = append(profile, temp)
	}

	userInfo := userInfoService.UserInfo{
		UserDetail:  detail,
		UserProfile: profile,
	}

	res := userInfoService.GetUserInfoServiceResponse{
		UserInfo: userInfo,
	}
	// 리스트 응답을 감싸는 구조체를 유지하는 게 더 낫습니다.
	// 즉, 지금처럼 GetUserDetailInfoResponse 안에 UserDetails []UserDetail 필드를 두는 구조가 더 바람직합니다.
	response.SendSuccess(c, res)

}
