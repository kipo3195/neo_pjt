package userInfoService

type UserInfoServiceDto struct {
	UserHash       string `json:"userHash" validate:"required"`
	DetailVersion  int64  `json:"detailVersion" validate:"required"`
	ProfileVersion int64  `json:"profileVersion" validate:"required"`
}
