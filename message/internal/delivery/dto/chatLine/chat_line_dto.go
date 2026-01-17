package chatLine

type ChatLineDto struct {
	LineKey   string `json:"lineKey"`
	Cmd       int    `json:"cmd"`
	EventType string `json:"eventType"`
	Contents  string `json:"contents"`
	SendDate  string `json:"sendDate"`
}
