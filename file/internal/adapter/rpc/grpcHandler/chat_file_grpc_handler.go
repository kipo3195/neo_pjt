package grpcHandler

import (
	"context"
	"file/internal/application/usecase"
	"file/internal/infrastructure/pb"
	"log"
)

type ChatFileGrpcHandler struct {
	//  생성된 인터페이스의 기본 구현체를 익명 필드로 포함합니다.
	// 왜 이렇게 해야 하나요?
	// gRPC-Go 버전이 업데이트되면서, 인터페이스에 메서드를 추가할 때 모든 구현체를 일일이 수정해야 하는 번거로움을 피하고자 이 방식이 도입되었습니다.
	// 인터페이스 보호: pb.UnimplementedFileServiceServer는 인터페이스의 모든 메서드를 "Not Implemented" 에러를 반환하는 함수로 이미 구현해 두었습니다.
	// 선택적 구현: 이를 상속(Embedding)받으면, 실제로 필요한 메서드(예: UpdateFileStatus)만 오버라이딩해서 구현하면 됩니다.
	// 안전성: 나중에 .proto 파일에 새로운 RPC가 추가되어도, 이 기본 구현체가 나머지를 방어해주기 때문에 서버가 컴파일 에러 없이 빌드됩니다.

	pb.UnimplementedFileServiceServer
	usecase usecase.ChatFileUsecase
}

func NewChatFileGrpcHandler(usecase usecase.ChatFileUsecase) *ChatFileGrpcHandler {
	return &ChatFileGrpcHandler{
		usecase: usecase,
	}
}

// Message 서비스에서 정의한 UpdateFileStatus RPC 구현
func (h *ChatFileGrpcHandler) UpdateFileStatus(ctx context.Context, req *pb.UpdateFileRequest) (*pb.UpdateFileResponse, error) {
	// 여기서 비즈니스 로직(UseCase)을 호출하면 됩니다.
	log.Printf("gRPC 요청 수신! TransactionId: %s", req.TransactionId)

	err := h.usecase.UpDateFileStatus(ctx, req.TransactionId)

	if err != nil {
		log.Println("[UpdateFileStatus] err :", err)
		return &pb.UpdateFileResponse{
			Success: false,
			Message: "File status updated fail",
		}, nil
	} else {
		log.Println("[UpdateFileStatus] update success transactionId  :", req.TransactionId)
		return &pb.UpdateFileResponse{
			Success: true,
			Message: "File status updated successfully",
		}, nil
	}
}
