package server

import (
	"common/internal/consts"
	requestDTO "common/internal/domains/skin/dto/server/requestDTO"
	"common/internal/domains/skin/entities"
	repositories "common/internal/domains/skin/repositories/server"
	storage "common/internal/infra/storage"
	"context"
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

type skinUsecase struct {
	repository  repositories.SkinRepository
	skinStorage storage.SkinStorage
}

type SkinUsecase interface {
	CreateSkinImg(ctx context.Context, body requestDTO.CreateSkinImgRequest) (interface{}, error)
	GetSkinInfo() (*entities.SkinInfo, error)
}

func NewSkinUsecase(repository repositories.SkinRepository, skinStorage storage.SkinStorage) SkinUsecase {
	return &skinUsecase{
		repository:  repository,
		skinStorage: skinStorage,
	}
}

func (r *skinUsecase) GetSkinInfo() (*entities.SkinInfo, error) {

	skinHash, err := r.skinStorage.GetSkinHash()
	if err != nil {
		log.Println("[GetSkinInfo] skinHash invalid")
		return nil, err
	}

	skinFileInfos, err := r.skinStorage.GetAllSkinFiles()
	if err != nil {
		log.Println("[GetSkinInfo] skinFileInfos invalid")
		return nil, err
	}
	return &entities.SkinInfo{
		SkinHash:      skinHash,
		SkinFileInfos: toSkinFileInfos(skinFileInfos),
	}, nil
}

func toSkinFileInfos(skinFiles []map[string]string) []entities.SkinFileInfoEntity {
	var skinFileInfos []entities.SkinFileInfoEntity
	for _, file := range skinFiles {
		skinFileInfos = append(skinFileInfos, entities.SkinFileInfoEntity{
			SkinType: file["skinType"],
			FileHash: file["hash"],
		})
	}
	return skinFileInfos
}

func (r *skinUsecase) CreateSkinImg(ctx context.Context, dto requestDTO.CreateSkinImgRequest) (interface{}, error) {

	log.Println("333")
	// 파일 저장
	entity, err := skinFileSaved(ctx, dto)
	if err != nil {
		log.Println("skinFileSaved error : ", err)
		return nil, err
	} else {
		log.Println("eee")
		// DB 저장
		_, err := r.repository.PutSkinFileInfo(ctx, entity)
		if err != nil {
			return nil, err
		} else {
			// 메모리 갱신
			r.skinStorage.SaveSkinHash(consts.SKIN, entity.FileHash)         // skin : 파일의 hash
			r.skinStorage.SaveSkinFilePath(entity.SkinType, entity.FilePath) // skinType : 파일 Path
			return nil, nil
		}
	}
}

func skinFileSaved(ctx context.Context, dto requestDTO.CreateSkinImgRequest) (*entities.SkinFileInfoEntity, error) {
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
