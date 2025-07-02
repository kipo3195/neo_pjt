package handlers

import (
	"admin/consts"
	commonDto "admin/dto/client/common"
	dto "admin/dto/common"
	"admin/usecases"
	"context"
	"encoding/json"
	"fmt"
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

	var res = commonDto.CreateSkinImgResponse{}

	// 파일 데이터 추출
	file, fileInfo, err := r.FormFile("file")

	if err == nil {
		res.Result = consts.FAIL
		res.Data = dto.ErrorResponse{
			Code:    consts.E_103,
			Message: consts.E_103_MSG,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	var req = commonDto.CreateSkinImgRequest{
		File:     file,
		FileInfo: fileInfo,
	}

	data, err := h.usecase.CreateSkinImg(ctx, req)

	fmt.Println(data)

	if err == nil {
		// http status code 200
		res.Result = consts.SUCCESS
		res.Data = data
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		res.Result = consts.ERROR
		res.Data = dto.ErrorResponse{
			Code:    consts.E_500,
			Message: consts.E_500_MSG,
		}
	}

}
