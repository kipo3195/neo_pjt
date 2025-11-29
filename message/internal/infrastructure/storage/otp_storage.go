package storage

import (
	"context"
	"log"
	"message/internal/consts"
	"message/internal/domain/otp/entity"
	"sync"
)

type otpStorage struct {
	mu              sync.RWMutex
	chKeyMap        map[string]entity.OtpKeyInfoEntity //
	noKeyMap        map[string]entity.OtpKeyInfoEntity //
	svKeyMap        map[string]string                  // 대칭키 저장 맵
	svKeyVersionMap map[string]string                  // 대칭키 저장 맵
}

type OtpStorage interface {
	SaveOtpKeyStorage(ctx context.Context, id string, uuid string, entity entity.OtpKeyInfoEntity) error
	GetMyOtpInfoStorage(ctx context.Context, en entity.MyOtpInfoEntity, version string, t string) (entity entity.OtpKeyInfoEntity, err error)
	PutServerKey(svKey, t string, version string) error
	GetServerKey(t string) (string, string)
}

func NewOtpStorage() OtpStorage {
	return &otpStorage{
		chKeyMap:        make(map[string]entity.OtpKeyInfoEntity),
		noKeyMap:        make(map[string]entity.OtpKeyInfoEntity),
		svKeyMap:        make(map[string]string),
		svKeyVersionMap: make(map[string]string),
	}
}

func (s *otpStorage) SaveOtpKeyStorage(ctx context.Context, id string, uuid string, entity entity.OtpKeyInfoEntity) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if entity.Kind == consts.CHAT {
		s.chKeyMap[id+":"+uuid+":"+entity.SvKeyVersion] = entity
	} else if entity.Kind == consts.NOTE {
		s.noKeyMap[id+":"+uuid+":"+entity.SvKeyVersion] = entity
	}
	log.Printf("[SaveOtpKeyStorage] id :%s, uuid :%s, kind:%s save success. \n", id, uuid, entity.Kind)
	return nil
}

func (s *otpStorage) GetMyOtpInfoStorage(ctx context.Context, en entity.MyOtpInfoEntity, version string, t string) (entity.OtpKeyInfoEntity, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := en.UserId + ":" + en.Uuid + ":" + version

	switch t {
	case consts.CHAT:
		value, ok := s.chKeyMap[key]
		if !ok {
			return entity.OtpKeyInfoEntity{}, consts.ErrOtpNotFound
		}
		return value, nil

	case consts.NOTE:
		value, ok := s.noKeyMap[key]
		if !ok {
			return entity.OtpKeyInfoEntity{}, consts.ErrOtpNotFound
		}
		return value, nil

	default:
		return entity.OtpKeyInfoEntity{}, consts.ErrInvalidOtpType
	}
}

func (s *otpStorage) PutServerKey(svKey string, t string, version string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.svKeyMap[t] = svKey
	s.svKeyVersionMap[t] = version
	return nil
}

func (s *otpStorage) GetServerKey(t string) (string, string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	serverKey := s.svKeyMap[t]
	serverKeyVersion := s.svKeyVersionMap[t]
	if serverKey == "" || serverKeyVersion == "" {
		return "neo_encrypt_key", "" // 키, 버전이 없을때
	}
	return serverKey, serverKeyVersion
}
