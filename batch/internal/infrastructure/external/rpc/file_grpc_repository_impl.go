package rpc

import (
	"batch/internal/domain/fileGrpc/repository"
	"batch/internal/infrastructure/pb"
	"context"
	"log"

	"google.golang.org/grpc"
)

type fileGrpcRepository struct {
	fileServiceGrpcClient pb.BatchFileServiceClient
}

func NewFileGrpcRepository(conn *grpc.ClientConn) repository.FileGrpcRepository {
	return &fileGrpcRepository{
		fileServiceGrpcClient: pb.NewBatchFileServiceClient(conn),
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
	log.Println("[CheckUploadFile] res message : ", res.Message)

	if res.Success {
		return nil
	} else {
		return err
	}
}

func (r *fileGrpcRepository) GetInvalidFileInfo(ctx context.Context, yesterday string) ([]string, error) {

	req := &pb.GetInvalidFileInfoRequest{
		Yesterday: yesterday,
	}

	res, err := r.fileServiceGrpcClient.GetInvalidFileInfo(ctx, req)

	if err != nil {
		log.Printf("[GetInvalidFileInfo] 통신 에러 발생: %v", err)
		return nil, err
		// 에러가 있으면 바로 return
	}

	result := make([]string, 0)

	for _, value := range res.InvalidFileId {
		result = append(result, value)
	}

	return result, nil

}

func (r *fileGrpcRepository) ClearFileStorage(ctx context.Context, clearFileId []string, sendedFileId []string) error {

	req := &pb.ClearFileStorageRequest{
		ClearFileIds:  clearFileId,
		SendedFileIds: sendedFileId,
	}

	res, err := r.fileServiceGrpcClient.ClearFileStorage(ctx, req)

	if err != nil {
		log.Printf("[ClearFileStorage] 통신 에러 발생: %v", err)
		return err // 에러가 있으면 여
		//
		// 기서 바로 리턴해야 합니다.
	}

	if res.Success {
		return nil
	} else {
		return err
	}
}
