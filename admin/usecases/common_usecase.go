package usecases

import (
	"admin/consts"
	commonDto "admin/dto/client/common"
	commonReqDto "admin/dto/client/common/request"
	"admin/repositories"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

type commonUsecase struct {
	repo repositories.CommonRepository
}

type CommonUsecase interface {
	CreateSkinImg(ctx context.Context, body commonReqDto.CreateSkinImgRequestBody) (interface{}, error)
}

func NewCommonUsecase(repo repositories.CommonRepository) CommonUsecase {
	return &commonUsecase{
		repo: repo,
	}
}

func (r *commonUsecase) CreateSkinImg(ctx context.Context, body commonReqDto.CreateSkinImgRequestBody) (interface{}, error) {

	fmt.Println("11")

	// 파일의 사이즈 검증
	sizeCheck := checkSkinImgSize(body.FileSize)
	if !sizeCheck {
		return nil, consts.ErrFileSizeExceeded
	}

	fmt.Println("22")
	// 파일의 확장자 검증 (이미지인지 판단.)
	detectedType, err := ValidateImageFile(body.FileName, body.File)
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

var (
	// 허용 가능한 확장자 및 MIME 타입 목록 -> 서버 설정 또는 메모리 로딩 필요
	allowedExtensions = map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".webp": "image/webp",
	}
)

func ValidateImageFile(fileName string, fileBytes []byte) error {
	// 확장자 소문자로 정리
	ext := strings.ToLower(filepath.Ext(fileName))

	// 1. 확장자 허용 여부 확인
	expectedMime, ok := allowedExtensions[ext]
	if !ok {
		return errors.New("허용되지 않은 파일 확장자입니다")
	}

	// 2. MIME 타입 확인 (내용 기반)
	detectedMime := http.DetectContentType(fileBytes[:min(512, len(fileBytes))])
	if detectedMime != expectedMime {
		return errors.New("파일 내용이 확장자와 일치하지 않습니다 (위조 가능성)")
	}

	// 3. (선택) 파일 이름이 이상하거나 잘못된 확장자를 포함하는 경우도 막을 수 있음
	if _, _, err := mime.ParseMediaType(detectedMime); err != nil {
		return errors.New("올바르지 않은 MIME 타입입니다")
	}

	return nil
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
