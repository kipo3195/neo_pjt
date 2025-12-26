package di

import (
	"notificator/internal/application/usecase"
	"notificator/internal/domain/socketSender/sender"
	"notificator/internal/infrastructure/storage"
)

type SocketSendModule struct {
	Usecase usecase.SocketSenderUsecase
}

func InitSocketSendModule(ss sender.SocketSender, sendConnectionStorage storage.SendConnectionStorage, chatUserStorage storage.ChatUserStorage) SocketSendModule {
	usecase := usecase.NewSocketSenderUsecase(ss, sendConnectionStorage, chatUserStorage)

	return SocketSendModule{
		Usecase: usecase,
	}
}
