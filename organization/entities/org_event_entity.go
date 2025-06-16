package entities

type OrgEventEntity struct {
	Seq        int    `json:"seq"`
	EventType  string `json:"eventType"`
	Kind       string `json:"kind"`
	OrgCode    string `json:"orgCode"`
	UpdateHash string `json:"updateHash"`
}
