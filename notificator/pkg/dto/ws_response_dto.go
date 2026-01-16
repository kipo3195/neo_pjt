package dto

// server to server 공통 response.
type WsResponseDTO[T any] struct {
	Type      string `json:"type"`
	EventType string `json:"eventType"`
	Data      T      `json:"data"`
}

// 20260116
// 	DTO를 만들어야 하는 정확한 위치
//  UseCase에서 Sender를 호출하기 직전 (r.messageSender.SendToClient(recvUser, out))
// 이 시점부터는 “이벤트 → 패킷” 변환
// transport contract 결정
// pkg/dto 사용 가능
