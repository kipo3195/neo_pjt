package otp

type OtpKeySvRegistRequest struct {
	Id              string            `json:"id"`
	Uuid            string            `json:"uuid"`
	DevicePubKeyDto []DevicePubKeyDto `json:"devicePubKey"`
}
