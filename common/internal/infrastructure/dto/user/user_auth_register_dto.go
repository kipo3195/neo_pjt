package user

type UserAuthRegisterDto struct {
	Id       string
	Salt     string
	UserHash string
	AuthHash string
}
