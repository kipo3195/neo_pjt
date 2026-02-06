package task

import (
	"batch/internal/domain/messageGrpc/repository"
)

type messageGrpcTask struct {
	repository repository.MessageGrpcRepository
}

type MessageGrpcTask interface {
}

func NewMessageGrpcTask(repository repository.MessageGrpcRepository) MessageGrpcTask {

	return &messageGrpcTask{
		repository: repository,
	}
}
