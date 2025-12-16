package adapter

import "org/internal/application/usecase/input"

func MakeRegistOrgBatchInput(orgFile *[]byte, fileName string, orgCode string) input.RegistOrgBatchInput {
	return input.RegistOrgBatchInput{
		OrgFile:     orgFile,
		OrgFileName: fileName,
		OrgCode:     orgCode,
	}
}
