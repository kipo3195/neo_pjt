package usecase

import (
	"context"
	"message/internal/application/usecase/input"
	"message/internal/domain/chatRoomTitle/entity"
	"message/internal/domain/chatRoomTitle/repository"
	"time"
)

type chatRoomTitleUsecase struct {
	repo repository.ChatRoomTitleRepository
}

type ChatRoomTitleUsecase interface {
	DeleteChatRoomTitle(ctx context.Context, input input.DeleteChatRoomTitleInput) (string, error)
	UpdateChatRoomTitle(ctx context.Context, input input.UpdateChatRoomTitleInput) (string, error)
}

func NewChatRoomTitleUsecase(repo repository.ChatRoomTitleRepository) ChatRoomTitleUsecase {
	return &chatRoomTitleUsecase{
		repo: repo,
	}
}

func (r *chatRoomTitleUsecase) UpdateChatRoomTitle(ctx context.Context, input input.UpdateChatRoomTitleInput) (string, error) {

	en := entity.MakeUpdateChatRoomTitleEntity(input.UserHash, input.Org, input.RoomKey, input.Type, input.Title)

	regDate := time.Now()
	updateDate := regDate.UTC().Format(time.RFC3339)

	en.EventDate = updateDate

	err := r.repo.UpdateChatRoomTitle(ctx, en)

	if err != nil {
		return "", err
	}

	return en.EventDate, nil
}

func (r *chatRoomTitleUsecase) DeleteChatRoomTitle(ctx context.Context, input input.DeleteChatRoomTitleInput) (string, error) {

	en := entity.MakeDeleteChatRoomTitleEntity(input.UserHash, input.Org, input.RoomKey, input.Type)

	regDate := time.Now()
	deleteDate := regDate.UTC().Format(time.RFC3339)

	en.EventDate = deleteDate

	err := r.repo.DeleteChatRoomTitle(ctx, en)

	if err != nil {
		return "", err
	}

	return en.EventDate, nil
}
