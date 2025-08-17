package requestDTO

type ConnectInfoRequest struct {
	WorksCode string `json:"worksCode"`
	Uuid      string `json:"uuid"`
	Device    string `json:"device"`
}
