package entity

type RegistOrgBatch struct {
	OrgFile     *[]byte
	OrgFileName string
}

func MakeRegistOrgBatchEntity(orgFile *[]byte, orgFileName string) RegistOrgBatch {

	return RegistOrgBatch{
		OrgFile:     orgFile,
		OrgFileName: orgFileName,
	}
}
