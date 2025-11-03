package entity

type GetProfileInfoEntity struct {
	UserHash []string
}

func MakeGetProfileInfoEntity(userHash []string) GetProfileInfoEntity {
	return GetProfileInfoEntity{
		UserHash: userHash,
	}
}
