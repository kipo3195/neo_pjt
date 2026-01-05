package di

import (
	"notificator/internal/application/usecase"
	"notificator/internal/domain/socketSender/sender"
	"notificator/internal/infrastructure/storage"
)

type SocketSendModule struct {
	Usecase usecase.SocketSenderUsecase
}

func InitSocketSendModule(chatDataSender sender.ChatDataSender, sendConnectionStorage storage.SendConnectionStorage, chatUserStorage storage.ChatUserStorage) SocketSendModule {
	usecase := usecase.NewSocketSenderUsecase(chatDataSender, sendConnectionStorage, chatUserStorage)

	return SocketSendModule{
		Usecase: usecase,
	}
}
