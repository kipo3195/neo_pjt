package entity

import "org/internal/sharedEntities"

type OrgInfo struct {
	DeptCode       string                    `json:"deptCode"`
	ParentDeptCode string                    `json:"parentDeptCode"`
	Name           sharedEntities.NameEntity `json:"name"`
	//UpdateHash     string     `json:"updateHash"`
	Kind   string `json:"kind"`
	Id     string `json:"id"`
	Header string `json:"header,omitempty"`
}
