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
	GetLineKey(ctx context.Context) (string, string, error)
}

func NewLineKeyUsecase(repository repository.LineKeyRepository, ulidGen *util.ULIDGenerator) LineKeyUsecase {
	return &lineKeyUsecase{
		repository: repository,
		ulidGen:    ulidGen,
	}
}

func (u *lineKeyUsecase) GetLineKey(ctx context.Context) (string, string, error) {

	lineKey := u.ulidGen.New()
	uid, _ := ulid.Parse(lineKey)

	t := time.UnixMilli(int64(uid.Time())).Format("2006-01-02 15:04:05.000")

	log.Printf("발급된 라인 : %s, 시간: %v\n", lineKey, t)

	return lineKey, t, ctx.Err()
}
