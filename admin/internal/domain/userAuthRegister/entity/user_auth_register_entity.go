package entity

type UserAuthRegisterEntity struct {
	UserHash string
	UserId   string
	UserAuth string
	Salt     string
}
