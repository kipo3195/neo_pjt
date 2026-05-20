package usecase

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"test/internal/application/usecase/input"
	"test/internal/infrastructure/config"
	"test/internal/infrastructure/external/http/dto"
	"time"

	"github.com/gorilla/websocket"
)

type Stats struct {
	Attempted        atomic.Int64
	Connected        atomic.Int64
	Failed           atomic.Int64
	Closed           atomic.Int64
	CurrentActive    atomic.Int64
	TokenError       atomic.Int64
	WSError          atomic.Int64
	EventReceived    atomic.Int64
	EventParseError  atomic.Int64
	Pending          atomic.Int64
	ChatReceived     atomic.Int64
	ChatLatencyError atomic.Int64

	latencyMu      sync.Mutex
	tokenLatencies []time.Duration
	wsLatencies    []time.Duration
	chatLatencies  []time.Duration
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

func (s *Stats) RecordChatLatency(d time.Duration) {
	s.latencyMu.Lock()
	defer s.latencyMu.Unlock()
	s.chatLatencies = append(s.chatLatencies, d)
}

type StatsSnapshot struct {
	Attempted       int64
	Connected       int64
	Failed          int64
	Closed          int64
	CurrentActive   int64
	TokenError      int64
	WSError         int64
	EventReceived   int64
	EventParseError int64
	ChatReceived    int64

	TokenAvgLatency time.Duration
	TokenP95Latency time.Duration

	WSAvgLatency time.Duration
	WSP95Latency time.Duration

	TotalAvgLatency  time.Duration
	TotalP95Latency  time.Duration
	chatLatencies    []time.Duration
	ChatLatencyError int64
	ChatAvgLatency   time.Duration
	ChatP95Latency   time.Duration
}

func (s *Stats) Snapshot() StatsSnapshot {
	s.latencyMu.Lock()
	tokenCopied := append([]time.Duration(nil), s.tokenLatencies...)
	wsCopied := append([]time.Duration(nil), s.wsLatencies...)
	totalCopied := append([]time.Duration(nil), s.totalLatencies...)
	chatCopied := append([]time.Duration(nil), s.chatLatencies...)
	s.latencyMu.Unlock()

	tokenAvg, tokenP95 := calcAvgP95(tokenCopied)
	wsAvg, wsP95 := calcAvgP95(wsCopied)
	totalAvg, totalP95 := calcAvgP95(totalCopied)
	chatAvg, chatP95 := calcAvgP95(chatCopied)

	return StatsSnapshot{
		Attempted:       s.Attempted.Load(),
		Connected:       s.Connected.Load(),
		Failed:          s.Failed.Load(),
		Closed:          s.Closed.Load(),
		CurrentActive:   s.CurrentActive.Load(),
		TokenError:      s.TokenError.Load(),
		WSError:         s.WSError.Load(),
		EventReceived:   s.EventReceived.Load(),
		EventParseError: s.EventParseError.Load(),
		ChatReceived:    s.ChatReceived.Load(),

		TokenAvgLatency: tokenAvg,
		TokenP95Latency: tokenP95,

		WSAvgLatency:   wsAvg,
		WSP95Latency:   wsP95,
		ChatAvgLatency: chatAvg,
		ChatP95Latency: chatP95,

		TotalAvgLatency: totalAvg,
		TotalP95Latency: totalP95,
	}
}

type loginUsecase struct {
	sfg *config.ServerConfig
}

type wsResponseDTO[T any] struct {
	Type      string `json:"type"`
	EventType string `json:"eventType"`
	Data      T      `json:"data"`
}

type chatMessageOutput struct {
	ChatSession  string               `json:"chatSession"`
	ChatRoomData chatRoomDataOutput   `json:"chatRoomData"`
	ChatLineData chatLineDataOutput   `json:"chatLineData"`
	ChatFileData []chatFileDataOutput `json:"chatFileData,omitempty"`
}

type chatRoomDataOutput struct {
	RoomType   string `json:"roomType"`
	RoomKey    string `json:"roomKey"`
	SecretFlag string `json:"secretFlag"`
}

type chatLineDataOutput struct {
	Cmd           int    `json:"cmd"`
	Contents      string `json:"contents"`
	LineKey       string `json:"lineKey"`
	TargetLineKey string `json:"targetLineKey"`
	SendUserHash  string `json:"sendUserHash"`
	SendDate      string `json:"sendDate"`
}

type chatFileDataOutput struct {
	FileId       string `json:"fileId"`
	FileType     string `json:"fileType"`
	FileName     string `json:"fileName"`
	ThumbnailUrl string `json:"thumbnailUrl,omitempty"`
}

type LoginUsecase interface {
	PutLogin(ctx context.Context, input input.PutLoginInput)
	PutLoginRampUp(ctx context.Context, input input.PutLoginInput)
	PutLoginRampUpAndRecvEvent(ctx context.Context, input input.PutLoginInput)
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
				//snapshot := stats.Snapshot()
				// log.Printf(
				// 	"stats attempted=%d connected=%d failed=%d closed=%d currentActive=%d tokenError=%d wsError=%d eventReceived=%d eventParseError=%d tokenAvg=%s tokenP95=%s wsAvg=%s wsP95=%s totalAvg=%s totalP95=%s",
				// 	snapshot.Attempted,
				// 	snapshot.Connected,
				// 	snapshot.Failed,
				// 	snapshot.Closed,
				// 	snapshot.CurrentActive,
				// 	snapshot.TokenError,
				// 	snapshot.WSError,
				// 	snapshot.EventReceived,
				// 	snapshot.EventParseError,
				// 	snapshot.TokenAvgLatency,
				// 	snapshot.TokenP95Latency,
				// 	snapshot.WSAvgLatency,
				// 	snapshot.WSP95Latency,
				// 	snapshot.TotalAvgLatency,
				// 	snapshot.TotalP95Latency,
				// )
			case <-statsCtx.Done():
				snapshot := stats.Snapshot()
				log.Printf(
					"stats end attempted=%d connected=%d failed=%d closed=%d currentActive=%d tokenError=%d wsError=%d eventReceived=%d eventParseError=%d tokenAvg=%s tokenP95=%s wsAvg=%s wsP95=%s totalAvg=%s totalP95=%s",
					snapshot.Attempted,
					snapshot.Connected,
					snapshot.Failed,
					snapshot.Closed,
					snapshot.CurrentActive,
					snapshot.TokenError,
					snapshot.WSError,
					snapshot.EventReceived,
					snapshot.EventParseError,
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

func (r *loginUsecase) PutLoginRampUpAndRecvEvent(ctx context.Context, input input.PutLoginInput) {
	fmt.Println("putLoginAndRecvEvent init!")

	testStart := time.Now()

	// 동일한 라인키의 채팅 개수 점검 (발신과 수신의 일치 점검)
	recvLine := make(chan string, input.ConnectionCount*10)
	recvDone := make(chan int, 1)

	go func() {
		receivedLines := make(map[string]struct{})
		for lineKey := range recvLine {
			receivedLines[lineKey] = struct{}{}
		}
		recvDone <- len(receivedLines)
	}()

	rps := 500

	var wg sync.WaitGroup
	stats := &Stats{}

	statsCtx, cancelStats := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancelStats()

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// snapshot := stats.Snapshot()
				// log.Printf(
				// 	"stats attempted=%d connected=%d failed=%d closed=%d currentActive=%d tokenError=%d wsError=%d eventReceived=%d eventParseError=%d tokenAvg=%s tokenP95=%s wsAvg=%s wsP95=%s totalAvg=%s totalP95=%s",
				// 	snapshot.Attempted,
				// 	snapshot.Connected,
				// 	snapshot.Failed,
				// 	snapshot.Closed,
				// 	snapshot.CurrentActive,
				// 	snapshot.TokenError,
				// 	snapshot.WSError,
				// 	snapshot.EventReceived,
				// 	snapshot.EventParseError,
				// 	snapshot.TokenAvgLatency,
				// 	snapshot.TokenP95Latency,
				// 	snapshot.WSAvgLatency,
				// 	snapshot.WSP95Latency,
				// 	snapshot.TotalAvgLatency,
				// 	snapshot.TotalP95Latency,
				// )
			case <-statsCtx.Done():
				snapshot := stats.Snapshot()
				log.Printf(
					"!!! stats end attempted=%d connected=%d failed=%d closed=%d currentActive=%d tokenError=%d wsError=%d eventReceived=%d chatReceived=%d chatLatencyError=%d eventParseError=%d chatAvg=%s chatP95=%s tokenAvg=%s tokenP95=%s wsAvg=%s wsP95=%s totalAvg=%s totalP95=%s",
					snapshot.Attempted,
					snapshot.Connected,
					snapshot.Failed,
					snapshot.Closed,
					snapshot.CurrentActive,
					snapshot.TokenError,
					snapshot.WSError,
					snapshot.EventReceived,
					snapshot.ChatReceived,
					snapshot.ChatLatencyError,
					snapshot.EventParseError,
					snapshot.ChatAvgLatency,
					snapshot.ChatP95Latency,
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
			log.Printf("putLoginAndRecvEvent launch stopped: %v", ctx.Err())
			wg.Wait()
			return
		case <-launchTicker.C:
			wg.Add(1)
			go loginAndRecvWithTokenIssuer(
				&wg,
				r.sfg.ServerIP,
				i,
				input.UserIdPrefix,
				stats,
				r.issueTestToken,
				recvLine,
			)
		}
	}

	wg.Wait()
	// 모든 고루틴 종료 후 recvLine 채널 닫고, map에 라인키의 수 조회
	close(recvLine)
	uniqueReceivedLineCount := <-recvDone
	testEnd := time.Since(testStart)

	log.Printf("testend=%d", testEnd)
	log.Printf("unique received chat line count=%d", uniqueReceivedLineCount)

	snapshot := stats.Snapshot()
	log.Printf(
		"!!! stats end attempted=%d connected=%d failed=%d closed=%d currentActive=%d tokenError=%d wsError=%d eventReceived=%d chatReceived=%d chatLatencyError=%d eventParseError=%d chatAvg=%s chatP95=%s tokenAvg=%s tokenP95=%s wsAvg=%s wsP95=%s totalAvg=%s totalP95=%s",
		snapshot.Attempted,
		snapshot.Connected,
		snapshot.Failed,
		snapshot.Closed,
		snapshot.CurrentActive,
		snapshot.TokenError,
		snapshot.WSError,
		snapshot.EventReceived,
		snapshot.ChatReceived,
		snapshot.ChatLatencyError,
		snapshot.EventParseError,
		snapshot.ChatAvgLatency,
		snapshot.ChatP95Latency,
		snapshot.TokenAvgLatency,
		snapshot.TokenP95Latency,
		snapshot.WSAvgLatency,
		snapshot.WSP95Latency,
		snapshot.TotalAvgLatency,
		snapshot.TotalP95Latency,
	)
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
				// snapshot := stats.Snapshot()
				// log.Printf(
				// 	"stats attempted=%d connected=%d failed=%d closed=%d currentActive=%d tokenError=%d wsError=%d tokenAvg=%s tokenP95=%s wsAvg=%s wsP95=%s totalAvg=%s totalP95=%s",
				// 	snapshot.Attempted,
				// 	snapshot.Connected,
				// 	snapshot.Failed,
				// 	snapshot.Closed,
				// 	snapshot.CurrentActive,
				// 	snapshot.TokenError,
				// 	snapshot.WSError,
				// 	snapshot.TokenAvgLatency,
				// 	snapshot.TokenP95Latency,
				// 	snapshot.WSAvgLatency,
				// 	snapshot.WSP95Latency,
				// 	snapshot.TotalAvgLatency,
				// 	snapshot.TotalP95Latency,
				// )

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

type tokenIssuer func(serverIP string, userID string) (string, error)

func login(wg *sync.WaitGroup, serverIP string, i int, userIdPrefix string, stats *Stats) {
	loginWithTokenIssuer(wg, serverIP, i, userIdPrefix, stats, issueToken)
}

func loginWithTokenIssuer(wg *sync.WaitGroup, serverIP string, i int, userIdPrefix string, stats *Stats, issue tokenIssuer) {
	defer wg.Done()

	stats.Attempted.Add(1)

	totalStart := time.Now()
	userID := fmt.Sprintf("%s%06d", userIdPrefix, i)

	tokenStart := time.Now()
	token, err := issue(serverIP, userID)
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

func loginAndRecvWithTokenIssuer(wg *sync.WaitGroup, serverIP string, i int, userIdPrefix string, stats *Stats, issue tokenIssuer, recvLine chan string) {
	defer wg.Done()

	stats.Attempted.Add(1)

	totalStart := time.Now()
	userID := fmt.Sprintf("%s%06d", userIdPrefix, i)

	tokenStart := time.Now()
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
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("read close userID=%s err=%v", userID, err)
			return
		}

		res, ok, err := parseChatMessageResponse(msg)
		if err != nil {
			stats.EventParseError.Add(1)
			log.Printf("read parse fail userID=%s err=%v payload=%s", userID, err, trimPayload(msg, 512))
			continue
		}
		if !ok {
			continue
		}

		stats.EventReceived.Add(1)
		// 미확인 건수 수신
		if res.Data.ChatSession == "" && res.Data.ChatRoomData.RoomKey == "" && res.Data.ChatLineData.LineKey == "" {
			log.Printf("read empty chat message userID=%s type=%s eventType=%s", userID, res.Type, res.EventType)
		} else {
			// 채팅 데이터 수신
			lineKey := res.Data.ChatLineData.LineKey

			if lineKey != "" {
				recvLine <- lineKey
			}

			log.Printf(
				"read chat message userID=%s type=%s eventType=%s lineKey=%s",
				userID,
				res.Type,
				res.EventType,
				lineKey,
			)

			stats.ChatReceived.Add(1)

			createdAt, err := parseULIDTime(lineKey)
			if err != nil {
				stats.ChatLatencyError.Add(1)
				//log.Printf("parse lineKey ulid fail userID=%s lineKey=%s err=%v", userID, lineKey, err)
				continue
			}
			now := time.Now()
			latency := time.Since(createdAt)
			if latency < 0 {
				stats.ChatLatencyError.Add(1)
				log.Printf(
					"chat latency negative userID=%s lineKey=%s clientNow=%s ulidCreatedAt=%s latency=%s",
					userID,
					lineKey,
					now.Format(time.RFC3339Nano),
					createdAt.Format(time.RFC3339Nano),
					latency,
				)
				continue
			}

			stats.RecordChatLatency(latency)
		}
	}
}

func parseULIDTime(lineKey string) (time.Time, error) {
	if len(lineKey) < 10 {
		return time.Time{}, fmt.Errorf("invalid ulid length")
	}

	var ms uint64
	for _, c := range lineKey[:10] {
		v, ok := crockfordBase32Value(c)
		if !ok {
			return time.Time{}, fmt.Errorf("invalid ulid char %q", c)
		}
		ms = (ms << 5) | uint64(v)
	}

	// ULID timestamp is 48 bits. The first 10 base32 chars contain 50 bits,
	// so the top 2 bits must be discarded.

	return time.UnixMilli(int64(ms)), nil
}

func crockfordBase32Value(c rune) (byte, bool) {
	switch {
	case c >= '0' && c <= '9':
		return byte(c - '0'), true
	case c >= 'A' && c <= 'Z':
		return crockfordBase32Upper(c)
	case c >= 'a' && c <= 'z':
		return crockfordBase32Upper(c - 'a' + 'A')
	default:
		return 0, false
	}
}

func crockfordBase32Upper(c rune) (byte, bool) {
	const alphabet = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"
	for i, a := range alphabet {
		if c == a {
			return byte(i), true
		}
	}
	return 0, false
}

func parseChatMessageResponse(msg []byte) (wsResponseDTO[chatMessageOutput], bool, error) {
	var raw wsResponseDTO[json.RawMessage]
	if err := json.Unmarshal(msg, &raw); err != nil {
		return wsResponseDTO[chatMessageOutput]{}, false, err
	}

	if raw.Type == "" && raw.EventType == "" {
		return wsResponseDTO[chatMessageOutput]{}, false, nil
	}

	var data chatMessageOutput
	if len(raw.Data) > 0 && string(raw.Data) != "null" {
		if err := json.Unmarshal(raw.Data, &data); err != nil {
			return wsResponseDTO[chatMessageOutput]{}, false, err
		}
	}

	return wsResponseDTO[chatMessageOutput]{
		Type:      raw.Type,
		EventType: raw.EventType,
		Data:      data,
	}, true, nil
}

func trimPayload(payload []byte, limit int) string {
	if len(payload) <= limit {
		return string(payload)
	}
	return string(payload[:limit]) + "...(truncated)"
}

func issueToken(serverIP string, userId string) (string, error) {

	url := fmt.Sprintf("http://172.16.10.114/auth/client/v1/user/auth/test")

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

func (r *loginUsecase) issueTestToken(_ string, userID string) (string, error) {
	if r.sfg.AccessTokenHash == "" {
		return "", fmt.Errorf("ACCESS_TOKEN_HASH is empty")
	}

	now := time.Now()
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}
	claims := map[string]interface{}{
		"id":   userID,
		"uuid": fmt.Sprintf("%s-test-device", userID),
		"hash": userID,
		"iss":  "device",
		"iat":  now.Unix(),
		"exp":  now.Add(24 * time.Hour).Unix(),
	}

	headerBytes, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	claimsBytes, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	unsignedToken := strings.Join([]string{
		base64.RawURLEncoding.EncodeToString(headerBytes),
		base64.RawURLEncoding.EncodeToString(claimsBytes),
	}, ".")

	mac := hmac.New(sha256.New, []byte(r.sfg.AccessTokenHash))
	if _, err := mac.Write([]byte(unsignedToken)); err != nil {
		return "", err
	}

	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return unsignedToken + "." + signature, nil
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
