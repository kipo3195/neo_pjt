package consts

const (

	// 채팅 발송 실패
	MESSAGE_F001     = "MESSAGE_F001"
	MESSAGE_F001_MSG = "Publish to message broker failed."

	// OTP 정보 없음
	MESSAGE_F002     = "MESSAGE_F002"
	MESSAGE_F002_MSG = "OTP information not found."

	// 방 생성시 참여자에 생성자 정보 누락
	MESSAGE_F003     = "MESSAGE_F003"
	MESSAGE_F003_MSG = "Chat room member invalid."

	// 방 생성시 룸 키 중복
	MESSAGE_F004     = "MESSAGE_F004"
	MESSAGE_F004_MSG = "Chat room key is already regist."

	// 방 생성시 룸 키 중복
	MESSAGE_F005     = "MESSAGE_F005"
	MESSAGE_F005_MSG = "Chat room type invalid."

	// 방 생성시 시크릿 구분값 에러
	MESSAGE_F006     = "MESSAGE_F006"
	MESSAGE_F006_MSG = "Chat room secret flag invalid."

	// 방 생성시 시크릿 구분값이 Y임에도 시크릿 정보 존재 하지않음.
	MESSAGE_F007     = "MESSAGE_F007"
	MESSAGE_F007_MSG = "Chat room secret empty."
)
