package input

type CreateOrgFileInput struct {
	OrgCode []string
}

func MakeCreateOrgFileInput(orgCode []string) CreateOrgFileInput {

	return CreateOrgFileInput{
		OrgCode: orgCode,
	}
}
