package entities

type OrgEventEntity struct {
	Id         string `json:"id"`
	EventType  string `json:"eventType"`
	Kind       string `json:"kind"`
	OrgCode    string `json:"orgCode"`
	UpdateHash string `json:"updateHash"`
}
