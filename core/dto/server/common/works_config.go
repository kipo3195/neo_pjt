package common

type WorksConfig struct {
	TimeZone   string         `json:"timeZone"`
	Language   string         `json:"language"`
	SkinHash   string         `json:"skinHash"`
	ConfigHash string         `json:"configHash"`
	Skin       []SkinFileInfo `json:"skin"`
}
