package deviceAuthService

type DeviceRegistResponse struct {
	AccessToken     string `json:"accessToken"`
	RefreshToken    string `json:"refreshToken"`
	RefreshTokenExp string `json:"refreshTokenExp"`
	// DeviceChallenge string `json:"deviceChallenge"`
	OtpRegDate   string `json:"otpRegDate"`
	SvKeyVersion string `json:"svKeyVersion"`
}
