package sender

import (
	"notificator/internal/consts"
	"notificator/internal/domain/port"
	"notificator/internal/infrastructure/persistence/storage"
)

type messageSenderImpl struct {
	storage storage.SendConnectionStorage
}

func NewMessageSender(storage storage.SendConnectionStorage) port.MessageSender {
	return &messageSenderImpl{
		storage: storage,
	}
}

func (r *messageSenderImpl) SendToClient(userHash string, payload interface{}) error {

	connection := r.storage.GetConnection(userHash)
	if connection != nil {

		select {
		case connection.Chan <- payload:

		default:
			return consts.ErrSenderChannelError
		}
	}
	return nil
}
