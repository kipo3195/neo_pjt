package profile

type GetProfileImgRequest struct {
	UserHash string `json:"userHash" validate:"required"`
}
