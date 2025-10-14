package entity

type GenerateAuthtokenEntity struct {
	Id   string
	Uuid string
}

func MakeGenerateAuthTokenEntity(id string, uuid string) GenerateAuthtokenEntity {
	return GenerateAuthtokenEntity{
		Id:   id,
		Uuid: uuid,
	}
}
