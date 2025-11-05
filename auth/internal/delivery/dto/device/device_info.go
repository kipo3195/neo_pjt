package device

type DeviceInfo struct {
	Uuid      string `json:"uuid"`
	Version   string `json:"version"`
	CreateAt  string `json:"createAt"`
	ModelName string `json:"modelName"`
}
