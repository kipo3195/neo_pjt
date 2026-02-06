package task

import (
	"batch/internal/domain/fileGrpc/repository"
)

type fileGrpcTask struct {
	repository repository.FileGrpcRepository
}

type FileGrpcTask interface {
}

func NewFileGrpcTask(repository repository.FileGrpcRepository) FileGrpcTask {

	return &fileGrpcTask{
		repository: repository,
	}
}
