package grpcHandler

import (
	"context"
	"log"
	"message/internal/application/service"
	"message/internal/infrastructure/pb"
)

type ChatLineServiceHandler struct {
	pb.UnimplementedBatchMessageServiceServer
	svc *service.ChatLineService
}

func NewChatLineServiceHandler(svc *service.ChatLineService) *ChatLineServiceHandler {
	return &ChatLineServiceHandler{
		svc: svc,
	}
}

func (r *ChatLineServiceHandler) GetSendFileInfo(ctx context.Context, in *pb.GetSendFileInfoRequest) (*pb.GetSendFileInfoResponse, error) {

	log.Println("gRPC 데이터 수신 :", in.FileIds)

	result, err := r.svc.Chat.GetSendFileInfo(ctx, in.FileIds)

	if err != nil {
		return &pb.GetSendFileInfoResponse{
			Success:  false,
			FileInfo: nil,
		}, err
	}

	fileInfo := make([]*pb.FileInfo, 0)
	for _, value := range result {

		temp := &pb.FileInfo{
			FileId:  value.FileId,
			LineKey: value.LineKey,
		}

		fileInfo = append(fileInfo, temp)
	}

	return &pb.GetSendFileInfoResponse{
		Success:  true,
		FileInfo: fileInfo,
	}, nil
}
