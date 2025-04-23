package models

type Rule struct {
	Tenant    string `json:"tenant"`
	Rule      string `json:"rule"`
	Timestamp string `json:"timestamp"`
}

func (Rule) TableName() string {
	return "rule"
}
