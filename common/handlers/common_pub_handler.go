package handlers

import (
	clDto "common/dto/client"
	dto "common/dto/common"
	"common/usecases"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	consts "common/consts"

	"github.com/go-playground/validator"
)

type CommonPubHandler struct {
	usecase usecases.CommonPubUsecase
}

func NewCommonPubHandler(uc usecases.CommonPubUsecase) *CommonPubHandler {
	return &CommonPubHandler{usecase: uc}
}

func (h *CommonPubHandler) AppValidation(w http.ResponseWriter, r *http.Request) {

	// context 생성
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// response
	var res dto.Response

	fmt.Println("1")
	// request body 데이터 -> dto로 변경
	var req = &clDto.AppValidationRequest{
		Uuid:       r.URL.Query().Get("uuid"),
		AppToken:   r.URL.Query().Get("appToken"),
		Device:     r.URL.Query().Get("device"),
		SkinHash:   r.URL.Query().Get("skinHash"),
		ConfigHash: r.URL.Query().Get("configHash"),
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		res.Result = consts.ERROR
		res.Data = dto.ErrorResponse{
			Code:    consts.E_103,
			Message: consts.E_103_MSG,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}
	fmt.Println("2")
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
	fmt.Println("3")
	// 검증
	//data, err := h.usecase.AppValidation(ctx, req)

	// if err != nil || !data { // 에러
	// 	switch {
	// 	case errors.Is(err, consts.ErrDbRowNotFound):
	// 		// 매핑된 hash 정보가 없음
	// 		res.Result = consts.FAIL
	// 		res.Data = newErrorResp(consts.AUTH_F001, consts.AUTH_F001_MSG)
	// 		w.WriteHeader(http.StatusBadRequest)
	// 	case errors.Is(err, consts.ErrTokenExpired):
	// 		res.Result = consts.ERROR
	// 		res.Data = newErrorResp(consts.E_107, consts.E_107_MSG)
	// 		w.WriteHeader(http.StatusBadRequest)
	// 	case errors.Is(err, consts.ErrTokenSignatureInvalid):
	// 		res.Result = consts.FAIL
	// 		res.Data = newErrorResp(consts.AUTH_F005, consts.AUTH_F005_MSG)
	// 		w.WriteHeader(http.StatusBadRequest)
	// 	case errors.Is(err, consts.ErrDB):
	// 		res.Result = consts.ERROR
	// 		res.Data = newErrorResp(consts.E_102, consts.E_102_MSG)
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 	case errors.Is(err, consts.ErrTokenParsing):
	// 		res.Result = consts.ERROR
	// 		res.Data = newErrorResp(consts.E_105, consts.E_105_MSG)
	// 		w.WriteHeader(http.StatusBadRequest)
	// 	case errors.Is(err, consts.ErrInvalidClaims):
	// 		res.Result = consts.FAIL
	// 		res.Data = newErrorResp(consts.AUTH_F002, consts.AUTH_F002_MSG)
	// 		w.WriteHeader(http.StatusBadRequest)
	// 	default:
	// 		res.Result = consts.ERROR
	// 		res.Data = newErrorResp(consts.E_500, consts.E_500_MSG)
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 	}
	// } else {
	// 	res.Result = consts.SUCCESS
	// 	// 데이터가 있어야할까?
	// }
	json.NewEncoder(w).Encode(res)
}
