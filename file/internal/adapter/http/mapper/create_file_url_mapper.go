package mapper

import (
	"file/internal/adapter/http/dto/fileUrl"
	"file/internal/application/usecase/input"
)

func MakeCreateFileUrlInput(reqUserHash string, eventType string, org string, fileInfo []fileUrl.FileInfoDto) input.CreateFileUrlInput {

	files := make([]input.FileInfo, 0)

	for _, f := range fileInfo {

		temp := input.FileInfo{
			FileId:   f.FileId,
			FileName: f.FileName,
			FileSize: f.FileSize,
			FileExt:  f.FileExt,
		}

		files = append(files, temp)
	}

	return input.CreateFileUrlInput{
		ReqUserHash: reqUserHash,
		EventType:   eventType,
		Org:         org,
		Files:       files,
	}

}
