package model

type ChatMessage struct {
	Linekey  string `json:"line_key"`
	Roomkey  string `json:"room_key"`
	Sender   string `json:"sender"`
	Message  string `json:"message"`
	SendDate string `json:"send_date"`
}

func (ChatMessage) TableName() string {
	return "chat_message"
}
