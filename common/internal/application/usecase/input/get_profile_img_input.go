package input

type GetProfileImgInput struct {
	UserId string
}

func MakeGetProfileImgInput(userId string) GetProfileImgInput {
	return GetProfileImgInput{
		UserId: userId,
	}
}
