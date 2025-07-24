package handlers

import (
	consts "common/consts"
	clDto "common/dto/client"
	dto "common/dto/common"
	cl
	"common/entities"
	"common/usecases"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type CommonHandler struct {
	usecase usecases.CommonUsecase
}

func NewCommonHandler(uc usecases.CommonUsecase) *CommonHandler {
	return &CommonHandler{usecase: uc}
}

func (h *CommonHandler) GetConfigHash(c *gin.Context) {

	// context 생성
	ctx := c.Request.Context()

	// 데이터 -> dto
	var req = clDto.GetConfigHash{
		SkinHash:   c.Query("skinHash"),
		ConfigHash: c.Query("configHash"),
		Device:     c.Query("device"),
	}

	// 유효성 검증
	if req.SkinHash == "" || req.ConfigHash == "" || req.Device == "" {
		res.Result = consts.FAIL
		res.Data = dto.ErrorResponse{
			Code:    consts.E_103,
			Message: consts.E_103_MSG,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	// usecase 호출
	data := h.usecase.GetConfigHash(toConfigHashEntity(req), ctx)

	res.Result = consts.SUCCESS
	res.Data = data

	// response.
	json.NewEncoder(w).Encode(res)

}

func toConfigHashEntity(dto clDto.GetConfigHash) entities.ConfigHashEntity {
	return entities.ConfigHashEntity{
		ConfigHash: dto.ConfigHash,
		SkinHash:   dto.SkinHash,
		Device:     dto.Device,
	}
}

func (h *CommonHandler) GetSkinImage(w http.ResponseWriter, r *http.Request) {
	// context 생성
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// response dto 생성
	var res = dto.Response{}

	// 데이터 -> dto
	var req = clDto.GetSkinImgRequest{
		SkinHash: r.URL.Query().Get("skinHash"),
		SkinType: r.URL.Query().Get("skinType"),
	}

	log.Println("2")
	// 유효성 검증 로직
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		// 검증 실패 처리
		w.WriteHeader(http.StatusBadRequest)
		res.Result = consts.ERROR
		res.Data = dto.ErrorResponse{
			Code:    consts.E_108,
			Message: consts.E_108_MSG,
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	log.Println("3")
	// 검증
	file, err := h.usecase.GetSkinImg(ctx, req)
	defer file.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	log.Println("4")
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "inline")

	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "failed to send file", http.StatusInternalServerError)
	}
	log.Println("5")

}
