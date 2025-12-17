package entity

type ServiceUserEntity struct {
	UserId   string
	UserHash string
	Salt     string
	UserAuth string
}
