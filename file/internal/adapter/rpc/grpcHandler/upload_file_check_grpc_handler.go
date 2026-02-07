package grpcHandler

import (
	"context"
	"file/internal/application/usecase"
	"file/internal/infrastructure/pb"
	"log"
)

type UploadFileCheckGrpcHandler struct {
	pb.UnimplementedUploadFileCheckServiceServer
	usecase usecase.UploadFileCheckUsecase
}

func NewUploadFileCheckGrpcHandler(usecase usecase.UploadFileCheckUsecase) *UploadFileCheckGrpcHandler {
	return &UploadFileCheckGrpcHandler{
		usecase: usecase,
	}
}

func (h *UploadFileCheckGrpcHandler) UploadFileCheck(ctx context.Context, req *pb.UploadFileCheckRequest) (*pb.UploadFileCheckResponse, error) {

	log.Println("gRPC 요청 수신! checkDate :", req.CheckDate)

	err := h.usecase.UploadFileCheck(ctx, req.CheckDate)

	if err != nil {
		return &pb.UploadFileCheckResponse{
			Success: false,
			Message: "upload file check fail.",
		}, nil
	}

	return &pb.UploadFileCheckResponse{
		Success: true,
		Message: "upload file check success.",
	}, nil

}
