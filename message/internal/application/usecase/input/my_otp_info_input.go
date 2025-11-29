package input

type MyOtpInfoInput struct {
	UserId      string
	VersionType string // 특정 버전 specific, 최신 버전 latest, 전체 버전 all
	VersionInfo []string
	Uuid        string
}
