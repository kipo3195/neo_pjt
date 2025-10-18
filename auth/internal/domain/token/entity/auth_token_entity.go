package entity

type AuthTokenEntity struct {
	Id    string `json:"id"`
	Uuid  string `json:"uuid"`
	At    string `json:"at"`
	Rt    string `json:"rt"`
	RtExp string `json:"rtExp"`
}
