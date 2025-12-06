package model

import "time"

type ChatRoomDetail struct {
	RoomKey         string    `gorm:"column:room_key;type:varchar(50);primaryKey;comment:고정된 생성 형식. 동일한 참여자 생성의 경우 체크할 수 있도록 함"`
	RoomTitle       string    `gorm:"column:room_title;type:varchar(100);comment:방 제목"`
	RoomSecretFlag  string    `gorm:"column:room_secret_flag;type:varchar(1);comment:방 비공개 여부 Y, N"`
	RoomSecret      string    `gorm:"column:room_secret;type:varchar(50);comment:방 비밀번호 암호화 값"`
	RoomDescription string    `gorm:"column:room_description;type:varchar(400);comment:방 상세 정보(설명) "`
	RoomState       string    `gorm:"column:room_state;type:varchar(1);default:1;comment:방 상태 1 : 사용, 2 : 중지 (조회만 가능), 3 : 폐쇄 (입장 불가)"`
	RoomWorksCode   string    `gorm:"column:room_works_code;type:varchar(50);comment:방 생성자 기준 works code"`
	RoomCreateDate  time.Time `gorm:"column:room_create_date;not null;comment:등록 일"`
	RoomUpdateDate  time.Time `gorm:"column:room_create_date;not null;comment:수정 일"`
	RoomCreateUser  string    `gorm:"column:room_create_user;not null;comment:등록 사용자"`
}

func (ChatRoomDetail) TableName() string {
	return "chat_room_detail"
}
