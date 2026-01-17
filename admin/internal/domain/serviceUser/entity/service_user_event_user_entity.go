package entity

type ServiceUsersEventUserEntity struct {
	Org      string `json:"org"`
	UserHash string `json:"userHash"`
	UserId   string `json:"userId"`
}
