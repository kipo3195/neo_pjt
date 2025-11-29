package adapter

import "message/internal/application/usecase/input"

func MakeMyOtpInfoInput(userId string, uuid string, versionType string, versionInfo []string) input.MyOtpInfoInput {
	return input.MyOtpInfoInput{
		UserId:      userId,
		Uuid:        uuid,
		VersionType: versionType,
		VersionInfo: versionInfo,
	}
}
