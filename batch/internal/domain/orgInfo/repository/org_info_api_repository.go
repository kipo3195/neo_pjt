package repository

import "context"

type OrgInfoApiRepository interface {
	SendJsonToOrg(ctx context.Context, fileName string, zipData []byte) error
}
