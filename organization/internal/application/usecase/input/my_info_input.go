package input

type MyInfoInput struct {
	MyHash string
}

func MakeMyInfoInput(myHash string) MyInfoInput {
	return MyInfoInput{
		MyHash: myHash,
	}
}
