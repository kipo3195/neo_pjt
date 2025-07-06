package dto

// server to server 공통 response.
type ServerResponse[T any] struct {
	Result string `json:"result"`
	Data   T      `json:"data"`
}
