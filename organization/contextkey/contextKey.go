package contextkey

type key string

const UserHashKey key = "userHash"

// Go 공식 문서에서도 context key는 string 대신 고유한 타입을 정의해서 사용하라고 권장, authMiddleware에서 사용
