package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

func MakeUpdateHash() string {
	// 현재 시간 밀리초 문자열
	now := time.Now().UnixNano() // 나노초 단위 시간값

	// 16바이트 랜덤 바이트 생성
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}

	// 시간값을 바이트 배열로 변환 (int64 -> []byte)
	timeBytes := []byte(fmt.Sprintf("%d", now))

	// 시간 + 랜덤 바이트 합치기
	data := append(timeBytes, randomBytes...)

	// SHA-256 해시 생성
	hash := sha256.Sum256(data)

	// 16진수 인코딩해서 문자열 반환
	return hex.EncodeToString(hash[:])
}
