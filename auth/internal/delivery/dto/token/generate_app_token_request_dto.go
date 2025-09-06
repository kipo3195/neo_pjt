package token

type GenerateAppTokenRequestDTO struct {
	Header GenerateAppTokenRequestHeader
	Body   GenerateAppTokenRequestBody
}

type GenerateAppTokenRequestBody struct {
	Uuid string `json:"uuid"`
}

type GenerateAppTokenRequestHeader struct {
	Token string
}
