package usecase

import (
	"context"
	"log"
	"message/internal/application/usecase/input"
	"message/internal/application/usecase/output"
	"message/internal/consts"
	"message/internal/domain/chatRoom/entity"
	"message/internal/domain/chatRoom/repository"
	"message/internal/infrastructure/storage"
	"message/internal/util"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

type chatRoomUsecase struct {
	repository repository.ChatRoomRepository
	connector  *nats.Conn
	storage    storage.ChatRoomStorage
}

type ChatRoomUsecase interface {
	CreateChatRoom(ctx context.Context, input input.CreateChatRoomInput) (string, error)
	GetChatRoomDetail(ctx context.Context, input input.GetChatRoomDetailInput) ([]output.GetChatRoomDetailOutput, error)
	GetChatRoomList(ctx context.Context, input input.GetChatRoomListInput) ([]output.GetChatRoomListOutput, error)
	GetChatRoomUpdateDate(ctx context.Context, input input.GetChatRoomUpdateInput) ([]output.GetChatRoomUpdateDateOutput, error)
}

func NewChatRoomUsecase(repository repository.ChatRoomRepository, connector *nats.Conn, storage storage.ChatRoomStorage) ChatRoomUsecase {

	return &chatRoomUsecase{
		repository: repository,
		storage:    storage,
		connector:  connector,
	}

}

func (u *chatRoomUsecase) CreateChatRoom(ctx context.Context, input input.CreateChatRoomInput) (string, error) {

	// member entity
	memberEntity := make([]entity.CreateChatRoomMemberEntity, 0)

	// 중복 제거
	unique := make(map[string]struct{}, 0)
	for _, member := range input.Member {

		if _, exists := unique[member.MemberHash]; exists {
			continue
		}
		temp := entity.CreateChatRoomMemberEntity{
			MemberHash:      member.MemberHash,
			MemberWorksCode: member.MemberWorksCode,
		}
		memberEntity = append(memberEntity, temp)
		unique[member.MemberHash] = struct{}{}
	}

	log.Println("[CreateChatRoom] chat room member : ", unique)

	// 시간 생성
	regDate := time.Now()
	regDateStr := regDate.UTC().Format(time.RFC3339)

	// chat room entity
	chatRoomEntity := entity.MakeCreateChatRoomEntity(input.CreateUserHash, regDate, input.RoomKey, input.RoomType, input.Title, input.Description, input.SecretFlag, input.Secret, input.WorksCode)

	// member + chat room
	CreateChatRoomEntity := entity.CreateChatRoomEntity{
		ChatRoomEntity:       chatRoomEntity,
		ChatRoomMemberEntity: memberEntity,
	}

	// 타입 에러 체크
	upperType, err := roomTypeCheck(CreateChatRoomEntity.ChatRoomEntity.RoomType)
	if err != nil {
		return "", err
	} else {
		CreateChatRoomEntity.ChatRoomEntity.RoomType = upperType
	}

	// 암호 처리 여부 체크
	upperSecret, err := secretCheck(CreateChatRoomEntity.ChatRoomEntity.SecretFlag, CreateChatRoomEntity.ChatRoomEntity.Secret)
	if err != nil {
		return "", err
	} else {
		CreateChatRoomEntity.ChatRoomEntity.SecretFlag = upperSecret
	}

	err = u.repository.PutChatRoom(ctx, CreateChatRoomEntity)
	if err != nil {
		return "", err
	}

	data, err := util.EntityMarshal(CreateChatRoomEntity)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	/* 채팅방 생성 이벤트 발송 Message Broker */
	msg, err := u.connector.Request("create.chat.room.message", data, 5*time.Second)
	if err != nil {
		if err == nats.ErrNoResponders {
			// 아무도 수신하지 않았으므로 재처리 혹은 server to server 처리 필요 혹은 별도 정책 정의하기.
			log.Fatal("NATS publish failed:", err)
		} else {
			log.Fatal("NATS publish failed:", err)
		}
		return "", consts.ErrPublishToMessageBrokerError
		// 이후에 server to server rest로 전송하는 API 추가 TODO 아마도 별도의 비동기 처리로?
	}
	// 기존 publish는 던지고 잊는 구조이므로 누가 수신을 했는지에 대한 정보가 없음.. 그래서 Request로 변경
	log.Println("[CreateChatRoom] recv notificator response :", string(msg.Data))

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
		owner := util.SplitAndMakeSet(r.Owner, ",")

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
			Owner:       owner,
			Type:        r.Type,
		}

		title := output.ChatRoomTitleOutput{
			Title:      r.MyRoomTitle,
			UpdateFlag: r.TitleUpdateFlag,
			UpdateDate: r.TitleUpdateDate,
		}

		roomInfo := output.GetChatRoomDetailOutput{
			ChatRoomDetail:  temp,
			Member:          memberSet,
			MyChatRoomTitle: title,
		}

		result = append(result, roomInfo)
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
		owner := util.SplitAndMakeSet(r.Owner, ",")

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
			Owner:       owner,
			Type:        r.Type,
		}

		title := output.ChatRoomTitleOutput{
			Title:      r.MyRoomTitle,
			UpdateFlag: r.TitleUpdateFlag,
			UpdateDate: r.TitleUpdateDate,
		}

		roomInfo := output.GetChatRoomListOutput{
			ChatRoomDetail:  temp,
			Member:          memberSet,
			MyChatRoomTitle: title,
		}

		result = append(result, roomInfo)
	}

	return result, nil

}

func (r *chatRoomUsecase) GetChatRoomUpdateDate(ctx context.Context, input input.GetChatRoomUpdateInput) ([]output.GetChatRoomUpdateDateOutput, error) {

	en := entity.MakeGetChatRoomUpdateDateEntity(input.ReqUserHash, input.Type, input.Date)

	updateDate, err := r.repository.GetChatRoomUpdateDate(ctx, en)

	if err != nil {
		return nil, err
	}

	result := make([]output.GetChatRoomUpdateDateOutput, 0)
	for _, rd := range updateDate {

		temp := output.GetChatRoomUpdateDateOutput{
			RoomKey:  rd.RoomKey,
			RoomType: rd.RoomType,
			Detail:   rd.Detail,
			Line:     rd.Line,
			Member:   rd.Member,
			Owner:    rd.Owner,
			Title:    rd.Title,
		}

		result = append(result, temp)
	}

	return result, nil
}
