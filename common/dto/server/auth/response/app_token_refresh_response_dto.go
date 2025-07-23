package auth

type AppTokenRefreshResponseDTO struct {
	Header AppTokenRefreshResponseHeader
	Body   AppTokenRefreshResponseBody
}

type AppTokenRefreshResponseHeader struct {
}

type AppTokenRefreshResponseBody struct {
	AppToken string `json:"AppToken"`
}
