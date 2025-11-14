package entity

type GetProfileMsgResultEntity struct {
	UserHash   string `column:"user_hash"`
	ProfileMsg string `column:"profile_msg"`
}
