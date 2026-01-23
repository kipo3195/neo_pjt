package input

type CreateFileUrlInput struct {
	ReqUserHash string
	EventType   string
	Org         string
	Files       []FileInfo
}
