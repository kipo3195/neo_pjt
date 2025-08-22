package client

import (
	usecases "org/internal/domains/department/usecases/client"

	"github.com/gin-gonic/gin"
)

type DepartmentHandler struct {
	usecase usecases.DepartmentUsecase
}

func NewDepartmentHandler(usecase usecases.DepartmentUsecase) *DepartmentHandler {
	return &DepartmentHandler{
		usecase: usecase,
	}
}

// 눌려서 일부 부서 조회 -> hash가 바뀌었는지 조회 필요.
func (h *DepartmentHandler) GetDept(c *gin.Context) {

	// context 생성

	// response dto 생성

	// request 데이터 파싱 header, body -> dto

	// usecase 호출

	// response.

}
