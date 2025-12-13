package storage

type orgInfoStorage struct {
}

type OrgInfoStorage interface {
}

func NewOrgInfoStorage() OrgInfoStorage {

	return &orgInfoStorage{}
}
