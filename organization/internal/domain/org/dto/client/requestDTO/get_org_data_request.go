package requestDTO

type GetOrgDataRequest struct {
	OrgCode string `json:"orgCode"`
	Type    string `json:"type"`
	OrgHash string `json:"orgHash"`
}
