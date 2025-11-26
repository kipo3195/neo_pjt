package otp

type OtpKeyRegistRequest struct {
	Id    string `json:"id"`
	Uuid  string `json:"uuid"`
	ChKey string `json:"chKey"`
	NoKey string `json:"noKey"`
}
