package usecase

import (
	"bytes"
	"context"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"admin/internal/application/usecase/input"
	"admin/internal/consts"
	"admin/internal/domain/skinImg/entity"
	"admin/internal/domain/skinImg/repository"
)

type skinImgUsecase struct {
	repository repository.SkinImgRepository
}

type SkinImgUsecase interface {
	CreateSkinImg(ctx context.Context, input input.CreateSkinImgInput) error
}

func NewSkinImgUsecase(repository repository.SkinImgRepository) SkinImgUsecase {
	return &skinImgUsecase{
		repository: repository,
	}
}

func (r *skinImgUsecase) CreateSkinImg(ctx context.Context, input input.CreateSkinImgInput) error {

	log.Println("CreateSkinImg 11")
	entity := entity.MakeCreateSkinImgEntity(input.SkinType, input.File, input.FileSize, input.FileName)

	// 파일의 사이즈 검증
	sizeCheck := checkSkinImgSize(entity.FileSize)
	if !sizeCheck {
		return consts.ErrFileSizeExceeded
	}

	log.Println("CreateSkinImg 22")
	// 파일의 확장자 검증 (이미지인지 판단.)
	extentionCheck := ValidateImageFile(entity.FileName, entity.File)
	if !extentionCheck {
		return consts.ErrFileExtentionDetect
	}

	log.Println("CreateSkinImg 33")
	// 파일을 common 서비스로 전송
	serverSending := skinImgforwardToCommon(ctx, entity)
	if !serverSending {
		return consts.ErrServerApiCallError
	}

	return nil
}

func checkSkinImgSize(size int64) bool {
	return size > 0
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

func skinImgforwardToCommon(ctx context.Context, entity entity.CreateSkinImgEntity) bool {
	log.Println("44")

	// 1. multipart/form-data 본문 구성
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	part, err := writer.CreateFormFile("File", entity.FileName) // common에서 수신받는 form-data key값
	if err != nil {
		log.Println("create form file error.")
		return false
	}
	log.Println("55")
	_, err = part.Write(entity.File)
	if err != nil {
		log.Println("file data write error.")
		return false
	}
	log.Println("66")
	writer.Close()

	// 3. HTTP 요청 생성
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://"+serverUrl+"/common/server/v1/skin-img", &requestBody)
	if err != nil {
		log.Println("request make error.")
		return false
	}
	log.Println("77")
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer serverToken") // 서버 to 서버 인증 처리 필요
	req.Header.Set("Skin-Type", entity.SkinType)

	// 4. 요청 전송
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("api call error.")
		return false
	}
	log.Println("88")
	defer resp.Body.Close()
	log.Println("99 resp.StatusCode : ", resp.StatusCode)

	// 5. 응답 상태 확인 200이 아니면.
	if resp.StatusCode != http.StatusOK {
		return false
	}

	log.Println("success.")
	return true
}
