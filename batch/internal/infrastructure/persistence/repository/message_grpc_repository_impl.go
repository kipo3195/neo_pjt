package repository

import (
	"batch/internal/consts"
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

func (r *messageGrpcRepository) GetSendFileInfo(ctx context.Context, checkDate string, fileIds []string) (map[string]string, error) {

	req := &pb.GetSendFileInfoRequest{
		FileIds: fileIds,
	}

	res, err := r.messageServiceGrpcClient.GetSendFileInfo(ctx, req)

	if !res.Success || err != nil {
		log.Panicf("[GetSendFileInfo] result : %v error :%s \n", res.Success, err)
		return nil, consts.ErrSendFileInfoError
	}

	result := make(map[string]string)
	for _, value := range res.FileInfo {
		result[value.FileId] = value.LineKey
	}

	log.Println("[GetSendFileInfo] result :", result)

	return result, nil
}
