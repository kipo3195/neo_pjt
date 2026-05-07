package login

type LoginRequest struct {
	UserIdPrefix    string `json:"userIdPrefix"`
	ConnectionCount int    `json:"connectionCount"`
}
