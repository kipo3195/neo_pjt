package otp

type OtpKeyRegistRequest struct {
	Id           string         `json:"id"`
	Uuid         string         `json:"uuid"`
	DevicePubKey []DevicePubKey `json:"devicePubKey"`
}
