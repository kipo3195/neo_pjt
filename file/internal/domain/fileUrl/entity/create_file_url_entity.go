package entity

type CreateFileUrlEntity struct {
	ReqUserHash string
	Org         string
	FileInfoMap map[string]FileInfoEntity
}

func MakeCreateFileUrlEntity(reqUserHash string, org string, fileInfoMap map[string]FileInfoEntity) CreateFileUrlEntity {

	return CreateFileUrlEntity{
		ReqUserHash: reqUserHash,
		Org:         org,
		FileInfoMap: fileInfoMap,
	}
}
