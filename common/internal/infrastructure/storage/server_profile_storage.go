package storage

import (
	"bytes"
	"common/internal/consts"
	"context"
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
)

// profile 도메인의 storage에 정의된 행동 계약(Contract)만 정의
type ServerProfileStorage struct {
	ServerUrl string
}

func NewServerProfileStorage(serverUrl string) *ServerProfileStorage {
	return &ServerProfileStorage{
		ServerUrl: serverUrl,
	}
}

func (s *ServerProfileStorage) Upload(ctx context.Context, profileImg []byte, fileHash string) (string, string, error) {
	// 1. []byte → io.Reader로 변환
	imgReader := bytes.NewReader(profileImg)

	// 2. 이미지 디코딩 (JPEG, PNG, GIF 등 자동 감지)
	srcImage, format, err := image.Decode(imgReader)
	if err != nil {
		return "", "", err
	}

	// 3. 저장 디렉터리 확인
	saveDir := "./user_profile/"
	log.Println("[profile - upload] 디렉터리 체크")
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return "", "", err
	}

	// 4. 파일 명칭 및 저장경로 생성
	fileName := fileHash + "." + format
	filePath := filepath.Join(saveDir, fileName)
	log.Println("[profile - upload] 파일 경로 생성:", filePath)

	// 5. 파일 생성
	outFile, err := os.Create(filePath)
	if err != nil {
		return "", "", err
	}
	defer outFile.Close()

	// 6. 이미지 리사이즈 (가로 300px, 세로 비율 유지)
	log.Println("[profile - upload] 이미지 리사이즈 중...")
	dstImage := imaging.Resize(srcImage, 300, 300, imaging.Lanczos)

	// 7. 포맷별 인코딩
	log.Println("[profile - upload] 파일 저장 중...")
	switch format {
	case "jpeg":
		err = imaging.Encode(outFile, dstImage, imaging.JPEG)
	case "png":
		err = imaging.Encode(outFile, dstImage, imaging.PNG)
	case "gif":
		err = imaging.Encode(outFile, dstImage, imaging.GIF)
	default:
		return "", "", fmt.Errorf("지원하지 않는 포맷: %s", format)
	}
	if err != nil {
		return "", "", err
	}

	log.Println("[profile - upload] 업로드 완료")
	saveFilePath := filePath

	return saveFilePath, fileName, nil
}

func (s *ServerProfileStorage) GetProfileUrl(ctx context.Context, fileName string) ([]byte, error) {

	saveDir := "./user_profile/"
	filePath := filepath.Join(saveDir, fileName)

	// 1. 파일 존재 여부 확인
	log.Printf("[profile - load] 프로필 이미지 존재 여부 확인 %s", fileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("파일이 존재하지 않습니다: %s", fileName)
	}

	// 2. 파일 읽기
	fileData, err := os.ReadFile(filePath)
	log.Printf("[profile - load] 프로필 이미지 로드 %s", filePath)
	if err != nil {
		return nil, fmt.Errorf("파일 읽기 실패: %w", err)
	}

	// 3. 로그 및 반환
	log.Printf("[profile - load] 프로필 이미지 로드 완료: %s", filePath)
	return fileData, nil
}

func (s *ServerProfileStorage) DeleteImg(ctx context.Context, fileName string) error {

	saveDir := "./user_profile/"
	filePath := filepath.Join(saveDir, fileName)

	// 존재하지 않는 파일
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return consts.ErrProfileImgNotExist
	}

	// 파일 삭제
	if err := os.Remove(filePath); err != nil {
		return consts.ErrProfileImgRemoveError
	}

	log.Printf("[DeleteImg] fileName : %s delete success.", fileName)
	return nil
}
