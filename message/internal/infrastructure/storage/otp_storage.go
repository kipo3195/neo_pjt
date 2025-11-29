package storage

import (
	"context"
	"log"
	"message/internal/consts"
	"message/internal/domain/otp/entity"
	"sync"
)

type otpStorage struct {
	mu         sync.RWMutex
	chKeyMap   map[string]string //
	noKeyMap   map[string]string //
	regDateMap map[string]string //
}

type OtpStorage interface {
	SaveOtpKeyStorage(ctx context.Context, entity entity.OTPKeyRegistEntity, version string, t string) error
	GetMyOtpInfoStorage(ctx context.Context, entity entity.MyOtpInfoEntity, version string, t string) (value string, err error)
}

func NewOtpStorage() OtpStorage {
	return &otpStorage{
		chKeyMap:   make(map[string]string),
		noKeyMap:   make(map[string]string),
		regDateMap: make(map[string]string),
	}
}

func (s *otpStorage) SaveOtpKeyStorage(ctx context.Context, entity entity.OTPKeyRegistEntity, version string, t string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if t == consts.CHAT {
		s.chKeyMap[entity.Id+":"+version+":"+entity.Uuid] = entity.ChatOtpKey
		log.Println("SaveOtpKeyStorage CHAT :", entity.Id+":"+version+":"+entity.Uuid, entity.ChatOtpKey)
	} else if t == consts.NOTE {
		s.noKeyMap[entity.Id+":"+version+":"+entity.Uuid] = entity.NoteOtpKey
		log.Println("SaveOtpKeyStorage NOTE :", entity.Id+":"+version+":"+entity.Uuid, entity.NoteOtpKey)
	} else if t == consts.DATE {
		s.regDateMap[entity.Id+":"+version+":"+entity.Uuid] = entity.OtpRegDate
		log.Println("SaveOtpKeyStorage DATE :", entity.Id+":"+version+":"+entity.Uuid, entity.OtpRegDate)
	}
	return nil
}

func (s *otpStorage) GetMyOtpInfoStorage(ctx context.Context, entity entity.MyOtpInfoEntity, version string, t string) (value string, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if t == consts.CHAT {
		chKey := s.chKeyMap[entity.UserId+":"+version+":"+entity.Uuid]
		return chKey, nil
	} else if t == consts.NOTE {
		noKey := s.noKeyMap[entity.UserId+":"+version+":"+entity.Uuid]
		return noKey, nil
	} else if t == "date" {
		date := s.regDateMap[entity.UserId+":"+version+":"+entity.Uuid]
		return date, nil
	}
	return "", nil
}
