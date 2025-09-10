package appToken

type AppTokenRefreshRequestDTO struct {
	Header AppTokenRefreshRequestHeader
	Body   AppTokenRefreshRequestBody
}

type AppTokenRefreshRequestHeader struct {
}

type AppTokenRefreshRequestBody struct {
	Uuid         string `json:"uuid" validate:"required"`
	RefreshToken string `json:"refreshToken" validate:"required"`
}
