package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"notificator/internal/application/usecase"
	"runtime"
	"time"
)

type StatsHandler struct {
	loginUsecase usecase.LoginUsecase
}

func NewStatsHandler(loginUsecase usecase.LoginUsecase) *StatsHandler {
	return &StatsHandler{
		loginUsecase: loginUsecase,
	}
}

func (h *StatsHandler) GetServerStats(w http.ResponseWriter, r *http.Request) {

	log.Println("[GetServerStats] call")

	ticker := time.NewTicker(1 * time.Second)
	const mb = 1024 * 1024
	go func() {
		for {
			select {
			case <-ticker.C:
				var m runtime.MemStats
				runtime.ReadMemStats(&m)
				active, connected, closed := h.loginUsecase.GetStats()

				log.Printf(
					"stats activeConnections=%d totalConnected=%d totalClosed=%d goroutines=%d alloc=%.2fMB heapAlloc=%.2fMB heapInuse=%.2fMB heapObjects=%d numGC=%d",
					active,
					connected,
					closed,
					runtime.NumGoroutine(),
					float64(m.Alloc)/mb,
					float64(m.HeapAlloc)/mb,
					float64(m.HeapInuse)/mb,
					m.HeapObjects,
					m.NumGC,
				)
			}
		}
	}()

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode("")

}

type StatsSnapshot struct {
	ActiveConnections int64  `json:"activeConnections"`
	TotalConnected    int64  `json:"totalConnected"`
	TotalClosed       int64  `json:"totalClosed"`
	Goroutines        int    `json:"goroutines"`
	AllocBytes        uint64 `json:"allocBytes"`
	HeapAllocBytes    uint64 `json:"heapAllocBytes"`
	HeapInuseBytes    uint64 `json:"heapInuseBytes"`
	HeapObjects       uint64 `json:"heapObjects"`
	NumGC             uint32 `json:"numGC"`
}
