package entities

type CreateDeptUserEntity struct {
	UserHash string `json:"userHash" validate:"required"`
	DeptCode string `json:"deptCode" validate:"required"`
	DeptOrg  string `json:"deptOrg" validate:"required"`
	// KrLang               string `json:"krLang"`
	// EnLang               string `json:"enLang"`
	// JpLang               string `json:"jpLang"`
	// CnLang               string `json:"cnLang"`
	PositionCode         string `json:"positionCode"`         // 직위 코드
	RoleCode             string `json:"roleCode"`             // 직책 코드
	IsConcurrentPosition string `json:"isConcurrentPosition"` // 겸직 여부
	UpdateHash           string `json:"updateHash"`
}
