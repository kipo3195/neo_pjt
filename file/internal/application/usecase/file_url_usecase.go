package usecase

import (
	"context"
	"file/internal/application/usecase/input"
	"file/internal/application/usecase/output"
	"file/internal/consts"
	"file/internal/domain/fileUrl/entity"
	"file/internal/domain/fileUrl/repository"
	"file/internal/domain/logger"
	"file/pkg/util"
	"log"
)

type fileUrlUsecase struct {
	repo        repository.FileUrlRepository
	storageRepo repository.FileUrlStorageRepository
	logger      logger.Logger
}

type FileUrlUsecase interface {
	CreateFileUrl(ctx context.Context, input input.CreateFileUrlInput) (output.CreateFileUrlOutput, error)
	FileUrlUploadEnd(ctx context.Context, input input.FileUrlUploadEndInput) error
}

func NewFileUrlUsecase(repo repository.FileUrlRepository, storageRepo repository.FileUrlStorageRepository, logger logger.Logger) FileUrlUsecase {
	return &fileUrlUsecase{
		repo:        repo,
		storageRepo: storageRepo,
		logger:      logger,
	}
}

func (r *fileUrlUsecase) CreateFileUrl(ctx context.Context, input input.CreateFileUrlInput) (output.CreateFileUrlOutput, error) {

	fileInfoMap := make(map[string]entity.FileInfoEntity)

	for _, v := range input.Files {

		if v.FileId == "" || v.FileName == "" || v.FileExt == "" || v.FileSize <= 0 {
			// 유효성 검증
			continue
		}

		_, exists := fileInfoMap[v.FileName]
		if exists {
			// 파일 명이 존재하는지 검증
			continue
		}

		temp := entity.FileInfoEntity{
			FileId:   v.FileId,
			FileName: v.FileName,
			FileSize: v.FileSize,
			FileExt:  v.FileExt,
		}

		// 동일한 파일명이 있을 수 있으므로
		fileInfoMap[temp.FileId] = temp
	}

	// transactionId 는 ULID로 사용함
	ulidGen, err := util.NewULIDGenerator()
	transactionId := ulidGen.New()
	if err != nil {
		return output.CreateFileUrlOutput{}, consts.ErrULIDGeneratorError
	}

	// url 생성
	entity := entity.MakeCreateFileUrlEntity(input.ReqUserHash, input.Org, fileInfoMap)
	result, err := r.storageRepo.CreateFileUrl(ctx, entity)

	if err != nil {
		return output.CreateFileUrlOutput{}, err
	}

	// DB 저장
	err = r.repo.SaveCreateFileUrl(ctx, input.ReqUserHash, transactionId, result)
	if err != nil {
		r.logger.Error(ctx, "file_url_save_fail",
			"save_url", err.Error())
		return output.CreateFileUrlOutput{}, err
	}

	fileUrlOutput := make([]output.FileUrlInfo, 0)
	for _, v := range result {

		// 파일 명이 존재하는지 검증
		_, exists := fileInfoMap[v.FileId]
		if !exists {
			continue
		}

		temp := output.FileUrlInfo{
			FileId:   v.FileId,
			FileName: v.FileName,
			Url:      v.CreatedUrl,
		}

		fileUrlOutput = append(fileUrlOutput, temp)
	}

	out := output.CreateFileUrlOutput{
		TransactionId: transactionId,
		FileUrlInfo:   fileUrlOutput,
	}

	return out, nil
}

func (r *fileUrlUsecase) FileUrlUploadEnd(ctx context.Context, input input.FileUrlUploadEndInput) error {

	en := entity.MakeFileUrlUploadEndEntity(input.ReqUserHash, input.TransactionId)

	files, err := r.repo.GetFileId(ctx, en)

	if err != nil {
		r.logger.Error(ctx, "upload_file_id_select_fail",
			"transactionId", err.Error())
		return err
	}

	fileIds := make([]string, 0)
	for _, f := range files {

		result, err := r.storageRepo.CheckFileExists(ctx, f.FileId)

		if err != nil {
			log.Printf("fileId : %s invalid.. err :%s", f, err)
			return err
		}

		if !result {
			log.Printf("fileId : %s not regist storage", f)
			return err
		}

		fileIds = append(fileIds, f.FileId)
	}

	// 업로드 한 파일 upload_flag 변경
	err = r.repo.UploadFlagUpdate(ctx, en.ReqUserHash, fileIds)
	if err != nil {
		return err
	}

	// 업로드 완료 후 transaction id 저장 - message 서비스 조회용
	err = r.repo.PutUploadEndFileInfo(ctx, en.TransactionId, files)
	if err != nil {
		return err
	}
	return nil
}
