package userInfoService

type UserDetail struct {
	UserHash      string `json:"userHash"`
	UserEmail     string `json:"userEmail"`
	UserPhoneNum  string `json:"userPhoneNum"`
	DetailVersion int64  `json:"detailVersion"`
}
