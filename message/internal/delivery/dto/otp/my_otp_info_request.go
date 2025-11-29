package otp

type MyOtpInfoRequest struct {
	VersionType string   `json:"versionType" validate:"required"` // 특정 버전 specific, 최신 버전 latest, 전체 버전 all
	VersionInfo []string `json:"versionInfo"`
	Uuid        string   `json:"uuid" validate:"required"`
}
