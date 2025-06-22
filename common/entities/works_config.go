package entities

type WorksConfig struct {
	TimeZone      string `json:"timeZone"`
	Language      string `json:"language"`
	SkinVersion   string `json:"skinVersion"`
	ConfigVersion string `json:"configVersion"`
}
