package fileUrl

type CreateFileUrlResponse struct {
	TransactionId string           `json:"transactionId"`
	FileUrlInfo   []FileUrlInfoDto `json:"fileUrlInfo"`
}
