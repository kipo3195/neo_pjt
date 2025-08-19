package requestDTO

type GenerateAppTokenRequestDTO struct {
	Body GenerateAppTokenBody
}

type GenerateAppTokenBody struct {
	Uuid string `json:"uuid"`
}
