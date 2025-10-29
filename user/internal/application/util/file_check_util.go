package util

import (
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

func CheckProfileImgSize(size int64) bool {

	if size <= 0 || size > 5242880 {
		return true
	}
	return false
}

func ValidateImageFile(fileName string, fileBytes []byte) bool {
	// 확장자 소문자로 정리
	ext := strings.ToLower(filepath.Ext(fileName))

	// 1. 확장자 허용 여부 확인
	expectedMime, ok := allowedExtensions[ext]
	if !ok {
		log.Println("허용되지 않은 파일 확장자 입니다.")
		return false
	}
	// 2. MIME 타입 확인 (내용 기반)
	detectedMime := http.DetectContentType(fileBytes[:min(512, len(fileBytes))])
	if detectedMime != expectedMime {
		log.Println("파일 내용이 확장자와 일치하지 않습니다.")
		return false
	}

	// 3. (선택) 파일 이름이 이상하거나 잘못된 확장자를 포함하는 경우도 막을 수 있음
	if _, _, err := mime.ParseMediaType(detectedMime); err != nil {
		log.Println("올바르지 않은 MIME 타입입니다.")
		return false
	}

	return true
}

var (
	// 허용 가능한 확장자 및 MIME 타입 목록 -> 서버 설정 또는 메모리 로딩 필요
	allowedExtensions = map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".webp": "image/webp",
	}
)
