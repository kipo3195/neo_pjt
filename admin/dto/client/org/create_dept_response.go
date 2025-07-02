package org

type CreateDeptResponse struct {
	Result string `json:"result"` // 상태 코드
	Data   any    `json:"data"`   // 응답 데이터 (map, struct 등 자유롭게 가능)
}
