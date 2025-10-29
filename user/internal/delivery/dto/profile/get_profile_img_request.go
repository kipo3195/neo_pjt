package profile

type GetProfileImgRequest struct {
	UserId string `json:"userId" validate:"required"`
}
