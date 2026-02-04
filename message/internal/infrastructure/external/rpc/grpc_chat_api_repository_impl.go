package rpc

import (
	"context"
	"message/internal/domain/chat/repository"
	"message/internal/infrastructure/pb"

	"google.golang.org/grpc"
)

type grpcChatApiRepositoryImpl struct {
	client pb.FileServiceClient
}

func NewGrpcChatApiRepositoryImpl(conn *grpc.ClientConn) repository.ChatApiRepository {
	return &grpcChatApiRepositoryImpl{
		client: pb.NewFileServiceClient(conn),
	}
}

func (r *grpcChatApiRepositoryImpl) NotifySendChatFile(ctx context.Context, transactionId string) error {

	req := &pb.UpdateFileRequest{
		TransactionId: transactionId,
	}

	// 실제 gRPC 호출
	res, err := r.client.UpdateFileStatus(ctx, req)

	if err != nil {
		return err
	}

	if res.Success {
		return nil
	} else {
		return err
	}
}
