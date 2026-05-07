package dto

type ResponseDTO[T any] struct {
	Result string `json:"result"` // 상태 코드
	Data   T      `json:"data"`   // 응답 데이터 (map, struct 등 자유롭게 가능)
}
