package task

import "batch/internal/domain/extendDbConnect/repository"

type extendDBConnectTask struct {
}

// 설정을 어디서 주입해야되는지 점검..
type ExtendDBConnectTask interface {
	GetOrgInfo() error
	GetUserDetail() error
}

func NewExtendDBConnectTask(repo repository.ExtendDBConnectRepository) ExtendDBConnectTask {
	return &extendDBConnectTask{}
}

func (r *extendDBConnectTask) GetOrgInfo() error {
	return nil
}

func (r *extendDBConnectTask) GetUserDetail() error {
	return nil
}
