package appValidation

type AppValidationRequestHeader struct {
	Hash   string // 배포 앱 해시
	Device string // device 종류 A, I, W
	Uuid   string // UUID
}

type AppValidationRequestBody struct {
	WorksCode string `json:"worksCode"`
}

type AppValidationRequestDTO struct {
	Header AppValidationRequestHeader `json:"header"`
	Body   AppValidationRequestBody   `json:"body"`
}
