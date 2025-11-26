package otp

type OtpKeySvRegistRequest struct {
	Id    string `json:"id"`
	Uuid  string `json:"uuid"`
	Chkey string `json:"chKey"`
	Nokey string `json:"noKey"`
}
