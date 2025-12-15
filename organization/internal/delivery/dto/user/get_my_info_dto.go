package user

type GetMyInfoResponse struct {
	MyInfo       MyDetailInfo `json:"myInfo"`
	WorksOrgCode []string     `json:"worksOrgCode"`
}
