package entity

type MyOtpInfoEntity struct {
	UserId      string
	VersionType string // 특정 버전 specific, 최신 버전 latest, 전체 버전 all
	VersionInfo []string
	Uuid        string
}

func MakeMyOtpInfoEntity(userId string, versionType string, versionInfo []string, uuid string) MyOtpInfoEntity {
	return MyOtpInfoEntity{
		UserId:      userId,
		VersionType: versionType,
		VersionInfo: versionInfo,
		Uuid:        uuid,
	}
}
