package requestDTO

// 인증 요청 body
type AuthRequestBody struct {
	Id       string `json:"id"`
	Password string `json:"password"`
}

type AuthRequestHeader struct {
	Token string `json:"token"`
	Uuid  string `json:"uuid"`
}

type AuthRequestDTO struct {
	Header AuthRequestHeader `json:"header"`
	Body   AuthRequestBody   `json:"body"`
}
