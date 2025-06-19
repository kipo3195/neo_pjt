package entities

type MyInfoEntity struct {
	UserHash     string         `json:"userHash"`
	UserPhoneNum string         `json:"userPhoneNum"`
	Username     UsernameEntity `json:"userName"`
}
