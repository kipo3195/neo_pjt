package storage

type orgInfoBatchStorage struct {
}

type OrgInfoBatchStorage interface {
}

func NewOrgInfoBatchStorage() OrgInfoBatchStorage {

	return &orgInfoBatchStorage{}
}
