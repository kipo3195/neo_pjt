package input

type RegistServiceUserInput struct {
	Org          string
	UserId       []string
	UserAuth     string
	UserIdPrefix string
	Start        int
	End          int
}
