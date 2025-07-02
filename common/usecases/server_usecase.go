package usecases

import (
	"bytes"
	"common/consts"
	dto "common/dto/common"
	adminDto "common/dto/server/admin"
	commonDto "common/dto/server/common"
	"common/entities"
	"common/infra/storage"
	"common/repositories"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/disintegration/imaging"
)

type serverUsecase struct {
	repo              repositories.ServerRepository
	configHashStorage storage.ConfigHashStorage
}

type ServerUsecase interface {
	DeviceInit(ctx context.Context, body *commonDto.DeviceInitRequest) (*entities.InitResult, *dto.ErrorResponse)
	CreateSkinImg(ctx context.Context, body adminDto.CreateSkinImgRequest) (interface{}, error)
}

func NewServerUsecase(repo repositories.ServerRepository, configHashStorage storage.ConfigHashStorage) ServerUsecase {
	return &serverUsecase{
		repo:              repo,
		configHashStorage: configHashStorage,
	}
}

func (u *serverUsecase) DeviceInit(ctx context.Context, body *commonDto.DeviceInitRequest) (*entities.InitResult, *dto.ErrorResponse) {

	// DB 조회
	result, err := u.repo.GetConnectInfo(body.WorksCode)
	if err != nil {
		return &entities.InitResult{}, &dto.ErrorResponse{
			Code:    consts.E_102,
			Message: consts.E_102_MSG,
		}
	}

	// AUTH에 JWT 요청
	result.AppToken, err = generateDeviceToken(body, result.ConnectInfo)
	if err != nil {
		return &entities.InitResult{}, &dto.ErrorResponse{
			Code:    consts.E_500,
			Message: consts.E_500_MSG,
		}
	}

	// 타임존, 언어, 앱 별 스킨 정보, 설정 정보
	worksConfig, err := u.repo.GetWorksConfig(toWorksConfigEntity(body.WorksCode, body.Device), ctx)
	if err != nil {
		return &entities.InitResult{}, &dto.ErrorResponse{
			Code:    consts.E_500,
			Message: consts.E_500_MSG,
		}
	}

	result.TimeZone = worksConfig.TimeZone
	result.Language = worksConfig.Language
	result.SkinVersion = worksConfig.SkinVersion
	result.ConfigVersion = worksConfig.ConfigVersion

	return result, nil
}

func generateDeviceToken(body *commonDto.DeviceInitRequest, serverUrl string) (string, error) {
	// 소스 모듈화 처리하기
	data := map[string]string{
		"uuid": body.Uuid,
	}

	// JSON 변환
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	fmt.Println("auth service 호출! 1")

	url := "http://" + serverUrl + "/auth/sv1/generate-device-token"
	//url := domain + "/auth/v1/generate-device-token"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer serverToken") // 서버 api key

	// POST 요청 보내기
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("auth service 호출 에러 1")
		return "", err
	}
	defer resp.Body.Close()

	// 구조체로 반환해야 하는거아닌가?
	// 서버간 통신에서 var result dto.ServerResponsed 이 구조를 사용할 것인지 고민

	// 응답 출력
	var result dto.Response // common/dto/server/auth/에 ~~~ResponseDto 생성할 것
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("serverReponse 파싱시 에러")
		return "", err
	}

	resultData, ok := result.Data.(map[string]interface{})
	if !ok {
		fmt.Println("Data 필드를 map으로 변환하는 데 실패했습니다.")
		return "", errors.New("invalid data format")
	}

	token, tokenOk := resultData["token"].(string)

	if !tokenOk {
		fmt.Println("token 또는 uuid를 string으로 변환하는 데 실패했습니다.")
		return "", errors.New("invalid token format")
	}
	fmt.Println("auth service 호출 후 발급 받은 토큰 : ", token)
	return token, nil
}

func (r *serverUsecase) CreateSkinImg(ctx context.Context, dto adminDto.CreateSkinImgRequest) (interface{}, error) {

	fmt.Println("333")
	entity, err := skinFileSaved(ctx, dto)
	if err != nil {
		fmt.Println("ddd")
	} else {
		fmt.Println("eee")
		// DB 저장
		_, err := r.repo.PutSkinFileInfo(ctx, entity)
		if err != nil {
			return nil, err
		} else {
			// 메모리 갱신
			r.configHashStorage.SkinRefresh(entity.FileHash)
			return nil, nil
		}
	}
	return nil, err
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

func skinFileSaved(ctx context.Context, dto adminDto.CreateSkinImgRequest) (*entities.SkinFileInfoEntity, error) {
	defer dto.File.Close()
	fmt.Println("444")
	// 파일의 사이즈 검증
	fileSize := dto.FileInfo.Size
	sizeCheck := checkSkinImgSize(fileSize)
	if !sizeCheck {
		return nil, consts.ErrFileSizeExceeded
	}
	fmt.Println("555")
	// 파일의 확장자 검증 (이미지인지 판단.)
	detectedType, err := detectContentType(dto.File)
	if err != nil {
		return nil, consts.ErrFileExtentionDetect
	}
	fmt.Println("666")
	if !strings.HasPrefix(detectedType, "image/") {
		return nil, consts.ErrFileExtentionInvalid
	}
	fmt.Println("777")
	// 1. 이미지 디코딩
	srcImage, format, err := image.Decode(dto.File)
	if err != nil {
		return nil, fmt.Errorf("이미지 디코딩 실패: %w", err)
	}
	fmt.Println("888")
	// 2. 리사이징 (너비 300, 비율 유지) 추후 서버 설정으로 뺄 것
	dstImage := imaging.Resize(srcImage, 300, 0, imaging.Lanczos)

	fmt.Println("999")
	// 3. 저장할 파일 이름 생성
	ext := strings.ToLower(filepath.Ext(dto.FileInfo.Filename))
	fileHash := getNow()
	fileName := fmt.Sprintf("%s_%s%s", dto.SkinType, fileHash, ext)
	localPath := filepath.Join("./skins", fileName) // 저장 경로
	fmt.Println("ext : ", ext)
	fmt.Println("fileHash : ", fileHash)
	fmt.Println("fileName : ", fileName)
	fmt.Println("localPath : ", localPath)

	if err := os.MkdirAll("./skins", 0755); err != nil {
		return nil, fmt.Errorf("디렉토리 생성 실패: %w", err)
	}
	fmt.Println("000")
	// 5. 로컬 저장
	outFile, err := os.Create(localPath)
	if err != nil {
		return nil, fmt.Errorf("파일 생성 실패: %w", err)
	}
	fmt.Println("000 format : ", format)
	defer outFile.Close()

	switch format {
	case "jpeg":
		err = imaging.Encode(outFile, dstImage, imaging.JPEG, imaging.JPEGQuality(100))
	case "png":
		err = imaging.Encode(outFile, dstImage, imaging.PNG)
	case "gif":
		err = imaging.Encode(outFile, dstImage, imaging.GIF)
	default:
		return nil, fmt.Errorf("지원하지 않는 포맷: %s", format)
	}
	fmt.Println("bbb")

	if err != nil {
		return nil, fmt.Errorf("이미지 저장 실패: %w", err)
	}
	fmt.Println("ccc")
	return &entities.SkinFileInfoEntity{
		FileHash: fileHash,
		FileName: fileName,
		SkinType: dto.SkinType,
		Device:   dto.Device,
	}, nil
}

func getNow() string {
	now := time.Now()
	formatted := now.Format(consts.YYYYMMDDHHMSS)
	return formatted
}
