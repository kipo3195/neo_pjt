package otp

type OtpKeyRegistResponse struct {
	OtpRegDate   string `json:"otpRegDate"`
	SvKeyVersion string `json:"svKeyVersion"`
}
