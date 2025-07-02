package org

type DeleteDeptRequest struct {
	DeptOrg  string `json:"deptOrg"  validate:"required"`
	DeptCode string `json:"deptCode"  validate:"required"`
}
