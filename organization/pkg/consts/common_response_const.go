package consts

const (
	SUCCESS = "success"
	FAIL    = "fail"
	ERROR   = "error"

	BAD_REQUEST        = 400
	SERVER_ERROR       = 500
	StatusUnauthorized = 401

	// ERROR
	E_101     = "E101"
	E_101_MSG = "Type invalid."

	E_102     = "E102"
	E_102_MSG = "Sql exception."

	// body 데이터가 잘못되었을때
	E_103     = "E103"
	E_103_MSG = "Data invalid."

	E_104     = "E104"
	E_104_MSG = "Header invalid."

	E_105     = "E105"
	E_105_MSG = "Token format invalid."

	E_106     = "E106"
	E_106_MSG = "Unauthorized."

	E_107     = "E107"
	E_107_MSG = "Token expired."

	E_108     = "E108"
	E_108_MSG = "Query string is invalid."

	E_500     = "E500"
	E_500_MSG = "Server error."
)
