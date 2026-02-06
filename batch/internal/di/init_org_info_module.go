package di

import (
	"batch/internal/application/task"
	"batch/internal/infrastructure/persistence/repository"
	"batch/internal/infrastructure/persistence/storage"

	"gorm.io/gorm"
)

type OrgInfoModule struct {
	Task task.OrgInfoTask
}

func InitOrgInfoModule(db *gorm.DB, storage storage.OrgInfoStorage, domain string) *OrgInfoModule {

	repo := repository.NewOrgInfoRepository(db)
	apiRepo := repository.NewOrgInfoApiRepository(domain)
	task := task.NewOrgInfoTask(repo, apiRepo, storage)

	return &OrgInfoModule{
		Task: task,
	}

}
