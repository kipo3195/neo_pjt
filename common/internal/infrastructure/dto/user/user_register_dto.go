package user

type UserRegisterRequest struct {
	Id   string `json:"id"`
	Salt string `json:"salt"`
	Hash string `json:"hash"`
}

type UserRegisterResponse struct {
	Result string `json:"result"`
}
