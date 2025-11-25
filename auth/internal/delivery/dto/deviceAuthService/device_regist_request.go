package deviceAuthService

type DeviceRegistRequest struct {
	Id        string `json:"id" validate:"required"`
	Uuid      string `json:"uuid" validate:"required"`
	ModelName string `json:"modelName" validate:"required"`
	Version   string `json:"version" validate:"required"`
	Challenge string `json:"challenge" validate:"required"`
	ChKey     string `json:"chKey" validate:"required"` // 채팅 암호화 키 암호화용 공개키
	NoKey     string `json:"noKey" validate:"required"` // 쪽지 암호화 키 암호화용 공개키
}
