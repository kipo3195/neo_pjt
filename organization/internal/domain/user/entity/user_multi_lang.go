package entity

type UserMultiLangEntity struct {
	UserHash string
	Ko       string
	En       string
	Vi       string
	Ru       string
	Zh       string
	Jp       string
}

func MakeUserMultilangEntity(userHash string, ko string, en string, vi string, ru string, jp string, zh string) UserMultiLangEntity {
	return UserMultiLangEntity{
		UserHash: userHash,
		Ko:       ko,
		En:       en,
		Vi:       vi,
		Ru:       ru,
		Zh:       zh,
		Jp:       jp,
	}
}
