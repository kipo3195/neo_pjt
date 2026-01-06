package usecase

import (
	"notificator/internal/domain/chatRoom/repository"
	"notificator/internal/infrastructure/storage"
)

type chatRoomUsecase struct {
	repo            repository.ChatRoomRepository
	chatRoomStorage storage.ChatRoomStorage
}

type ChatRoomUsecase interface {
	SubscribeChat(userHash string) error
	UnSubscribeChat(userHash string)
}

func NewChatRoomUsecase(repo repository.ChatRoomRepository, chatRoomStorage storage.ChatRoomStorage) ChatRoomUsecase {
	return &chatRoomUsecase{
		repo:            repo,
		chatRoomStorage: chatRoomStorage,
	}
}

func (u *chatRoomUsecase) SubscribeChat(userHash string) error {

	myChatRoom, err := u.repo.GetMyChatRoom(userHash)
	if err != nil {
		return err
	}

	u.chatRoomStorage.InitMyRoom(myChatRoom.RoomKey, userHash)
	return nil
}

func (u *chatRoomUsecase) UnSubscribeChat(userHash string) {
	u.chatRoomStorage.CleanUpMyRoom(userHash)
}
