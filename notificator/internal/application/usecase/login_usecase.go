package usecase

import (
	"log"
	"notificator/internal/application/usecase/input"
	"notificator/internal/domain/login/repository"
	"sync/atomic"
)

type loginUsecase struct {
	repository  repository.LoginRepository
	serverStats *ServerStats
}

type LoginUsecase interface {
	LoginProcess(input input.LoginInput)
	AddStats()
	DelStats()
	GetStats() (active int64, connected int64, closed int64)
}

func NewLoginUsecase(repository repository.LoginRepository) LoginUsecase {
	return &loginUsecase{
		repository:  repository,
		serverStats: &ServerStats{},
	}
}

func (r *loginUsecase) AddStats() {
	r.serverStats.ActiveConnections.Add(1)
	r.serverStats.TotalConnected.Add(1)
}

func (r *loginUsecase) DelStats() {
	r.serverStats.ActiveConnections.Add(-1)
	r.serverStats.TotalClosed.Add(1)
}

func (r *loginUsecase) GetStats() (active int64, connected int64, closed int64) {
	return r.serverStats.ActiveConnections.Load(),
		r.serverStats.TotalConnected.Load(),
		r.serverStats.TotalClosed.Load()
}

func (r *loginUsecase) LoginProcess(input input.LoginInput) {

	log.Printf("login api call uuid : %s, deviceType :%s \n", input.Uuid, input.DeviceType)
}

type ServerStats struct {
	ActiveConnections atomic.Int64
	TotalConnected    atomic.Int64
	TotalClosed       atomic.Int64
}
