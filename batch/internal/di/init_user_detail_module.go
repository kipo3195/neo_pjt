package di

import (
	"batch/internal/application/task"
	"batch/internal/infrastructure/persistence/repository"
	"batch/internal/infrastructure/persistence/storage"

	"gorm.io/gorm"
)

type UserDetailModule struct {
	Task task.UserDetailTask
}

func InitUserDetailModule(db *gorm.DB, storage storage.UserDetailStorage, domain string) *UserDetailModule {

	repo := repository.NewUserDetailRepository(db)
	apiRepo := repository.NewUserDetailApiRepository(domain)

	task := task.NewUserDetailTask(repo, apiRepo)
	return &UserDetailModule{
		Task: task,
	}

}
