package entity

type RegistUserDetailBatchEntity struct {
	File     *[]byte
	FileName string
	OrgCode  string
}

func MakeRegistUserDetailBatchEntity(file *[]byte, fileName string, orgCode string) RegistUserDetailBatchEntity {

	return RegistUserDetailBatchEntity{
		File:     file,
		FileName: fileName,
		OrgCode:  orgCode,
	}
}
