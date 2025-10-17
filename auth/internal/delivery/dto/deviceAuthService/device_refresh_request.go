package deviceAuthService

type DeviceRefreshRequest struct {
	Id           string `json:"id"`
	AppToken     string `json:"appToken"`
	Uuid         string `json:"uuid"`
	RefreshToken string `json:"refreshToken"`
	ModelName    string `json:"modelName"`
	Version      string `json:"version"`
}
