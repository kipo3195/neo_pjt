package otp

type MyOtpInfo struct {
	Version    string `json:"version"`
	KeyType    string `json:"keyType"`
	Key        string `json:"key"`
	OtpRegDate string `json:"otpRegDate"`
}
