package di

import (
	"notificator/internal/application/usecase"
	"notificator/internal/domain/socketSender/sender"
	"notificator/internal/infrastructure/storage"
)

type SocketSendModule struct {
	Usecase usecase.SocketSenderUsecase
}

func InitSocketSendModule(chatDataSender sender.ChatDataSender, sendConnectionStorage storage.SendConnectionStorage, chatRoomStorage storage.ChatRoomStorage) SocketSendModule {
	usecase := usecase.NewSocketSenderUsecase(chatDataSender, sendConnectionStorage, chatRoomStorage)

	return SocketSendModule{
		Usecase: usecase,
	}
}
