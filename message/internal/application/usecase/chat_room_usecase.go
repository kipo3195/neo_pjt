package usecase

import (
	"context"
	"message/internal/application/usecase/input"
	"message/internal/application/usecase/output"
	"message/internal/consts"
	"message/internal/domain/chatRoom/entity"
	"message/internal/domain/chatRoom/repository"
	"message/internal/infrastructure/storage"
	"message/internal/util"
	"strings"
	"time"
)

type chatRoomUsecase struct {
	repository repository.ChatRoomRepository
	storage    storage.ChatRoomStorage
}

type ChatRoomUsecase interface {
	CreateChatRoom(ctx context.Context, input input.CreateChatRoomInput) (string, error)
	GetChatRoomDetail(ctx context.Context, input input.GetChatRoomDetailInput) ([]output.GetChatRoomDetailOutput, error)
	GetChatRoomList(ctx context.Context, input input.GetChatRoomListInput) ([]output.GetChatRoomListOutput, error)
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

func (r *chatRoomUsecase) GetChatRoomDetail(ctx context.Context, input input.GetChatRoomDetailInput) ([]output.GetChatRoomDetailOutput, error) {

	entity := entity.MakeGetChatRoomDetailEntity(input.ReqUserHash, input.RoomType, input.RoomKey)
	detail, err := r.repository.GetChatRoomDetail(ctx, entity)

	if err != nil {
		return nil, err
	}

	result := make([]output.GetChatRoomDetailOutput, 0)

	for _, r := range detail {

		// member 를 ','로 split하여 리스트 생성
		memberSet := util.SplitAndMakeSet(r.Member, ",")

		temp := output.ChatRoomDetail{
			RoomKey:     r.RoomKey,
			Title:       r.Title,
			SecretFlag:  r.SecretFlag,
			Secret:      r.Secret,
			Description: r.Description,
			State:       r.State,
			WorksCode:   r.WorksCode,
			CreateDate:  r.CreateDate,
			CreateUser:  r.CreateUser,
			Hash:        r.Hash,
		}

		detail := output.GetChatRoomDetailOutput{
			ChatRoomDetail: temp,
			Member:         memberSet,
		}

		result = append(result, detail)
	}

	return result, nil

}

func (r *chatRoomUsecase) GetChatRoomList(ctx context.Context, input input.GetChatRoomListInput) ([]output.GetChatRoomListOutput, error) {

	entity := entity.MakeGetChatRoomListEntity(input.ReqUserHash, input.RoomType, input.Hash, input.ReqCount, input.Filter, input.Sorting)

	// filter, sorting에 따른 처리 필요

	list, err := r.repository.GetChatRoomList(ctx, entity)
	if err != nil {
		return nil, err
	}

	result := make([]output.GetChatRoomListOutput, 0)
	for _, r := range list {

		// member 를 ','로 split하여 리스트 생성
		memberSet := util.SplitAndMakeSet(r.Member, ",")

		temp := output.ChatRoomDetail{
			RoomKey:     r.RoomKey,
			Title:       r.Title,
			SecretFlag:  r.SecretFlag,
			Secret:      r.Secret,
			Description: r.Description,
			State:       r.State,
			WorksCode:   r.WorksCode,
			CreateDate:  r.CreateDate,
			CreateUser:  r.CreateUser,
			Hash:        r.Hash,
		}

		detail := output.GetChatRoomListOutput{
			ChatRoomDetail: temp,
			Member:         memberSet,
		}

		result = append(result, detail)
	}

	return result, nil

}
