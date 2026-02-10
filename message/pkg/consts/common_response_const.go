package consts

const (
	SUCCESS = "success"
	FAIL    = "fail"
	ERROR   = "error"
)
const (
	BAD_REQUEST     = 400
	SERVER_ERROR    = 500
	UNAUTHORIZED    = 401
	GATEWAY_TIMEOUT = 504
)
const (
	// ERROR Code
	E_101     = "E101"
	E_101_MSG = "Type invalid."

	E_102     = "E102"
	E_102_MSG = "Sql exception."

	E_103     = "E103"
	E_103_MSG = "Data invalid."

	E_104     = "E104"
	E_104_MSG = "Header invalid."

	E_105     = "E105"
	E_105_MSG = "Token format invalid"

	E_106     = "E106"
	E_106_MSG = "Unauthorized."

	E_107     = "E107"
	E_107_MSG = "Token expired."

	E_108     = "E108"
	E_108_MSG = "Check client parameters."

	E_109     = "E109"
	E_109_MSG = "There is no user mapped to the app token."

	E_110     = "E110"
	E_110_MSG = "User hash invalid. check user auth info."

	E_500     = "E500"
	E_500_MSG = "Server error."

	E_504     = "E504"
	E_504_MSG = "Gateway timeout."
)
