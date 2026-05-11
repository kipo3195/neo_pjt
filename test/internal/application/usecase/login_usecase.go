package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
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
	Pending       atomic.Int64

	latencyMu      sync.Mutex
	tokenLatencies []time.Duration
	wsLatencies    []time.Duration
	totalLatencies []time.Duration
}

func (s *Stats) RecordTokenLatency(d time.Duration) {
	s.latencyMu.Lock()
	defer s.latencyMu.Unlock()
	s.tokenLatencies = append(s.tokenLatencies, d)
}

func (s *Stats) RecordWSLatency(d time.Duration) {
	s.latencyMu.Lock()
	defer s.latencyMu.Unlock()
	s.wsLatencies = append(s.wsLatencies, d)
}

func (s *Stats) RecordTotalLatency(d time.Duration) {
	s.latencyMu.Lock()
	defer s.latencyMu.Unlock()
	s.totalLatencies = append(s.totalLatencies, d)
}

type StatsSnapshot struct {
	Attempted     int64
	Connected     int64
	Failed        int64
	Closed        int64
	CurrentActive int64
	TokenError    int64
	WSError       int64

	TokenAvgLatency time.Duration
	TokenP95Latency time.Duration

	WSAvgLatency time.Duration
	WSP95Latency time.Duration

	TotalAvgLatency time.Duration
	TotalP95Latency time.Duration
}

func (s *Stats) Snapshot() StatsSnapshot {
	s.latencyMu.Lock()
	tokenCopied := append([]time.Duration(nil), s.tokenLatencies...)
	wsCopied := append([]time.Duration(nil), s.wsLatencies...)
	totalCopied := append([]time.Duration(nil), s.totalLatencies...)
	s.latencyMu.Unlock()

	tokenAvg, tokenP95 := calcAvgP95(tokenCopied)
	wsAvg, wsP95 := calcAvgP95(wsCopied)
	totalAvg, totalP95 := calcAvgP95(totalCopied)

	return StatsSnapshot{
		Attempted:     s.Attempted.Load(),
		Connected:     s.Connected.Load(),
		Failed:        s.Failed.Load(),
		Closed:        s.Closed.Load(),
		CurrentActive: s.CurrentActive.Load(),
		TokenError:    s.TokenError.Load(),
		WSError:       s.WSError.Load(),

		TokenAvgLatency: tokenAvg,
		TokenP95Latency: tokenP95,

		WSAvgLatency: wsAvg,
		WSP95Latency: wsP95,

		TotalAvgLatency: totalAvg,
		TotalP95Latency: totalP95,
	}
}

type loginUsecase struct {
	sfg *config.ServerConfig
}
type LoginUsecase interface {
	PutLogin(ctx context.Context, input input.PutLoginInput)
	PutLoginRampUp(ctx context.Context, input input.PutLoginInput)
}

func NewLoginUsecase(sfg *config.ServerConfig) LoginUsecase {

	return &loginUsecase{
		sfg: sfg,
	}
}

func (r *loginUsecase) PutLoginRampUp(ctx context.Context, input input.PutLoginInput) {
	fmt.Println("putLogin init!")

	rps := 500

	var wg sync.WaitGroup
	stats := &Stats{}

	statsCtx, cancelStats := context.WithCancel(context.Background())
	defer cancelStats()

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				snapshot := stats.Snapshot()
				log.Printf(
					"stats attempted=%d connected=%d failed=%d closed=%d currentActive=%d tokenError=%d wsError=%d tokenAvg=%s tokenP95=%s wsAvg=%s wsP95=%s totalAvg=%s totalP95=%s",
					snapshot.Attempted,
					snapshot.Connected,
					snapshot.Failed,
					snapshot.Closed,
					snapshot.CurrentActive,
					snapshot.TokenError,
					snapshot.WSError,
					snapshot.TokenAvgLatency,
					snapshot.TokenP95Latency,
					snapshot.WSAvgLatency,
					snapshot.WSP95Latency,
					snapshot.TotalAvgLatency,
					snapshot.TotalP95Latency,
				)
			case <-statsCtx.Done():
				snapshot := stats.Snapshot()
				log.Printf(
					"stats end attempted=%d connected=%d failed=%d closed=%d currentActive=%d tokenError=%d wsError=%d tokenAvg=%s tokenP95=%s wsAvg=%s wsP95=%s totalAvg=%s totalP95=%s",
					snapshot.Attempted,
					snapshot.Connected,
					snapshot.Failed,
					snapshot.Closed,
					snapshot.CurrentActive,
					snapshot.TokenError,
					snapshot.WSError,
					snapshot.TokenAvgLatency,
					snapshot.TokenP95Latency,
					snapshot.WSAvgLatency,
					snapshot.WSP95Latency,
					snapshot.TotalAvgLatency,
					snapshot.TotalP95Latency,
				)
				return
			}
		}
	}()

	launchTicker := time.NewTicker(time.Second / time.Duration(rps))
	defer launchTicker.Stop()

	for i := 1; i <= input.ConnectionCount; i++ {
		select {
		case <-ctx.Done():
			log.Printf("putLogin launch stopped: %v", ctx.Err())
			wg.Wait()
			return
		case <-launchTicker.C:
			wg.Add(1)
			go login(&wg, r.sfg.ServerIP, i, input.UserIdPrefix, stats)
		}
	}

	wg.Wait()
	fmt.Println("putLogin complete!")
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
					"stats attempted=%d connected=%d failed=%d closed=%d currentActive=%d tokenError=%d wsError=%d tokenAvg=%s tokenP95=%s wsAvg=%s wsP95=%s totalAvg=%s totalP95=%s",
					snapshot.Attempted,
					snapshot.Connected,
					snapshot.Failed,
					snapshot.Closed,
					snapshot.CurrentActive,
					snapshot.TokenError,
					snapshot.WSError,
					snapshot.TokenAvgLatency,
					snapshot.TokenP95Latency,
					snapshot.WSAvgLatency,
					snapshot.WSP95Latency,
					snapshot.TotalAvgLatency,
					snapshot.TotalP95Latency,
				)

			case <-ctx.Done():
				snapshot := stats.Snapshot()
				log.Printf(
					"stats end. attempted=%d connected=%d failed=%d closed=%d currentActive=%d tokenError=%d wsError=%d tokenAvg=%s tokenP95=%s wsAvg=%s wsP95=%s totalAvg=%s totalP95=%s",
					snapshot.Attempted,
					snapshot.Connected,
					snapshot.Failed,
					snapshot.Closed,
					snapshot.CurrentActive,
					snapshot.TokenError,
					snapshot.WSError,
					snapshot.TokenAvgLatency,
					snapshot.TokenP95Latency,
					snapshot.WSAvgLatency,
					snapshot.WSP95Latency,
					snapshot.TotalAvgLatency,
					snapshot.TotalP95Latency,
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

	stats.Attempted.Add(1)

	totalStart := time.Now()
	userID := fmt.Sprintf("%s%06d", userIdPrefix, i)

	tokenStart := time.Now()
	// token, err := issueToken(serverIP, userID)
	token, err := issueTokenInTest(userID)
	tokenLatency := time.Since(tokenStart)
	stats.RecordTokenLatency(tokenLatency)

	if err != nil {
		stats.TokenError.Add(1)
		stats.Failed.Add(1)
		log.Printf("issueToken fail userID=%s err=%v", userID, err)
		return
	}

	wsURL := fmt.Sprintf("ws://%s/notificator/ws/connect", serverIP)

	wsStart := time.Now()
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, http.Header{
		"Authorization": []string{"Bearer " + token},
	})
	wsLatency := time.Since(wsStart)
	stats.RecordWSLatency(wsLatency)

	if err != nil {
		stats.WSError.Add(1)
		stats.Failed.Add(1)
		log.Printf("ws dial fail userID=%s err=%v", userID, err)
		return
	}

	totalLatency := time.Since(totalStart)
	stats.RecordTotalLatency(totalLatency)

	stats.Connected.Add(1)
	stats.CurrentActive.Add(1)

	defer func() {
		conn.Close()
		stats.Closed.Add(1)
		stats.CurrentActive.Add(-1)
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf("read close userID=%s err=%v", userID, err)
			return
		}
	}
}

func issueToken(serverIP string, userId string) (string, error) {

	var url string
	if serverIP == "" {
		url = fmt.Sprintf("http://172.16.10.114/auth/client/v1/user/auth/test")
	} else {
		url = fmt.Sprintf("http://%s/auth/client/v1/user", serverIP)
	}

	reqBody := dto.TokenRequest{
		UserId: userId,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	client := &http.Client{
		Timeout: 20 * time.Second,
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

func issueTokenInTest(userID string) (accessToken string, err error) {

	return userID, nil
}

func calcAvgP95(values []time.Duration) (time.Duration, time.Duration) {
	if len(values) == 0 {
		return 0, 0
	}

	copied := append([]time.Duration(nil), values...)
	sort.Slice(copied, func(i, j int) bool {
		return copied[i] < copied[j]
	})

	var sum time.Duration
	for _, v := range copied {
		sum += v
	}
	avg := sum / time.Duration(len(copied))

	idx := int(math.Ceil(float64(len(copied))*0.95)) - 1
	if idx < 0 {
		idx = 0
	}
	if idx >= len(copied) {
		idx = len(copied) - 1
	}

	return avg, copied[idx]
}
