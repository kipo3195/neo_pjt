package profile

type ProfileUploadRequest struct {
	ProfileImg     *[]byte
	ProfileImgSize int64
	ProfileImgName string
}
