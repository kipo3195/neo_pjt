package handler

import (
	"context"
	"encoding/json"
	"log"
	"message/internal/adapter/http/dto/chatRoom"
	"message/internal/application/usecase"
	"message/internal/application/usecase/input"
	commonConsts "message/pkg/consts"
	response "message/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type ChatRoomHandler struct {
	usecase usecase.ChatRoomUsecase
}

func NewChatRoomHandler(usecase usecase.ChatRoomUsecase) *ChatRoomHandler {
	return &ChatRoomHandler{
		usecase: usecase,
	}
}

func (r *ChatRoomHandler) CreateBulkChatRoom(c *gin.Context) {

	var req chatRoom.BulkCreateChatRoomRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_108, commonConsts.E_108_MSG)
		return
	}

	makeCount := req.MakeCount
	go func() {
		workerCount := 20
		if makeCount < workerCount {
			workerCount = makeCount
		}

		baseCount := makeCount / workerCount
		remainder := makeCount % workerCount

		for workerID := 0; workerID < workerCount; workerID++ {
			workerMakeCount := baseCount
			if workerID < remainder {
				workerMakeCount++
			}

			go func(workerID int, workerMakeCount int) {
				output, err := r.usecase.CreateBulkChatRoom(context.Background(), input.BulkCreateChatRoomInput{
					MakeCount: workerMakeCount,
				})
				if err != nil {
					log.Println("[CreateBulkChatRoom] async worker failed. workerID :", workerID, "err :", err)
					return
				}
				log.Println("[CreateBulkChatRoom] async worker finished. workerID :", workerID, "created count :", output.CreatedCount)
			}(workerID, workerMakeCount)
		}
	}()

	res := chatRoom.BulkCreateChatRoomResponse{
		MakeCount:    makeCount,
		CreatedCount: 0,
		Room:         make([]chatRoom.BulkCreateChatRoomDto, 0),
	}

	response.SendSuccess(c, res)
}
