package entity

type MyInfoHashEntity struct {
	MyHash string `json:"myHash"`
}

func MakeMyInfoHashEntity(myHash string) MyInfoHashEntity {
	return MyInfoHashEntity{
		MyHash: myHash,
	}
}
