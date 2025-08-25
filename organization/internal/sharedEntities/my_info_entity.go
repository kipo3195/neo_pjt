package sharedEntities

type MyInfoEntity struct {
	UserHash     string       `json:"userHash"`
	UserPhoneNum string       `json:"userPhoneNum"`
	Username     NameEntity   `json:"userName"`
	ProfileUrl   string       `json:"profileUrl"`
	ProfileMsg   string       `json:"profileMsg"`
	DeptInfo     []DeptEntity `json:"deptInfo"`
}
