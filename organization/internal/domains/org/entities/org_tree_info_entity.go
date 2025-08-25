package entities

import (
	"org/internal/sharedEntities"
)

type OrgTreeInfo struct {
	DeptCode       string                    `json:"deptCode"`
	ParentDeptCode string                    `json:"parentDeptCode"`
	Name           sharedEntities.NameEntity `json:"name"`
	//UpdateHash     string         `json:"updateHash"`
	SubDept []OrgTreeInfo `json:"subDept,omitempty"` // 재귀 구조
	Id      string        `json:"id"`
	Kind    string        `json:"kind"`             // 사용자, 부서 구분
	Header  string        `json:"header,omitempty"` // omitempty -> 빈값인 경우 생략됨. zero value
}
