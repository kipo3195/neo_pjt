package usecases

import (
	"bytes"
	"common/consts"
	dto "common/dto/common"
	adminDto "common/dto/server/admin"
	authDto "common/dto/server/auth"
	commonDto "common/dto/server/common"
	"common/entities"
	"common/infra/storage"
	"common/repositories"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/disintegration/imaging"
)

type serverUsecase struct {
	repo          repositories.ServerRepository
	configStorage storage.ConfigStorage
}

type ServerUsecase interface {
	DeviceInit(ctx context.Context, body *commonDto.DeviceInitRequest) (*entities.InitResult, *dto.ErrorResponse)
	CreateSkinImg(ctx context.Context, body adminDto.CreateSkinImgRequest) (interface{}, error)
}

func NewServerUsecase(repo repositories.ServerRepository, configStorage storage.ConfigStorage) ServerUsecase {
	return &serverUsecase{
		repo:          repo,
		configStorage: configStorage,
	}
}

func (u *serverUsecase) DeviceInit(ctx context.Context, body *commonDto.DeviceInitRequest) (*entities.InitResult, *dto.ErrorResponse) {

	// DB 조회 connectInfo(접속 url)은 관리해야할 필요있음. 최초 이후에 클라이언트가 정보가 필요할때를 대비해서.
	connectInfo, err := u.repo.GetConnectInfo(body.WorksCode)
	if err != nil {
		return nil, &dto.ErrorResponse{
			Code:    consts.E_102,
			Message: consts.E_102_MSG,
		}
	}

	// AUTH에 JWT 요청
	issuedAppToken, err := generateAppToken(body, connectInfo.ServerUrl) //  serverUrl은 이후에 .env 또는 k8s의 secrets에서 읽기
	if err != nil {
		return nil, &dto.ErrorResponse{
			Code:    consts.E_500,
			Message: consts.E_500_MSG,
		}
	}

	// 타임존, 언어, 앱 별 스킨 정보, 설정 정보 - GetConnectInfo와 합쳐서 트랜잭션 처리
	worksConfig, err := u.repo.GetWorksConfig(toWorksConfigEntity(body.WorksCode, body.Device), ctx)
	if err != nil {
		return nil, &dto.ErrorResponse{
			Code:    consts.E_500,
			Message: consts.E_500_MSG,
		}
	}

	log.Println("connectInfo:", connectInfo)
	log.Println("issuedAppToken:", issuedAppToken)
	log.Println("worksConfig ", worksConfig)

	return &entities.InitResult{
		ConnectInfo:    connectInfo,
		IssuedAppToken: issuedAppToken,
		WorksConfig:    worksConfig,
	}, nil
}

func generateAppToken(body *commonDto.DeviceInitRequest, serverUrl string) (*entities.IssuedAppToken, error) {
	// 소스 모듈화 처리하기
	data := map[string]string{
		"uuid": body.Uuid,
	}

	// JSON 변환
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	log.Println("auth service 호출! 1")

	url := "http://" + serverUrl + "/auth/sv1/generate-app-token"
	//url := domain + "/auth/v1/generate-device-token"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer serverToken") // 서버 api key

	// POST 요청 보내기
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("auth service 호출 에러 1")
		return nil, err
	}
	defer resp.Body.Close()

	// 응답 출력
	var result dto.ServerResponse[*authDto.DeviceInitAuthResponse]

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("serverReponse 파싱시 에러")
		return nil, err
	}

	return toIssuedAppTokenEntity(result.Data), nil
}

func toIssuedAppTokenEntity(dto *authDto.DeviceInitAuthResponse) *entities.IssuedAppToken {

	return &entities.IssuedAppToken{
		AppToken:     dto.AppToken,
		RefreshToken: dto.RefreshToken,
	}
}

func (r *serverUsecase) CreateSkinImg(ctx context.Context, dto adminDto.CreateSkinImgRequest) (interface{}, error) {

	log.Println("333")
	// 파일 저장
	entity, err := skinFileSaved(ctx, dto)
	if err != nil {
		log.Println("skinFileSaved error : ", err)
		return nil, err
	} else {
		log.Println("eee")
		// DB 저장
		_, err := r.repo.PutSkinFileInfo(ctx, entity)
		if err != nil {
			return nil, err
		} else {
			// 메모리 갱신
			r.configStorage.SaveConfigHash(consts.SKIN, entity.FileHash)       // skin : 파일의 hash
			r.configStorage.SaveSkinFilePath(entity.SkinType, entity.FilePath) // skinType : 파일 Path
			return nil, nil
		}
	}
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
	log.Println("444")
	// 파일의 사이즈 검증
	fileSize := dto.FileInfo.Size
	sizeCheck := checkSkinImgSize(fileSize)
	if !sizeCheck {
		return nil, consts.ErrFileSizeExceeded
	}
	log.Println("555")
	// 파일의 확장자 검증 (이미지인지 판단.)
	detectedType, err := detectContentType(dto.File)
	if err != nil {
		return nil, consts.ErrFileExtentionDetect
	}
	log.Println("666")
	if !strings.HasPrefix(detectedType, "image/") {
		return nil, consts.ErrFileExtentionInvalid
	}
	log.Println("777")
	// 1. 이미지 디코딩
	srcImage, format, err := image.Decode(dto.File)
	if err != nil {
		return nil, fmt.Errorf("이미지 디코딩 실패: %w", err)
	}
	log.Println("888")
	// 2. 리사이징 (너비 300, 비율 유지) 추후 서버 설정으로 뺄 것
	dstImage := imaging.Resize(srcImage, 300, 0, imaging.Lanczos)

	log.Println("999")
	// 3. 저장할 파일 이름 생성
	ext := strings.ToLower(filepath.Ext(dto.FileInfo.Filename))
	fileHash := getNow()
	skinType := dto.SkinType
	fileName := fmt.Sprintf("%s%s", fileHash, ext)
	filePath := filepath.Join("./skins/"+skinType, fileName) // 저장 경로 skins/skinType
	log.Println("ext : ", ext)
	log.Println("fileHash : ", fileHash)
	log.Println("fileName : ", fileName)
	log.Println("filePath : ", filePath)

	if err := os.MkdirAll("./skins/"+skinType, 0755); err != nil {
		return nil, fmt.Errorf("디렉토리 생성 실패: %w", err)
	}
	log.Println("000")
	// 5. 로컬 저장
	outFile, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("파일 생성 실패: %w", err)
	}
	log.Println("000 format : ", format)
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
	log.Println("bbb")

	if err != nil {
		return nil, fmt.Errorf("이미지 저장 실패: %w", err)
	}

	// 별도 함수로 뺄것(update시에도 사용.)
	// 1. 저장된 경로 기준 디렉토리 경로 구성
	dirPath := filepath.Join("./skins", dto.SkinType)
	// 2. 파일 목록 읽기
	files, err := os.ReadDir(dirPath)
	if err == nil {
		for _, file := range files {
			// 3. 현재 저장한 파일은 제외하고 삭제
			if !file.IsDir() && file.Name() != fileName {
				oldFilePath := filepath.Join(dirPath, file.Name())
				if err := os.Remove(oldFilePath); err != nil {
					fmt.Printf("기존 파일 삭제 실패: %s\n", oldFilePath)
				} else {
					fmt.Printf("기존 파일 삭제 성공: %s\n", oldFilePath)
				}
			}
		}
	}
	log.Println("ccc")
	// 서버 설정화 필요
	return &entities.SkinFileInfoEntity{
		FileHash: fileHash,
		SkinType: dto.SkinType,
	}, nil
}

func getNow() string {
	now := time.Now()
	formatted := now.Format(consts.YYYYMMDDHHMSS)
	return formatted
}
