package device

type DeviceRegistRequest struct {
	Id        string `json:"id"`
	Uuid      string `json:"uuid"`
	ModelName string `json:"modelName"`
	Version   string `json:"version"`
	Challenge string `json:"challenge"`
}

type DeviceRegistResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
