package server

import (
	serverUsecase "common/internal/domains/skin/usecases/server"
	"common/pkg/consts"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type SkinHandler struct {
	usecase serverUsecase.SkinUsecase
}

func NewSkinHandler(usecase serverUsecase.SkinUsecase) *SkinHandler {
	return &SkinHandler{
		usecase: usecase,
	}
}

func (h *SkinHandler) PutSkinImg(w http.ResponseWriter, r *http.Request) {

	// context 생성
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// server 토큰 검증은 미들웨어에서

	var res = adminDto.CreateSkinImgResponse{}

	log.Println("111")

	// 파일 데이터 추출
	file, fileInfo, err := r.FormFile("File")
	skinType := r.Header.Get("Skin-Type")

	if err != nil || file == nil || fileInfo.Size == 0 || skinType == "" {
		res.Result = consts.FAIL
		res.Data = dto.ErrorResponse{
			Code:    consts.E_103,
			Message: consts.E_103_MSG,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	var req = adminDto.CreateSkinImgRequest{
		File:     file,
		FileInfo: fileInfo,
		SkinType: skinType,
	}

	data, err := h.usecase.CreateSkinImg(ctx, req)

	log.Println(data)

	if err == nil {
		// http status code 200
		res.Result = consts.SUCCESS
		res.Data = data
	} else {
		// TODO 모든 에러 세분화 할 것.
		w.WriteHeader(http.StatusInternalServerError)
		res.Result = consts.ERROR
		res.Data = dto.ErrorResponse{
			Code:    consts.E_500,
			Message: consts.E_500_MSG,
		}
	}

}
