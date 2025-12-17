package entity

type UserAuthInfoEntity struct {
	UserId   string
	Salt     string
	UserAuth string
	UserHash string
}
