package util

import "encoding/json"

// 비즈니스 로직(usecase)은 어떤 방식(라이브러리)로 marshal 하는지 몰라도 되니까
func EntityMarshal(entity interface{}) ([]byte, error) {
	data, err := json.Marshal(entity)
	return data, err
}
