package userInfoService

type GetUserInfoServiceRequest struct {
	//UserHashs []string `json:"userHashs"`
	ReqUsers []UserInfoServiceDto `json:"reqUsers"`
}
