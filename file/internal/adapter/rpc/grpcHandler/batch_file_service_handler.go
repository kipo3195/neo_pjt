package grpcHandler

import (
	"context"
	"file/internal/application/service"
	"file/internal/infrastructure/pb"
	"log"
)

type UploadFileCheckServiceHandler struct {
	pb.UnimplementedBatchFileServiceServer
	svc service.UploadFilecheckService
}

func NewUploadFileCheckServiceHandler(svc service.UploadFilecheckService) *UploadFileCheckServiceHandler {
	return &UploadFileCheckServiceHandler{
		svc: svc,
	}
}

func (h *UploadFileCheckServiceHandler) UploadFileCheck(ctx context.Context, req *pb.UploadFileCheckRequest) (*pb.UploadFileCheckResponse, error) {

	log.Println("gRPC 요청 수신! checkDate :", req.CheckDate)

	invalidFileIds, err := h.svc.UploadFileCheck.InvalidFileCheck(ctx, req.CheckDate)
	if err != nil {
		return &pb.UploadFileCheckResponse{
			Success: false,
			Message: "upload file check fail 1",
		}, nil
	}

	if len(invalidFileIds) == 0 {
		// 비정상 파일 존재 X
		return &pb.UploadFileCheckResponse{
			Success: true,
			Message: "There are not invalid file at " + req.CheckDate,
		}, nil
	}

	err = h.svc.UploadFileCheck.UpDateInvalidFileState(ctx, invalidFileIds)

	if err != nil {
		return &pb.UploadFileCheckResponse{
			Success: false,
			Message: "upload file check fail 2",
		}, nil
	}

	return &pb.UploadFileCheckResponse{
		// 모든 과정 처리 완료
		Success: true,
		Message: "upload file check success.",
	}, nil

}

func (h *UploadFileCheckServiceHandler) GetInvalidFileInfo(ctx context.Context, req *pb.GetInvalidFileInfoRequest) (*pb.GetInvalidFileInfoResponse, error) {

	log.Println("[GetInvalidFileInfo] req :", req.Yesterday)
	result, err := h.svc.ChatFile.GetInvalidFileInfo(ctx, req.Yesterday)

	log.Println("[GetInvalidFileInfo] result :", result)
	if err != nil {
		return &pb.GetInvalidFileInfoResponse{
			InvalidFileId: nil,
			Message:       "error",
		}, err
	}

	return &pb.GetInvalidFileInfoResponse{
		InvalidFileId: result,
		Message:       "success",
	}, nil

}

func (h *UploadFileCheckServiceHandler) ClearFileStorage(ctx context.Context, req *pb.ClearFileStorageRequest) (*pb.ClearFileStorageResponse, error) {

	log.Println("[ClearFileStorage] clearFileIds :", req.ClearFileIds)
	log.Println("[ClearFileStorage] sendedFileIds :", req.SendedFileIds)

	err := h.svc.ChatFile.ClearFileStorage(ctx, req.ClearFileIds, req.SendedFileIds)

	if err != nil {
		return &pb.ClearFileStorageResponse{
			Success: false,
			Message: "error",
		}, err
	}

	return &pb.ClearFileStorageResponse{
		// 모든 과정 처리 완료
		Success: true,
		Message: "success",
	}, nil

}
