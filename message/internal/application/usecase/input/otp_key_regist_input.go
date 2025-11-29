package input

type OtpKeyRegistInput struct {
	Id           string
	Uuid         string
	DevicePubKey []DevicePubKeyInput
}
