package input

type EventUsersInput struct {
	Org      string `json:"org"`
	UserId   string `json:"userId"`
	UserHash string `json:"userHash"`
}
