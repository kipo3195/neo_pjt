package requestDTO

type GenerateAppTokenRequest struct {
	Body GenerateAppTokenBody
}

type GenerateAppTokenBody struct {
	Uuid string `json:"uuid"`
}
