package usecase

import (
	"context"
	"log"
	"message/internal/domain/lineKey/repository"
	"message/pkg/util"
	"time"

	"github.com/oklog/ulid/v2"
)

type lineKeyUsecase struct {
	repository repository.LineKeyRepository
	ulidGen    *util.ULIDGenerator
}

type LineKeyUsecase interface {
	GetLineKey(ctx context.Context) string
}

func NewLineKeyUsecase(repository repository.LineKeyRepository, ulidGen *util.ULIDGenerator) LineKeyUsecase {
	return &lineKeyUsecase{
		repository: repository,
		ulidGen:    ulidGen,
	}
}

func (u *lineKeyUsecase) GetLineKey(ctx context.Context) string {

	lineKey := u.ulidGen.New()
	uid, _ := ulid.Parse(lineKey)

	log.Println("발급된 라인키 : ", lineKey)
	log.Printf(" 발급된 시간: %v\n", time.UnixMilli(int64(uid.Time())))

	return lineKey
}
