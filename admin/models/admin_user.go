package models

type Admin struct {
	Id       string `json:"id"`
	Password string `json:"password"`
}

func (Admin) TableName() string {
	return "admin_users"
}
