package requestDTO

type DeleteDeptRequestDTO struct {
	Body   DeleteDeptRequestBody
	Header DeleteDeptRequestHeader
}

type DeleteDeptRequestBody struct {
	DeptOrg  string `json:"deptOrg"  validate:"required"`
	DeptCode string `json:"deptCode"  validate:"required"`
}

type DeleteDeptRequestHeader struct {
}
