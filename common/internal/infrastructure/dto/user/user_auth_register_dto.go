package user

type UserAuthRegisterDto struct {
	UserId   string
	Salt     string
	UserHash string
	UserAuth string
}
