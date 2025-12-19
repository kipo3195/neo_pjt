package adapter

func MakeGetMyInfoInput(userHash string) []string {

	result := make([]string, 0)

	result = append(result, userHash)

	return result
}
