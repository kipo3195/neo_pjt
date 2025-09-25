package entity

type MyInfoEntity struct {
	UserHash     string           `json:"userHash"`
	UserPhoneNum string           `json:"userPhoneNum"`
	Username     UserNameEntity   `json:"userName"`
	ProfileUrl   string           `json:"profileUrl"`
	ProfileMsg   string           `json:"profileMsg"`
	DeptInfo     []DeptInfoEntity `json:"deptInfo"`
}
