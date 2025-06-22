package dto

type SvDeviceInitRequest struct {
	WorksCode string `json:"worksCode"`
	Uuid      string `json:"uuid"`
	Device    string `json:"device"`
}
