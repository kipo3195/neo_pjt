package entity

type UserRegisterInfoEntity struct {
	Salt string `json:"salt"`
	Hash string `json:"hash"`
}
