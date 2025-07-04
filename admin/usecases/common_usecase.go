package usecases

import (
	"admin/consts"
	commonDto "admin/dto/client/common"
	"admin/repositories"
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

type commonUsecase struct {
	repo repositories.CommonRepository
}

type CommonUsecase interface {
	CreateSkinImg(ctx context.Context, dto commonDto.CreateSkinImgRequest) (interface{}, error)
}

func NewCommonUsecase(repo repositories.CommonRepository) CommonUsecase {
	return &commonUsecase{
		repo: repo,
	}
}

func (r *commonUsecase) CreateSkinImg(ctx context.Context, dto commonDto.CreateSkinImgRequest) (interface{}, error) {
	defer dto.File.Close()

	fmt.Println("11")
	// 파일의 사이즈 검증
	fileSize := dto.FileInfo.Size
	sizeCheck := checkSkinImgSize(fileSize)
	if !sizeCheck {
		return nil, consts.ErrFileSizeExceeded
	}

	fmt.Println("22")
	// 파일의 확장자 검증 (이미지인지 판단.)
	detectedType, err := detectContentType(dto.File)
	if err != nil {
		return nil, consts.ErrFileExtentionDetect
	}
	fmt.Println("33")
	if !strings.HasPrefix(detectedType, "image/") {
		return nil, consts.ErrFileExtentionInvalid
	}

	fmt.Println("44")
	// 파일을 common 서비스로 전송
	err = skinImgforwardToCommon(ctx, dto)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func checkSkinImgSize(size int64) bool {
	return true
}

func detectContentType(file multipart.File) (string, error) {
	// 처음 몇 바이트를 읽어 content-type 추론
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return "", err
	}

	// 원위치로 되돌리기 (seek back to beginning)
	file.Seek(0, io.SeekStart)

	return http.DetectContentType(buffer), nil
}

func skinImgforwardToCommon(ctx context.Context, dto commonDto.CreateSkinImgRequest) error {
	fmt.Println("55")
	// 1. 파일 내용을 버퍼에 복사
	var buf bytes.Buffer
	_, err := io.Copy(&buf, dto.File)
	if err != nil {
		return fmt.Errorf("파일 버퍼 복사 실패: %w", err)
	}
	fmt.Println("66")
	// 2. multipart/form-data 본문 구성
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("File", dto.FileInfo.Filename) // common에서 수신받는 form-data key값
	if err != nil {
		return fmt.Errorf("multipart 작성 실패: %w", err)
	}
	fmt.Println("77")
	_, err = io.Copy(part, &buf)
	if err != nil {
		return fmt.Errorf("파일 데이터 삽입 실패: %w", err)
	}

	writer.Close()
	fmt.Println("88")
	// 3. HTTP 요청 생성
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://172.16.10.114/common/sv1/skin-img", &body)
	if err != nil {
		return fmt.Errorf("요청 생성 실패: %w", err)
	}
	fmt.Println("99")
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer serverToken") // 서버 to 서버 인증 처리 필요
	req.Header.Set("Skin-Type", dto.SkinType)

	// 4. 요청 전송
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("common 서비스 요청 실패: %w", err)
	}
	defer resp.Body.Close()
	fmt.Println("00 resp.StatusCode : ", resp.StatusCode)

	// 5. 응답 상태 확인
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("common 응답 실패 [%d]: %s", resp.StatusCode, string(respBody))
	}
	fmt.Println("aa")

	return nil
}
