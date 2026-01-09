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
	GetLineKey(ctx context.Context) (string, string)
}

func NewLineKeyUsecase(repository repository.LineKeyRepository, ulidGen *util.ULIDGenerator) LineKeyUsecase {
	return &lineKeyUsecase{
		repository: repository,
		ulidGen:    ulidGen,
	}
}

func (u *lineKeyUsecase) GetLineKey(ctx context.Context) (string, string) {

	lineKey := u.ulidGen.New()
	uid, _ := ulid.Parse(lineKey)

	t := time.UnixMilli(int64(uid.Time())).Format("2006-01-02 15:04:05.000")

	log.Println("발급된 라인키 : ", lineKey)
	log.Printf("발급된 시간: %v\n", t)

	return lineKey, t
}
