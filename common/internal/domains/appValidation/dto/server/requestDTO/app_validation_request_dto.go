package requestDTO

type AppValidationRequestDTO struct {
	Body   AppValidationRequestBody
	Header AppValidationRequestHeader
}

type AppValidationRequestBody struct {
	AppToken   string `json:"appToken" validate:"required"`
	Uuid       string `json:"uuid" validate:"required"`
	Device     string `json:"device" validate:"required"`
	SkinHash   string `json:"skinHash" validate:"required"`
	ConfigHash string `json:"configHash" validate:"required"`
}

type AppValidationRequestHeader struct {
}
