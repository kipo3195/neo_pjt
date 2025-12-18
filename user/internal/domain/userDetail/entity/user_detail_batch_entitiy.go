package entity

type UserDetailBatchEntity struct {
	Org          string `json:"org"`
	UserId       string `json:"userId"`
	UserHash     string `json:"userHash"`
	UserEmail    string `json:"userEmail"`
	UserPhoneNum string `json:"userPhoneNum"`
}
