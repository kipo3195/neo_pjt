package util

import "strings"

func SplitAndMakeSet(s string, separator string) []string {
	parts := strings.Split(s, separator)
	seen := make(map[string]struct{})
	result := make([]string, 0)

	for _, p := range parts {
		p = strings.TrimSpace(p) // 공백 제거
		if p == "" {
			continue
		}
		if _, exists := seen[p]; !exists {
			seen[p] = struct{}{}
			result = append(result, p)
		}
	}
	return result
}
