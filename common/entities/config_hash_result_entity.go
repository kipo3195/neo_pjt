package entities

type ConfigHashResultEntity struct {
	ConfigExist bool `json:"configExist"`
	ConfigSame  bool `json:"configSame"`
	SkinExist   bool `json:"skinExist"`
	SkinSame    bool `json:"skinSame"`
}
