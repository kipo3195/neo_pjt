package adapter

import "user/internal/application/usecase/input"

func MakeRegistUserDetailBatchInput(file *[]byte, fileName string, orgCode string) input.RegistUserDetailBatchInput {

	return input.RegistUserDetailBatchInput{
		File:     file,
		FileName: fileName,
		OrgCode:  orgCode,
	}
}
