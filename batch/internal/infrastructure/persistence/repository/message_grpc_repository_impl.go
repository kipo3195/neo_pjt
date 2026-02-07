package repository

import (
	"batch/internal/domain/messageGrpc/repository"
	"batch/internal/infrastructure/pb"
	"context"
	"log"

	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type messageGrpcRepository struct {
	db                       *gorm.DB
	messageServiceGrpcClient pb.BatchMessageServiceClient
}

func NewChatFileRepository(db *gorm.DB, messageServiceGrpcClient *grpc.ClientConn) repository.MessageGrpcRepository {
	return &messageGrpcRepository{
		db:                       db,
		messageServiceGrpcClient: pb.NewBatchMessageServiceClient(messageServiceGrpcClient),
	}
}

func (r *messageGrpcRepository) GetSendFileInfo(ctx context.Context, checkDate string, fileIds []string) error {

	req := &pb.GetSendFileInfoRequest{
		FileIds: fileIds,
	}

	res, err := r.messageServiceGrpcClient.GetSendFileInfo(ctx, req)

	if !res.Success || err != nil {
		log.Panicln("[GetSendFileInfo] error :", err)
	}

	log.Println("[GetSendFileInfo] res :", res.FileInfo)

	return nil
}
