package entity

type FileUrlUploadEndEntity struct {
	ReqUserHash   string
	TransactionId string
}

func MakeFileUrlUploadEndEntity(reqUserHash string, transactionId string) FileUrlUploadEndEntity {
	return FileUrlUploadEndEntity{
		ReqUserHash:   reqUserHash,
		TransactionId: transactionId,
	}
}
