package entities

type WorksConfig struct {
	TimeZone   string `json:"timeZone"`
	Language   string `json:"language"`
	ConfigHash string `json:"configHash"`
}
