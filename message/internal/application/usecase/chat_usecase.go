package usecase

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"log"
	"message/internal/application/usecase/input"
	"message/internal/application/usecase/output"
	"message/internal/consts"
	"message/internal/domain/chat/entity"
	"message/internal/domain/chat/job"
	"message/internal/domain/chat/repository"
	"message/internal/domain/logger"
	"message/internal/infrastructure/workerPool"
	"message/internal/util"
	"time"

	"github.com/nats-io/nats.go"
)

type chatUsecase struct {
	repository repository.ChatRepository
	connector  *nats.Conn
	workerPool workerPool.ChatWorkerPool
	logger     logger.Logger
}

type ChatUsecase interface {
	SendChat(ctx context.Context, in input.SendChatInput) (output.SendChatOutput, error)
	addTaskChatLineJob(chatEntity entity.SendChatEntity, chatCountEventEntity entity.ChatCountEventEntity) error
	ReadChat(ctx context.Context, in input.ReadChatInput) error
	GetChatLineEvent(ctx context.Context, in input.GetChatLineEventInput) ([]output.GetChatLineEventOutput, error)
}

func NewChatUsecase(repository repository.ChatRepository, connector *nats.Conn, workerPool workerPool.ChatWorkerPool, logger logger.Logger) ChatUsecase {
	// domain layer
	return &chatUsecase{
		repository: repository,
		connector:  connector,
		workerPool: workerPool, // Usecase는 ChatWorkerPool이라는 인터페이스에 의존하고, 이 인터페이스의 구현체가 chatWorkerPool 구조체라는 사실을 전혀 알지 못합니다.
		logger:     logger,
	}
}
func (u *chatUsecase) ReadChat(ctx context.Context, in input.ReadChatInput) error {

	readChatEntity := entity.MakeReadChatEntity(in.RoomKey, in.RoomType, in.UserHash, in.ReadDate)

	err := u.repository.ReadChatLine(ctx, readChatEntity)

	if err != nil {
		log.Println(err)
		return err
	}

	chatCountEventEntity := entity.MakeChatCountEventEntity(readChatEntity.RoomKey, readChatEntity.RoomType, "read", readChatEntity.UserHash, 0)

	data, err := util.EntityMarshal(chatCountEventEntity)
	if err != nil {
		log.Println(err)
		return err
	}

	/* 미확인 건수 발송 Message Broker */
	err = u.connector.Publish("chat.count.broadcast", data)
	if err != nil {
		log.Println("NATS publish failed:", err)
		return consts.ErrPublishToMessageBrokerError
		// 이후에 server to server rest로 전송하는 API 추가 TODO 아마도 별도의 비동기 처리로?
	}

	return nil
}

func (u *chatUsecase) SendChat(ctx context.Context, in input.SendChatInput) (output.SendChatOutput, error) {

	// 채팅 라인 entity
	chatLineEntity := entity.MakeChatLineEntity(in.ChatLine.Cmd, in.ChatLine.Contents, in.ChatLine.LineKey, in.ChatLine.TargetLineKey, in.ChatLine.SendUserHash, in.ChatLine.SendDate)
	// 채팅 룸 entity
	chatRoomEntity := entity.MakeChatRoomEntity(in.ChatRoom.RoomKey, in.ChatRoom.RoomType, in.ChatRoom.SecretFlag)
	// 채팅 파일 entity
	chatFileEntity := make([]entity.ChatFileEntity, 0)

	// 썸네일 url output용 map
	thumbnailMap := make(map[string]string)

	for _, f := range in.ChatFile {

		temp := entity.ChatFileEntity{
			FileId:   f.FileId,
			FileName: f.FileName,
			FileExt:  f.FileExt,
		}
		if temp.FileId != "" && (temp.FileExt == "jpg" || temp.FileExt == "jpeg" || temp.FileExt == "png") {
			temp.ThumbnailUrl = generateThumborURL("uc898911", temp.FileId, "300x300")
		}
		chatFileEntity = append(chatFileEntity, temp)
		thumbnailMap[temp.FileId] = temp.ThumbnailUrl
	}
	// 채팅 전송용 entity
	chatEntity := entity.MakeSendChatEntity(in.EventType, in.ChatSession, chatLineEntity, chatRoomEntity, chatFileEntity)

	// 미확인 건수 전송용 entity
	chatCountEventEntity := entity.MakeChatCountEventEntity(in.ChatRoom.RoomKey, chatRoomEntity.RoomType, "unread", in.ChatLine.SendUserHash, 1)

	log.Println("[SendChat] entity : ", chatEntity)

	data, err := util.EntityMarshal(chatEntity)
	if err != nil {
		log.Fatal(err)
		return output.SendChatOutput{}, err
	}

	log.Println("[SendChat] data : ", data)

	/* 채팅 발송 Message Broker */
	err = u.connector.Publish("chat.broadcast", data)
	if err != nil {
		log.Println("NATS publish failed:", err)
		return output.SendChatOutput{}, consts.ErrPublishToMessageBrokerError

		// 이후에 server to server rest로 전송하는 API 추가 TODO 아마도 별도의 비동기 처리로?
	}

	/* DB 저장 Task */
	err = u.addTaskChatLineJob(chatEntity, chatCountEventEntity)

	// output 생성
	out := output.SendChatOutput{
		ThumbnailMap: thumbnailMap,
	}

	return out, nil
}

func (u *chatUsecase) addTaskChatLineJob(entity entity.SendChatEntity, chatCountEventEntity entity.ChatCountEventEntity) error {

	// Context 전달
	// ※ http의 context 전달시 job 에서 repository 호출하는 DB 처리 프로세스에서 context canceled 발생.
	// 백그라운드 처리이므로 실제 http 요청에 대한 context가 아니기 때문
	// 그러므로 새로운 context 생성하여 전달하고 job의 호출이 끝난 시점에 cancel 호출로 수명 주기 관리 = AddTask 내부에서 defer cancle 필수!
	// 이후 다른 비동기 처리가 추가될때 context 를 밖에서 생성하고 각각의 task에 주입하므로써 하나의 context로 모든 task를 아우를수 있는지 점검해보기
	jobCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	job := &job.ChatLineJob{
		SendChatEntity:       entity,
		ChatCountEventEntity: chatCountEventEntity,
		Ctx:                  jobCtx,
		Cancel:               cancel,
		Connector:            u.connector,
	}

	u.workerPool.AddTask(job)
	return nil
}

func (r *chatUsecase) GetChatLineEvent(ctx context.Context, in input.GetChatLineEventInput) ([]output.GetChatLineEventOutput, error) {

	entity := entity.MakeGetChatLineEventEntity(in.ReqUserHash, in.Org, in.RoomType, in.RoomKey, in.LineKey)

	lineEvent, err := r.repository.GetChatLineEvent(ctx, entity)

	if err != nil {
		r.logger.Error(ctx, "chat_line_event_select_fail",
			"line_event", err.Error())
		return nil, consts.ErrDB
	}

	result := make([]output.GetChatLineEventOutput, 0)

	for _, v := range lineEvent {

		temp := output.GetChatLineEventOutput{
			EventType:     v.EventType,
			Cmd:           v.Cmd,
			LineKey:       v.LineKey,
			TargetLineKey: v.TargetLineKey,
			Contents:      v.Contents,
			SendUserHash:  v.SendUserHash,
			SendDate:      v.SendDate,
			FileId:        v.FileId,
			FileName:      v.FileName,
			FileType:      v.FileType,
		}

		if temp.FileId != "" && temp.FileType == "img" {
			temp.ThumbnailUrl = generateThumborURL("uc898911", v.FileId, "300x300")
			log.Println("thumbnailurl : ", temp.ThumbnailUrl)
		}

		result = append(result, temp)
	}

	return result, nil
}

func generateThumborURL(key string, imagePath string, options string) string {
	// 1. 보안 키와 경로를 조합해 HMAC-SHA1 연산 (이 과정이 매우 빠름)
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(options + "/" + imagePath))

	// 2. 결과를 Base64로 인코딩 (Thumbor 표준 규격)
	signature := base64.URLEncoding.EncodeToString(mac.Sum(nil))

	// 3. 최종 URL 완성
	return "http://172.16.10.114/thumbnail/" + signature + "/" + options + "/" + imagePath
}
