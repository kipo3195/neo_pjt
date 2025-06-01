package dto

type AppValidationRequestHeader struct {
	Hash   string // 배포 앱 해시
	Device string // device 종류 A, I, W
	Uuid   string // UUID
}
