package entity

import sharedEntity "org/internal/domain/shared/entity"

type OrgInfo struct {
	DeptCode       string                  `json:"deptCode"`
	ParentDeptCode string                  `json:"parentDeptCode"`
	Name           sharedEntity.NameEntity `json:"name"`
	//UpdateHash     string     `json:"updateHash"`
	Kind        string `json:"kind"`
	UserHash    string `json:"userHash,omitempty"`
	UserId      string `json:"userId,omitempty"`
	Header      string `json:"header,omitempty"`
	Description string `json:"description,omitempty"`
}
