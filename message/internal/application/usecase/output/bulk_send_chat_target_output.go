package output

type BulkSendChatTargetOutput struct {
	RoomKey    string
	RoomType   string
	SecretFlag string
	MemberHash []string
}
