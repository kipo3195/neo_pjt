package entity

type ChatRoomMemberReadDateEntity struct {
	MemberHash string `gorm:"column:member_hash"`
	ReadDate   string `gorm:"column:read_date"`
}
