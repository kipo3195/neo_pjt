package handler

import (
	"encoding/json"
	"message/internal/application/usecase"
	"message/internal/consts"
	"message/internal/delivery/adapter"
	"message/internal/delivery/dto/note"
	commonConsts "message/pkg/consts"
	response "message/pkg/response"

	"github.com/gin-gonic/gin"
)

type NoteHandler struct {
	usecase usecase.NoteUsecase
}

func NewNoteHandler(usecase usecase.NoteUsecase) *NoteHandler {
	return &NoteHandler{
		usecase: usecase,
	}
}

func (r *NoteHandler) SendNote(c *gin.Context) {
	ctx := c.Request.Context()

	// 사용자 정보 파싱
	hash := c.Value(consts.USER_HASH)
	sendUserHash, ok := hash.(string)
	if !ok {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_110, commonConsts.E_110_MSG)
		return
	}

	var req note.SendNoteRequest

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_104, commonConsts.E_104_MSG)
		return
	}

	input := adapter.MakeSendNoteInput(sendUserHash, req.NoteKey, req.Contents, req.Type, req.RecvUserHash, req.RefeUserHash)
	err := r.usecase.SendNote(ctx, input)

	if err != nil {
		if err == consts.ErrPublishToMessageBrokerError {
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.MESSAGE_F001, consts.MESSAGE_F001_MSG)
		} else {
			response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		}
		return
	}

	response.SendSuccess(c, "")

}
