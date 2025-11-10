package usecase

import (
	"context"
	"log"
	"message/internal/domain/lineKey/repository"
	"message/pkg/util"
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
	log.Println("발급된 라인키 : ", lineKey)
	return lineKey
}
