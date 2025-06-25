package dto

type GetConfigHash struct {
	SkinHash   string `json:"skinHash"`
	ConfigHash string `json:"configHash"`
	Device     string `json:"device"`
}
