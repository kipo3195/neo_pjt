package adapter

import "org/internal/application/usecase/input"

func MakeWorksDeptInput(org string, maxDepth int, deptCount int) input.CreateWorksDeptInput {
	return input.CreateWorksDeptInput{
		Org:       org,
		MaxDepth:  maxDepth,
		DeptCount: deptCount,
	}
}
