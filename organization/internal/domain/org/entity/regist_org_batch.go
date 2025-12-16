package entity

type RegistOrgBatch struct {
	OrgFile     *[]byte
	OrgFileName string
	OrgCode     string
}

func MakeRegistOrgBatchEntity(orgFile *[]byte, orgFileName string, orgCode string) RegistOrgBatch {

	return RegistOrgBatch{
		OrgFile:     orgFile,
		OrgFileName: orgFileName,
		OrgCode:     orgCode,
	}
}
