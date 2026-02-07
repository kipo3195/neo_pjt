package repository

import (
	"batch/internal/domain/fileGrpc/repository"
	"batch/internal/infrastructure/pb"
	"context"
	"log"

	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type fileGrpcRepository struct {
	db                    *gorm.DB
	fileServiceGrpcClient pb.UploadFileCheckServiceClient
}

func NewFileGrpcRepository(db *gorm.DB, conn *grpc.ClientConn) repository.FileGrpcRepository {
	return &fileGrpcRepository{
		db:                    db,
		fileServiceGrpcClient: pb.NewUploadFileCheckServiceClient(conn),
	}
}

func (r *fileGrpcRepository) CheckUploadFile(ctx context.Context, checkDate string) error {

	req := &pb.UploadFileCheckRequest{
		CheckDate: checkDate,
	}

	res, err := r.fileServiceGrpcClient.UploadFileCheck(ctx, req)

	if err != nil {
		log.Printf("[CheckUploadFile] 통신 에러 발생: %v", err)
		return err // 에러가 있으면 여기서 바로 리턴해야 합니다.
	}

	log.Println("[CheckUploadFile] res : ", res.Success)

	if res.Success {
		return nil
	} else {
		return err
	}
}
