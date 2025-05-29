package handlers

import (
	"admin/usecases"
	"net/http"
)

type AdminHandler struct {
	usecase usecases.AdminUsecase
}

func NewAdminHandler(r usecases.AdminUsecase) *AdminHandler {
	return &AdminHandler{usecase: r}
}

func (h *AdminHandler) CreateDept(w http.ResponseWriter, r *http.Request) {

}

func (h *AdminHandler) GetDept(w http.ResponseWriter, r *http.Request) {

}

func (h *AdminHandler) UpdateDept(w http.ResponseWriter, r *http.Request) {

}

func (h *AdminHandler) DeleteDept(w http.ResponseWriter, r *http.Request) {

}
