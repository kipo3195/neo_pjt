package di

import (
	"notificator/internal/application/usecase"
	"notificator/internal/infrastructure/persistence/storage"
)

type SocketSendModule struct {
	Usecase usecase.SocketSenderUsecase
}

func InitSocketSendModule(sendConnectionStorage storage.SendConnectionStorage) SocketSendModule {
	usecase := usecase.NewSocketSenderUsecase(sendConnectionStorage)

	return SocketSendModule{
		Usecase: usecase,
	}
}
