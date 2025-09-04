package output

type AppValidationOutput struct {
	WorksCommonInfo WorksCommonInfo `json:"worksCommonInfo"`
	WorksInfo       WorksInfo       `json:"worksInfo"`
	WorksConfig     WorksConfig     `json:"worksConfig"`
}

type WorksCommonInfo struct {
	ServerUrl string `json:"serverUrl"`
	WorksCode string `json:"worksCode"`
	WorksName string `json:"worksName"`
	UseYn     string `json:"useYn"`
	RegDate   string `json:"regDate"`
}

type WorksInfo struct {
	IssuedAppToken IssuedAppToken `json:"issuedAppToken"`
}

type IssuedAppToken struct {
	AppToken     string `json:"appToken"`
	RefreshToken string `json:"refreshToken"`
}

type WorksConfig struct {
	TimeZone   string `json:"timeZone"`
	Language   string `json:"language"`
	SkinHash   string `json:"skinHash"`
	ConfigHash string `json:"configHash"`
	Skin       []skin `json:"skin"`
}

type skin struct {
	FileHash string `json:"fileHash"`
	SkinType string `json:"worksImg"`
}
