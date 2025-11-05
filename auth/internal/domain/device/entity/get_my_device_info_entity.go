package entity

type GetMyDeviceInfoEntity struct {
	UserHash string
}

func MakeGetMyDeviceInfoEntity(userHash string) GetMyDeviceInfoEntity {
	return GetMyDeviceInfoEntity{
		UserHash: userHash,
	}
}
