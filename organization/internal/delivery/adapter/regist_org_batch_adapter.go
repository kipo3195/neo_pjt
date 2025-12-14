package adapter

import "org/internal/application/usecase/input"

func MakeRegistOrgBatchInput(orgFile *[]byte, fileName string) input.RegistOrgBatchInput {
	return input.RegistOrgBatchInput{
		OrgFile:     orgFile,
		OrgFileName: fileName,
	}
}
