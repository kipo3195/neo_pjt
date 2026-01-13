package handler

import (
	"encoding/json"
	"message/internal/application/orchestrator"
	"message/internal/consts"
	"message/internal/delivery/adapter"
	"message/internal/delivery/dto/chatRoom"
	"message/internal/delivery/dto/chatRoomConfig"
	"message/internal/delivery/dto/chatRoomFixed"
	"message/internal/delivery/dto/chatRoomTitle"
	"message/internal/delivery/util"
	commonConsts "message/pkg/consts"
	response "message/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type ChatRoomServiceHandler struct {
	svc *orchestrator.ChatRoomService
}

func NewChatRoomServiceHandler(svc *orchestrator.ChatRoomService) *ChatRoomServiceHandler {
	return &ChatRoomServiceHandler{
		svc: svc,
	}
}

func (r *ChatRoomServiceHandler) CreateChatRoom(c *gin.Context) {

	ctx := c.Request.Context()

	// 사용자 정보 파싱
	createUserHash := util.GetUserHashByAccessToken(c)
	if createUserHash == "" {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_110, commonConsts.E_110_MSG)
		return
	}

	var req chatRoom.CreateChatRoomRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	// 필수 데이터 검증
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_108, commonConsts.E_108_MSG)
		return
	}

	input, err := adapter.MakeCreateChatRoomInput(createUserHash, req.RoomKey, req.Type, req.Title, req.SecretFlag, req.Secret, req.Description, req.WorksCode, req.Member)

	if err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.MESSAGE_F003, consts.MESSAGE_F003_MSG)
		return
	}

	regDate, err := r.svc.ChatRoom.CreateChatRoom(ctx, input)

	if err != nil {
		if err == consts.ErrRoomKeyAlreadyExist {
			// 룸 키 중복
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.MESSAGE_F004, consts.MESSAGE_F004_MSG)
		} else if err == consts.ErrRoomTypeCheckError {
			// 룸 타입 에러
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.MESSAGE_F005, consts.MESSAGE_F005_MSG)
		} else if err == consts.ErrRoomSecretFlagCheckError {
			// 시크릿 구분값 에러
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.MESSAGE_F006, consts.MESSAGE_F006_MSG)
		} else if err == consts.ErrRoomSecretCheckError {
			// 시크릿 데이터 에러
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.MESSAGE_F007, consts.MESSAGE_F007_MSG)
		} else {
			response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		}
		return
	}

	res := chatRoom.CreateChatRoomResponse{
		RegDate: regDate,
		RoomKey: input.RoomKey,
		Type:    input.RoomType,
	}

	response.SendSuccess(c, res)
}

func (r *ChatRoomServiceHandler) GetChatRoomDetail(c *gin.Context) {

	ctx := c.Request.Context()

	// 사용자 정보 파싱
	reqUserHash := util.GetUserHashByAccessToken(c)
	if reqUserHash == "" {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_110, commonConsts.E_110_MSG)
		return
	}

	var req chatRoom.GetChatRoomDetailRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	// 필수 데이터 검증
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_108, commonConsts.E_108_MSG)
		return
	}

	input := adapter.MakeGetChatRoomDetailInput(reqUserHash, req.RoomType, req.RoomKey)
	output, err := r.svc.ChatRoom.GetChatRoomDetail(ctx, input)

	if err != nil {
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		return
	}

	room := make([]chatRoom.GetChatRoomDetailDto, 0)
	for _, o := range output {

		detail := chatRoom.ChatRoomDetail{
			RoomKey:     o.ChatRoomDetail.RoomKey,
			Title:       o.ChatRoomDetail.Title,
			SecretFlag:  o.ChatRoomDetail.SecretFlag,
			Secret:      o.ChatRoomDetail.Secret,
			Description: o.ChatRoomDetail.Description,
			State:       o.ChatRoomDetail.State,
			WorksCode:   o.ChatRoomDetail.WorksCode,
			CreateDate:  o.ChatRoomDetail.CreateDate,
			CreateUser:  o.ChatRoomDetail.CreateUser,
			Hash:        o.ChatRoomDetail.Hash,
			Owner:       o.ChatRoomDetail.Owner,
			Type:        o.ChatRoomDetail.Type,
		}

		chatRoomFixed := chatRoomFixed.ChatRoomFixed{
			FixedFlag:  "N",
			FixedOrder: 0,
		}

		myChatRoomTitle := chatRoomTitle.ChatRoomTitle{
			Title:      o.MyChatRoomTitle.Title,
			ChangeFlag: o.MyChatRoomTitle.UpdateFlag,
			UpdateDate: o.MyChatRoomTitle.UpdateDate,
		}

		myChatRoomConfig := chatRoomConfig.ChatRoomConfig{
			MyNotiState:   "Y",
			AttentionFlag: "N",
		}

		dto := chatRoom.GetChatRoomDetailDto{
			ChatRoomDetail:   detail,
			ChatRoomFixed:    chatRoomFixed,
			MyChatRoomTitle:  myChatRoomTitle,
			MyChatRoomConfig: myChatRoomConfig,
			Member:           o.Member,
		}
		room = append(room, dto)
	}

	res := chatRoom.GetChatRoomDetailResponse{
		Room: room,
	}

	response.SendSuccess(c, res)

}

func (r *ChatRoomServiceHandler) GetChatRoomList(c *gin.Context) {
	ctx := c.Request.Context()

	// 사용자 정보 파싱
	reqUserHash := util.GetUserHashByAccessToken(c)
	if reqUserHash == "" {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_110, commonConsts.E_110_MSG)
		return
	}

	var req chatRoom.GetChatRoomListRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	// 필수 데이터 검증
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_108, commonConsts.E_108_MSG)
		return
	}

	input := adapter.MakeGetChatRoomListInput(reqUserHash, req.RoomType, req.Hash, req.ReqCount, req.Filter, req.Sorting)
	output, err := r.svc.ChatRoom.GetChatRoomList(ctx, input)

	if err != nil {
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		return
	}

	room := make([]chatRoom.GetChatRoomListDto, 0)
	for _, o := range output {

		detail := chatRoom.ChatRoomDetail{
			RoomKey:     o.ChatRoomDetail.RoomKey,
			Title:       o.ChatRoomDetail.Title,
			SecretFlag:  o.ChatRoomDetail.SecretFlag,
			Secret:      o.ChatRoomDetail.Secret,
			Description: o.ChatRoomDetail.Description,
			State:       o.ChatRoomDetail.State,
			WorksCode:   o.ChatRoomDetail.WorksCode,
			CreateDate:  o.ChatRoomDetail.CreateDate,
			CreateUser:  o.ChatRoomDetail.CreateUser,
			Hash:        o.ChatRoomDetail.Hash,
			Owner:       o.ChatRoomDetail.Owner,
			Type:        o.ChatRoomDetail.Type,
		}

		chatRoomFixed := chatRoomFixed.ChatRoomFixed{
			FixedFlag:  "N",
			FixedOrder: 0,
		}

		myChatRoomTitle := chatRoomTitle.ChatRoomTitle{
			Title:      o.MyChatRoomTitle.Title,
			ChangeFlag: o.MyChatRoomTitle.UpdateFlag,
			UpdateDate: o.MyChatRoomTitle.UpdateDate,
		}

		myChatRoomConfig := chatRoomConfig.ChatRoomConfig{
			MyNotiState:   "Y",
			AttentionFlag: "N",
		}

		dto := chatRoom.GetChatRoomListDto{
			ChatRoomDetail:   detail,
			ChatRoomFixed:    chatRoomFixed,
			MyChatRoomTitle:  myChatRoomTitle,
			MyChatRoomConfig: myChatRoomConfig,
			Member:           o.Member,
		}

		room = append(room, dto)
	}

	res := chatRoom.GetChatRoomListResponse{
		Room: room,
	}

	response.SendSuccess(c, res)

}

func (r *ChatRoomServiceHandler) GetChatRoomUpdateDate(c *gin.Context) {

	ctx := c.Request.Context()

	reqUserHash := util.GetUserHashByAccessToken(c)
	if reqUserHash == "" {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_110, commonConsts.E_110_MSG)
		return
	}

	var req chatRoom.GetChatRoomUpdateDateRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	input := adapter.MakeGetChatRoomUpdateDateInput(reqUserHash, req.Type, req.Date)
	output, err := r.svc.ChatRoom.GetChatRoomUpdateDate(ctx, input)

	if err != nil {
		if err == consts.ErrRoomUpdateDateTypeError {
			response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.FAIL, consts.MESSAGE_F012, consts.MESSAGE_F012_MSG)
		} else {
			response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		}
		return
	}

	updateDate := make([]chatRoom.ChatRoomUpdateDateDto, 0)

	for _, o := range output {

		temp := chatRoom.ChatRoomUpdateDateDto{
			RoomKey:  o.RoomKey,
			RoomType: o.RoomType,
			Detail:   o.Detail,
			Line:     o.Line,
			Owner:    o.Owner,
			Member:   o.Member,
			Unread:   o.Unread,
			Title:    o.Title,
		}
		updateDate = append(updateDate, temp)
	}

	res := chatRoom.GetChatRoomUpdateDateResponse{
		RoomUpdateDate: updateDate,
	}

	response.SendSuccess(c, res)

}

func (r *ChatRoomServiceHandler) GetChatRoomMemberReadDate(c *gin.Context) {

	ctx := c.Request.Context()

	reqUserHash := util.GetUserHashByAccessToken(c)
	if reqUserHash == "" {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_110, commonConsts.E_110_MSG)
		return
	}

	var req chatRoom.GetChatRoomMemberReadDateRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	input := adapter.MakeGetChatRoomMemberReadDateInput(req.RoomKey, req.RoomType, reqUserHash)
	out, err := r.svc.ChatRoom.GetChatRoomMemberReadDate(ctx, input)

	if err != nil {
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		return
	}

	memberReadDate := make([]chatRoom.ChatRoomMemberReadDateDto, 0)

	for _, o := range out {

		temp := chatRoom.ChatRoomMemberReadDateDto{
			MemberHash: o.MemberHash,
			ReadDate:   o.ReadDate,
		}

		memberReadDate = append(memberReadDate, temp)
	}

	res := chatRoom.GetChatRoomMemberReadDateResponse{
		MemberReadDate: memberReadDate,
	}

	response.SendSuccess(c, res)

}
