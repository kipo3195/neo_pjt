package output

import "user/internal/domain/userDetail/entity"

type GetUserDetailInfoOutput struct {
	UserInfos []entity.UserDetailInfoEntity
}
