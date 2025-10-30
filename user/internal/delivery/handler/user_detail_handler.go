package handler

import (
	"encoding/json"
	"user/internal/application/usecase"
	"user/internal/delivery/adapter"
	commonConsts "user/pkg/consts"
	"user/pkg/response"

	"user/internal/delivery/dto/userDetail"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type UserDetailHandler struct {
	usecase usecase.UserDetailUsecase
}

func NewUserDetailHandler(usecase usecase.UserDetailUsecase) *UserDetailHandler {
	return &UserDetailHandler{
		usecase: usecase,
	}
}

func (h *UserDetailHandler) GetUserDetailInfo(c *gin.Context) {

	ctx := c.Request.Context()

	var req userDetail.GetUserDetailInfoRequest

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

	input := adapter.MakeGetUserDetailInfoInput(req.UserHashs)
	output, err := h.usecase.GetUserDetailInfo(ctx, input)

	if err != nil {
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		return
	}

	res := make([]userDetail.GetUserDetailInfoResponse, 0) // 빈 배열을 추가해주면 null은 아님.

	for i := 0; i < len(output.UserInfos); i++ {

		temp := userDetail.GetUserDetailInfoResponse{
			UserHash:     output.UserInfos[i].UserHash,
			UserEmail:    output.UserInfos[i].UserEmail,
			UserPhoneNum: output.UserInfos[i].UserPhoneNum,
			ProfileMsg:   "",
		}
		res = append(res, temp)
	}
	// dto 배열을 response할건지 구조체로 배열을 감싼 걸 response 할건지?
	// 20251030
	// 리스트 응답을 감싸는 구조체를 유지하는 게 더 낫습니다.
	// 즉, 지금처럼 GetUserDetailInfoResponse 안에 UserDetails []UserDetail 필드를 두는 구조가 더 바람직합니다.
	response.SendSuccess(c, res)

}
