package sender

import (
	"context"
	"notificator/internal/consts"
	"notificator/internal/domain/socketSender/entity"
	"notificator/internal/domain/socketSender/sender"
	"notificator/internal/infrastructure/dto"
)

type chatRoomDataSenderImpl struct {
}

func NewChatRoomDataSender() sender.ChatRoomDataSender {
	return &chatRoomDataSenderImpl{}
}

func (r *chatRoomDataSenderImpl) SendCreateChatRoom(ctx context.Context, recv string, entity *entity.SendConnectionEntity, en entity.CreateChatRoomEntity) error {

	res := dto.MakeCreateChatRoomResponse(
		en.CreateUserHash,
		en.RegDate,
		en.RoomKey,
		en.RoomType,
		en.Title,
		en.SecretFlag,
		en.Secret,
		en.Description,
		en.WorksCode,
	)

	select {
	case entity.Chan <- res:

	default:
		return consts.ErrSenderChannelError
	}

	return nil

}
