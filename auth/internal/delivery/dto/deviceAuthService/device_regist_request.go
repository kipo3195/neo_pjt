package deviceAuthService

type DeviceRegistRequest struct {
	Id           string         `json:"id" validate:"required"`
	Uuid         string         `json:"uuid" validate:"required"`
	ModelName    string         `json:"modelName" validate:"required"`
	Version      string         `json:"version" validate:"required"`
	Challenge    string         `json:"challenge" validate:"required"`
	DevicePubKey []DevicePubKey `json:"devicePubKey"` // 채팅, 쪽지 내용암호화를 위한 대칭키를 암호화할 디바이스(uuid)별 구분 키
}
