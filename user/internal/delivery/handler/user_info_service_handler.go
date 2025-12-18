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

	temp := make([]string, 0)
	temp = append(temp, myHash)

	req := userInfoService.GetUserInfoServiceRequest{
		UserHashs: temp,
	}

	log.Println(req)

	// 어차피 adapter layer는 원시 값을 usecase에 맞는 input으로 변환하는 처리를 하므로
	detailInput := adapter.MakeGetUserDetailInfoInput(req.UserHashs)
	detailOutput, err := h.svc.UserDetail.GetUserDetailInfo(ctx, detailInput)

	if err != nil {
		// 사용자 정보 조회 에러
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		return
	}

	// 프로필 정보 조회
	profileInput := adapter.MakeGetProfileInfoInput(req.UserHashs)
	profileOutput, err := h.svc.Profile.GetProfileInfo(ctx, profileInput)

	if err != nil {
		// 프로필 이미지 조회 에러
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		return
	}

	// 사용자 상세 정보로 기반한 response 구성
	var res userInfoService.GetMyInfoServiceResponse

	// userHash : 프로필 정보 구조체 {} map 생성 필요

	for i := 0; i < len(detailOutput.UserInfos); i++ {

		// detail 정보 생성
		detail := userInfoService.UserDetail{
			UserHash:     detailOutput.UserInfos[i].UserHash,
			UserEmail:    detailOutput.UserInfos[i].UserEmail,
			UserPhoneNum: detailOutput.UserInfos[i].UserPhoneNum,
		}

		profileInfoEntity := profileOutput.ResultMap[detailOutput.UserInfos[i].UserHash]

		profile := userInfoService.UserProfile{
			ProlfileHash: profileInfoEntity.ProfileImgHash,
			ProfileMsg:   profileInfoEntity.ProfileMsg,
		}

		info := userInfoService.UserInfo{
			UserDetail:  detail,
			UserProfile: profile,
		}

		res.UserInfo = append(res.UserInfo, info)
	}
	// dto 배열을 response할건지 구조체로 배열을 감싼 걸 response 할건지?
	// 20251030
	// 리스트 응답을 감싸는 구조체를 유지하는 게 더 낫습니다.
	// 즉, 지금처럼 GetUserDetailInfoResponse 안에 UserDetails []UserDetail 필드를 두는 구조가 더 바람직합니다.
	response.SendSuccess(c, res)

}

func (h *UserInfoServiceHandler) GetUserDetailInfo(c *gin.Context) {

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

	detailInput := adapter.MakeGetUserDetailInfoInput(req.ReqUsers)
	detailOutput, err := h.svc.UserDetail.GetUserDetailInfo(ctx, detailInput)

	if err != nil {
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		return
	}

	// 프로필 정보 조회
	profileInput := adapter.MakeGetProfileInfoInput(req.ReqUsers)
	profileOutput, err := h.svc.Profile.GetProfileInfo(ctx, profileInput)

	if err != nil {
		// 프로필 이미지 조회 에러
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		return
	}

	// 사용자 상세 정보로 기반한 response 구성
	res := userInfoService.GetUserInfoServiceResponse{
		UserInfo: []userInfoService.UserInfo{},
	}

	for i := 0; i < len(detailOutput.UserInfos); i++ {

		// detail 정보 생성
		detail := userInfoService.UserDetail{
			UserHash:     detailOutput.UserInfos[i].UserHash,
			UserEmail:    detailOutput.UserInfos[i].UserEmail,
			UserPhoneNum: detailOutput.UserInfos[i].UserPhoneNum,
		}

		profileInfoEntity := profileOutput.ResultMap[detailOutput.UserInfos[i].UserHash]

		profile := userInfoService.UserProfile{
			ProlfileHash: profileInfoEntity.ProfileImgHash,
			ProfileMsg:   profileInfoEntity.ProfileMsg,
		}

		info := userInfoService.UserInfo{
			UserDetail:  detail,
			UserProfile: profile,
		}

		res.UserInfo = append(res.UserInfo, info)
	}
	// dto 배열을 response할건지 구조체로 배열을 감싼 걸 response 할건지?
	// 20251030
	// 리스트 응답을 감싸는 구조체를 유지하는 게 더 낫습니다.
	// 즉, 지금처럼 GetUserDetailInfoResponse 안에 UserDetails []UserDetail 필드를 두는 구조가 더 바람직합니다.
	response.SendSuccess(c, res)

}
