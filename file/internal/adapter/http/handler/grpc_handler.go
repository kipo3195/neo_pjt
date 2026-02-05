package handler

import (
	"context"
	"file/internal/infrastructure/pb" // 생성된 pb 패키지 임포트
	"log"
)

type GrpcHandler struct {
	pb.UnimplementedFileServiceServer
}

// Message 서비스에서 정의한 UpdateFileStatus RPC 구현
func (h *GrpcHandler) UpdateFileStatus(ctx context.Context, req *pb.UpdateFileRequest) (*pb.UpdateFileResponse, error) {
	// 여기서 비즈니스 로직(UseCase)을 호출하면 됩니다.
	log.Printf("gRPC 요청 수신! TransactionId: %s", req.TransactionId)

	// 예: 파일 상태를 '전송완료'로 업데이트하는 UseCase 호출 로직...

	return &pb.UpdateFileResponse{
		Success: true,
		Message: "File status updated successfully",
	}, nil
}
