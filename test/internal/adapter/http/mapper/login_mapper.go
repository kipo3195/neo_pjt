package mapper

import "test/internal/application/usecase/input"

func MakePutLoginInput(userIdPrefix string, connectionCount int) input.PutLoginInput {

	return input.PutLoginInput{
		UserIdPrefix:    userIdPrefix,
		ConnectionCount: connectionCount,
	}
}
