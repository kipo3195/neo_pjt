package appValidation

type WorksConfigDTO struct {
	TimeZone   string            `json:"timeZone"`
	Language   string            `json:"language"`
	SkinHash   string            `json:"skinHash"`
	ConfigHash string            `json:"configHash"`
	Skin       []SkinFileInfoDTO `json:"skin"`
}
