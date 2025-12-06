package usecase

import (
	"context"
	"message/internal/application/usecase/input"
	"message/internal/consts"
	"message/internal/domain/chatRoom/entity"
	"message/internal/domain/chatRoom/repository"
	"message/internal/infrastructure/storage"
	"strings"
	"time"
)

type chatRoomUsecase struct {
	repository repository.ChatRoomRepository
	storage    storage.ChatRoomStorage
}

type ChatRoomUsecase interface {
	CreateChatRoom(ctx context.Context, input input.CreateChatRoomInput) (string, error)
}

func NewChatRoomUsecase(repository repository.ChatRoomRepository, storage storage.ChatRoomStorage) ChatRoomUsecase {

	return &chatRoomUsecase{
		repository: repository,
		storage:    storage,
	}

}

func (u *chatRoomUsecase) CreateChatRoom(ctx context.Context, input input.CreateChatRoomInput) (string, error) {

	memberEntity := make([]entity.CreateChatRoomMemberEntity, 0)

	for _, member := range input.Member {

		temp := entity.CreateChatRoomMemberEntity{
			MemberHash:      member.MemberHash,
			MemberWorksCode: member.MemberWorksCode,
		}

		memberEntity = append(memberEntity, temp)
	}
	// 시간 생성
	regDate := time.Now()
	regDateStr := regDate.UTC().Format(time.RFC3339)

	chatRoomEntity := entity.MakeCreateChatRoomEntity(input.CreateUserHash, regDate, input.RoomKey, input.RoomType, input.Title, input.Description, input.SecretFlag, input.Secret, input.WorksCode)

	// 타입 에러 체크
	upperType, err := roomTypeCheck(chatRoomEntity.RoomType)
	if err != nil {
		return "", err
	} else {
		chatRoomEntity.RoomType = upperType
	}

	// 암호 처리 여부 체크
	upperSecret, err := secretCheck(chatRoomEntity.SecretFlag, chatRoomEntity.Secret)
	if err != nil {
		return "", err
	} else {
		chatRoomEntity.SecretFlag = upperSecret
	}

	err = u.repository.PutChatRoom(ctx, memberEntity, chatRoomEntity)
	if err != nil {
		return "", err
	}

	return regDateStr, nil
}

func roomTypeCheck(roomType string) (string, error) {

	upper := strings.ToUpper(roomType)

	if upper == "O" || upper == "N" {

	} else {
		return "", consts.ErrRoomTypeCheckError
	}

	return upper, nil
}

func secretCheck(secretFlag string, secret string) (string, error) {

	upper := strings.ToUpper(secretFlag)

	if upper != "Y" && upper != "N" {
		// 시크릿 타입 에러
		return "", consts.ErrRoomSecretFlagCheckError
	}

	if upper == "Y" && secret == "" {
		return "", consts.ErrRoomSecretCheckError
	}

	return upper, nil
}
