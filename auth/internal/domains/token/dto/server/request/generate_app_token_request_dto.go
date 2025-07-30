package common

type GenerateAppTokenRequestDTO struct {
	Header GenerateAppTokenRequestHeader
	Body   GenerateAppTokenRequestBody
}

type GenerateAppTokenRequestBody struct {
	Uuid string
}

type GenerateAppTokenRequestHeader struct {
	Token string
}
