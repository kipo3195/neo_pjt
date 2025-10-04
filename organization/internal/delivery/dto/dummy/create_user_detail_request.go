package dummy

type CreateUserDetailRequest struct {
	Keyword string `json:"keyword" validate:"required"`
	Type    string `json:"type" validate:"required"`
}
