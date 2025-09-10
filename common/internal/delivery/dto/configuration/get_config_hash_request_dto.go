package configuration

type GetConfigHashRequestDTO struct {
	Body   GetConfigHashRequestBody
	Header GetConfigHashRequestHeader
}

type GetConfigHashRequestBody struct {
	SkinHash   string `json:"skinHash"`
	ConfigHash string `json:"configHash"`
	Device     string `json:"device"`
}

type GetConfigHashRequestHeader struct {
}
