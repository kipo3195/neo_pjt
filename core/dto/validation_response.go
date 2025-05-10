package dto

type ValidationResponse struct {
	Code int // 상태 코드
	Data any // 응답 데이터 (map, struct 등 자유롭게 가능)
}
