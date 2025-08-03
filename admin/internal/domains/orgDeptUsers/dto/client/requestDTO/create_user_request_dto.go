package requestDTO

type CreateDeptUserRequestDTO struct {
	Body   CreateDeptUserRequestBody
	Header CreateDeptUserRequestHeader
}

type CreateDeptUserRequestBody struct {
	UserHash             string `json:"userHash" validate:"required"`
	DeptCode             string `json:"deptCode" validate:"required"`
	DeptOrg              string `json:"deptOrg" validate:"required"`
	PositionCode         string `json:"positionCode"`         // 직위 코드
	RoleCode             string `json:"roleCode"`             // 직책 코드
	IsConcurrentPosition string `json:"isConcurrentPosition"` // 겸직 여부
}

type CreateDeptUserRequestHeader struct {
}
