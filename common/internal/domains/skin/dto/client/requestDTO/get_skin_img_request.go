package requestDTO

type GetSkinImgRequest struct {
	SkinHash string `json:"skinHash" validate:"required"`
	SkinType string `json:"skinType" validate:"required"`
}
