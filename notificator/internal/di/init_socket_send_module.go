package di

import (
	"notificator/internal/application/usecase"
	"notificator/internal/domain/socketSender/sender"
	"notificator/internal/infrastructure/storage"
)

type SocketSendModule struct {
	Usecase usecase.SocketSenderUsecase
}

func InitSocketSendModule(ss sender.SocketSender, chatUserStorage storage.ChatUserStorage) SocketSendModule {
	usecase := usecase.NewSocketSenderUsecase(ss, chatUserStorage)

	return SocketSendModule{
		Usecase: usecase,
	}
}
