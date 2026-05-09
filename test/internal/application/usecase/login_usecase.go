package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"sync"
	"sync/atomic"
	"test/internal/application/usecase/input"
	"test/internal/infrastructure/config"
	"test/internal/infrastructure/external/http/dto"
	"time"

	"github.com/gorilla/websocket"
)

type Stats struct {
	Attempted     atomic.Int64
	Connected     atomic.Int64
	Failed        atomic.Int64
	Closed        atomic.Int64
	CurrentActive atomic.Int64
	TokenError    atomic.Int64
	WSError       atomic.Int64

	latencyMu sync.Mutex
	latencies []time.Duration
}

func (s *Stats) RecordConnectLatency(d time.Duration) {
	s.latencyMu.Lock()
	defer s.latencyMu.Unlock()

	s.latencies = append(s.latencies, d)
}

type StatsSnapshot struct {
	Attempted     int64
	Connected     int64
	Failed        int64
	Closed        int64
	CurrentActive int64
	TokenError    int64
	WSError       int64
	AvgLatency    time.Duration
	P95Latency    time.Duration // 시작 시점 ~ 웹소켓 연결 시점까지의 기준 시간을 측정
}

func (s *Stats) Snapshot() StatsSnapshot {
	s.latencyMu.Lock()
	copied := append([]time.Duration(nil), s.latencies...)
	s.latencyMu.Unlock()

	sort.Slice(copied, func(i, j int) bool {
		return copied[i] < copied[j]
	})

	var avg time.Duration
	for _, d := range copied {
		avg += d
	}
	if len(copied) > 0 {
		avg /= time.Duration(len(copied))
	}

	var p95 time.Duration
	if len(copied) > 0 {
		idx := int(float64(len(copied))*0.95) - 1
		if idx < 0 {
			idx = 0
		}
		p95 = copied[idx]
	}

	return StatsSnapshot{
		Attempted:     s.Attempted.Load(),
		Connected:     s.Connected.Load(),
		Failed:        s.Failed.Load(),
		Closed:        s.Closed.Load(),
		CurrentActive: s.CurrentActive.Load(),
		TokenError:    s.TokenError.Load(),
		WSError:       s.WSError.Load(),
		AvgLatency:    avg,
		P95Latency:    p95,
	}
}

type loginUsecase struct {
	sfg *config.ServerConfig
}
type LoginUsecase interface {
	PutLogin(ctx context.Context, input input.PutLoginInput)
}

func NewLoginUsecase(sfg *config.ServerConfig) LoginUsecase {

	return &loginUsecase{
		sfg: sfg,
	}
}

func (r *loginUsecase) PutLogin(ctx context.Context, input input.PutLoginInput) {

	fmt.Println("putLogin init!")
	var wg sync.WaitGroup
	stats := &Stats{}
	wg.Add(input.ConnectionCount)
	for i := 1; i <= input.ConnectionCount; i++ {
		go login(&wg, r.sfg.ServerIP, i, input.UserIdPrefix, stats)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-ticker.C:
				snapshot := stats.Snapshot()
				log.Printf(
					"stats attempted=%d connected=%d failed=%d closed=%d currentActive=%d tokenError=%d wsError=%d avgLatency=%s p95Latency=%s",
					snapshot.Attempted,
					snapshot.Connected,
					snapshot.Failed,
					snapshot.Closed,
					snapshot.CurrentActive,
					snapshot.TokenError,
					snapshot.WSError,
					snapshot.AvgLatency,
					snapshot.P95Latency,
				)

			case <-ctx.Done():
				snapshot := stats.Snapshot()
				log.Printf(
					"stats end. attempted=%d connected=%d failed=%d closed=%d currentActive=%d tokenError=%d wsError=%d avgLatency=%s p95Latency=%s",
					snapshot.Attempted,
					snapshot.Connected,
					snapshot.Failed,
					snapshot.Closed,
					snapshot.CurrentActive,
					snapshot.TokenError,
					snapshot.WSError,
					snapshot.AvgLatency,
					snapshot.P95Latency,
				)
				return
			}
		}
	}()
	wg.Wait()
	fmt.Println("putLogin time out!")
	cancel()

}

func login(wg *sync.WaitGroup, serverIP string, i int, userIdPrefix string, stats *Stats) {
	defer wg.Done()
	// WebSocket 연결을 시도한 총 수
	stats.Attempted.Add(1)

	start := time.Now()
	userID := fmt.Sprintf("%s%06d", userIdPrefix, i)

	token, err := issueToken(serverIP, userID)
	if err != nil {
		// JWT 발급 API 실패 수
		stats.TokenError.Add(1)
		// 전체 실패 수
		stats.Failed.Add(1)
		return
	}
	wsUrl := fmt.Sprintf("ws://%s/notificator/ws/connect", serverIP)
	conn, _, err := websocket.DefaultDialer.Dial(wsUrl, http.Header{
		"Authorization": []string{"Bearer " + token},
	})
	if err != nil {
		// WebSocket 연결 실패 수
		stats.WSError.Add(1)
		// 전체 실패 수
		stats.Failed.Add(1)
		return
	}
	// 시작 시점 ~ 인증 토큰 발급 + 웹소켓 연결 시점까지의 기준 시간을 측정
	latency := time.Since(start)
	stats.RecordConnectLatency(latency)

	// WebSocket 연결까지 성공한 수
	stats.Connected.Add(1)
	// 현재 살아있는 WebSocket 연결 수
	stats.CurrentActive.Add(1)

	defer func() {
		conn.Close()
		// 연결됐다가 끊어진 수
		stats.Closed.Add(1)
		// 현재 살아있는 WebSocket 연결 수
		stats.CurrentActive.Add(-1)
	}()
	// 연결 유지
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			return
		}
	}
}

func issueToken(serverIP string, userId string) (string, error) {

	var url string
	if serverIP == "" {
		url = fmt.Sprintf("http://172.16.10.114/auth/client/v1/user/auth/test")
	} else {
		url = fmt.Sprintf("http://%s/auth/client/v1/user/auth/test", serverIP)
	}

	reqBody := dto.TokenRequest{
		UserId: userId,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("token API failed. status=%d body=%s", resp.StatusCode, string(respBody))
	}

	var tokenResp dto.TokenResponse
	if err := json.Unmarshal(respBody, &tokenResp); err != nil {
		return "", err
	}

	if tokenResp.Data.AccessToken == "" {
		return "", fmt.Errorf("empty token response")
	}

	return tokenResp.Data.AccessToken, nil
}
