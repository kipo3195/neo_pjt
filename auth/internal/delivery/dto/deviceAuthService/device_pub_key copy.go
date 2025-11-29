package deviceAuthService

type DeviceOtp struct {
	OtpRegDate       string `json:"otpRegDate"`
	SvChatKeyVersion string `json:"svChatKeyVersion"`
	SvNoteKeyVersion string `json:"svNoteKeyVersion"`
}
