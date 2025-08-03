package requestDTO

type DeleteDeptUserRequestDTO struct {
	Body   DeleteDeptUserRequestBody
	Header DeleteDeptUserRequestHeader
}

type DeleteDeptUserRequestBody struct {
	UserHash string `json:"userHash" validate:"required"`
	DeptCode string `json:"deptCode" validate:"required"`
	DeptOrg  string `json:"deptOrg" validate:"required"`
}

type DeleteDeptUserRequestHeader struct {
}
