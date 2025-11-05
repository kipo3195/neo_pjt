package model

type DeviceInfo struct {
	Uuid      string `column:"uuid"`
	ModelName string `column:"model_name"`
	Version   string `column:"version"`
	CreateAt  string `column:"create_at"`
}
