package storage

import (
	"context"
	"message/internal/domain/otp/entity"
)

type otpStorage struct {
	chKeyMap map[string]string //
	noKeyMap map[string]string //
}

type OtpStorage interface {
	SaveOtpKeyStorage(ctx context.Context, entity *entity.OTPKeyRegistEntity) error
}

func NewOtpStorage() OtpStorage {
	return &otpStorage{
		chKeyMap: make(map[string]string),
		noKeyMap: make(map[string]string),
	}
}

func (s *otpStorage) SaveOtpKeyStorage(ctx context.Context, entity *entity.OTPKeyRegistEntity) error {
	s.chKeyMap[entity.Id+":"+entity.Uuid] = entity.ChatOtpKey
	s.noKeyMap[entity.Id+":"+entity.Uuid] = entity.NoteOtpKey
	return nil
}
