package output

import "user/internal/domain/profile/entity"

type GetProfileMsgOutput struct {
	ProfileMsg []entity.GetProfileMsgResultEntity
}
