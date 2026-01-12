package port

// MessageSender는 모든 도메인 유스케이스에서 공통으로 사용합니다.
type MessageSender interface {
	SendToClient(userHash string, payload interface{}) error
}
