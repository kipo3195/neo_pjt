package handlers

import (
	"net/http"
	"org/config"
	"org/usecases"
)

type OrgHandler struct {
	usecase usecases.OrgUsecase
	sfg     *config.ServerConfig
}

func NewOrgHandler(sfg *config.ServerConfig, uc usecases.OrgUsecase) *OrgHandler {
	return &OrgHandler{usecase: uc, sfg: sfg}
}
func (h *OrgHandler) Test(w http.ResponseWriter, r *http.Request) {

}
