package profile

type RegistProfileMsgRequest struct {
	ProfileMsg string `json:"profileMsg" validate:"required"`
}
