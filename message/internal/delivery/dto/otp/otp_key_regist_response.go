package otp

type OtpKeyRegistResponse struct {
	OtpRegDate       string `json:"otpRegDate"`
	SvChatKeyVersion string `json:"svChatKeyVersion"`
	SvNoteKeyVersion string `json:"svNoteKeyVersion"`
}
