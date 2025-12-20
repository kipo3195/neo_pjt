package userInfoService

type UserProfile struct {
	UserHash       string `json:"userHash"`
	ProfileVersion int64  `json:"profileVersion"`
	ProfileMsg     string `json:"profileMsg"`
}
