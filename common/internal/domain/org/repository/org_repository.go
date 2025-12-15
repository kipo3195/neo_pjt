package repository

type OrgRepository interface {
	GetOrgCode() ([]string, error)
}
