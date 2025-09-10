package entity

type SkinInfo struct {
	SkinHash      string               `json:"skinHash"`
	SkinFileInfos []SkinFileInfoEntity `json:"skinFileInfos"`
}
