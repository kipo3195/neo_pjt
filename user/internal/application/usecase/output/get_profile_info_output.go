package output

import "user/internal/domain/profile/entity"

type GetProfileInfoOutput struct {
	ResultMap map[string]entity.GetProfileInfoResultEntity
}
