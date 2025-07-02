package handlers

import (
	"admin/usecases"
	"context"
	"net/http"
	"time"
)

// 관리자 - common service 연계 handler
type CommonHandler struct {
	usecase usecases.CommonUsecase
}

func NewCommonHandler(uc usecases.CommonUsecase) *CommonHandler {
	return &CommonHandler{usecase: uc}
}

func (h *CommonHandler) CreateSkinImg(w http.ResponseWriter, r *http.Request) {

	// context 생성 - admin_route에 정의된 middleware에서 context에 관여함.
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// var res = commonDto.CreateSkinImgResponse{}

	// var req = commonDto.CreateSkinImgRequest{}

}
