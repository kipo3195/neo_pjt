package entity

type CreateOrgFileEntity struct {
	OrgCode []string `json:"orgCode"`
}

func MakeCreateOrgFileEntity(orgCode []string) CreateOrgFileEntity {
	return CreateOrgFileEntity{
		OrgCode: orgCode,
	}
}
