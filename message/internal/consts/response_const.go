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

	// 방 update date 조회 타입 에러
	MESSAGE_F012     = "MESSAGE_F012"
	MESSAGE_F012_MSG = "Chat room update date request type invalid."

	// 방 read count 갱신 실패 (모두 읽은 상황)
	MESSAGE_F013     = "MESSAGE_F013"
	MESSAGE_F013_MSG = "Already read chat all."

	// 채팅에 파일을 첨부하여 전송시 transactionId에 등록된 파일이 존재하지 않거나 만료된 transacionId
	MESSAGE_F014     = "MESSAGE_F014"
	MESSAGE_F014_MSG = "Transaction id is not regist or expired upload file info."

	// 채팅에 파일을 첨부하여 전송한 후, 전송한 파일의 정보를 file 서비스에 호출.
	MESSAGE_F015     = "MESSAGE_F015"
	MESSAGE_F015_MSG = "Send chat file transactionId sending error to file service."
)
