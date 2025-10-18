package input

type RefreshTokenCheckInput struct {
	UserId       string
	Uuid         string
	RefreshToken string
	WithoutId    bool
}
