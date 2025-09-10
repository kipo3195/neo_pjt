package appToken

type GenerateAppTokenRequestDTO struct {
	Body GenerateAppTokenBody
}

type GenerateAppTokenBody struct {
	Uuid string `json:"uuid"`
}
