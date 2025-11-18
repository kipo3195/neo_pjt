package usecase

import (
	"context"
	"log"
	"notificator/internal/application/usecase/input"
	"notificator/internal/domain/note/entity"
	"notificator/internal/infrastructure/storage"

	"github.com/gorilla/websocket"
)

type noteUsecase struct {
	//repo            repository.NoteRepository
	noteUserStorage storage.NoteUserStorage
}

type NoteUsecase interface {
	SubscribeNote(input input.NoteConnectInput, conn *websocket.Conn)
	RecvChatMessage(ctx context.Context, in input.NoteMessageInput)
}

func NewNoteUsecase(noteUserStorage storage.NoteUserStorage) NoteUsecase {
	return &noteUsecase{
		noteUserStorage: noteUserStorage,
	}
}

func (u *noteUsecase) SubscribeNote(input input.NoteConnectInput, conn *websocket.Conn) {

	entity := entity.MakeSubscribeNoteEntity(input.UserHash)
	u.noteUserStorage.PutNoteConnect(entity.UserHash, conn)
}

func (u *noteUsecase) RecvChatMessage(ctx context.Context, in input.NoteMessageInput) {
	en := entity.MakeRecvNoteMessageEntity(in.Type, in.NoteKey, in.Contents, in.SendUserHash, in.RecvUserHash, in.RefeUserHash)

	// 수신자, 참조자 배열 합치기
	destUserHash := append(in.RecvUserHash, in.RefeUserHash...)

	for i := 0; i < len(destUserHash); i++ {

		// 수신자의 웹소켓 connection 객체 조회
		conn := u.noteUserStorage.GetNoteConnect(destUserHash[i])

		if conn == nil {
			continue
		}

		if err := conn.WriteJSON(en); err != nil {
			log.Printf("websocket write error to %s: %v", destUserHash[i], err)
			conn.Close()
			u.noteUserStorage.RemoveNoteConnect(destUserHash[i])
		}
	}
}
