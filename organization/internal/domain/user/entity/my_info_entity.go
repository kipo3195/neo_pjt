package entity

type MyInfoEntity struct {
	UserHash     string           `json:"userHash"`
	UserPhoneNum string           `json:"userPhoneNum"`
	MyOrg        string           `json:"myOrg"`
	Username     UserNameEntity   `json:"userName"`
	UserEmail    string           `json:"userEmail"`
	DeptInfo     []DeptInfoEntity `json:"deptInfo"`
}
