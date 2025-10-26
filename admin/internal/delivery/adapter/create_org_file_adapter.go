package adapter

import "admin/internal/application/usecase/input"

func MakeCreateOrgFileInput(orgCode []string) input.CreateOrgFileInput {

	return input.CreateOrgFileInput{
		OrgCode: orgCode,
	}
}
