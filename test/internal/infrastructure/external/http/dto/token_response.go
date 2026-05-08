package dto

type TokenResponse struct {
	Result string `json:"result"`
	Data   Token  `json:"data"`
}
