package storage

type otpStorage struct {
	chKeyMap map[string]string //
	noKeyMap map[string]string //
}

type OtpStorage interface {
}

func NewOtpStorage() OtpStorage {
	return &otpStorage{
		chKeyMap: make(map[string]string),
		noKeyMap: make(map[string]string),
	}
}
