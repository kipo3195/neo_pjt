package mapper

import "file/internal/application/usecase/input"

func MakeFileUrlUploadEndInput(reqUserhash string, transactionId string) input.FileUrlUploadEndInput {

	return input.FileUrlUploadEndInput{
		ReqUserHash:   reqUserhash,
		TransactionId: transactionId,
	}
}
