package entity

import sharedEntity "org/internal/domain/shared/entity"

type OrgInfo struct {
	DeptCode       string                  `json:"deptCode"`
	ParentDeptCode string                  `json:"parentDeptCode"`
	Name           sharedEntity.NameEntity `json:"name"`
	//UpdateHash     string     `json:"updateHash"`
	Kind   string `json:"kind"`
	Id     string `json:"id"`
	Header string `json:"header,omitempty"`
}
