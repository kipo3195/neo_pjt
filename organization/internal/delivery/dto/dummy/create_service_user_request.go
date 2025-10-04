package dummy

type CreateServiceUserRequest struct {
	UserCount int    `json:"userCount" validate:"required"`
	Keyword   string `json:"keyword" validate:"required"`
}
