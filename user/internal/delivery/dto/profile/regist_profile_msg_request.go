package profile

type RegistProfileMsgRequest struct {
	Msg string `json:"msg" validate:"required"`
}
