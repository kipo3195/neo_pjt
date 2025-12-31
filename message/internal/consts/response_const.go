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

	// 방 제목 변경 실패 - 조건에 일치하는 데이터가 없는 경우
	MESSAGE_F008     = "MESSAGE_F008"
	MESSAGE_F008_MSG = "Chat room title update error."

	// 방 제목 삭제 실패 - 조건에 일치하는 데이터가 없는 경우
	MESSAGE_F009     = "MESSAGE_F009"
	MESSAGE_F009_MSG = "Chat room title delete error."

	// 방 제목 변경, 삭제 실패 - 조건에 일치하는 데이터가 없거나 참여중인 방이 아님.
	MESSAGE_F010     = "MESSAGE_F010"
	MESSAGE_F010_MSG = "Chat room that do not exist or that you are not participating in."

	// 방 제목 변경, 삭제 실패 - 요청 룸 타입과 서버의 룸 타입이 일치하지 않음.
	MESSAGE_F011     = "MESSAGE_F011"
	MESSAGE_F011_MSG = "Chat room type mismatch."

	// 방 생성시 생성자가 참여자에 포함되어 있지 않음.
	MESSAGE_F012     = "MESSAGE_F012"
	MESSAGE_F012_MSG = "Room creator is not included in member"
)
