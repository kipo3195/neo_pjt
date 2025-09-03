package entity

type ValidationEntity struct {
	Hash      string
	Device    string
	Uuid      string
	WorksCode string
}

func NewAppValidationEntity(hash string, device string, uuid string, worksCode string) ValidationEntity {
	return ValidationEntity{
		Hash:      hash,
		Device:    device,
		Uuid:      uuid,
		WorksCode: worksCode,
	}
}
