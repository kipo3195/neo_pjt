package input

type CreateWorksDeptInput struct {
	Org       string
	MaxDepth  int
	DeptCount int
}

func MakeWorksDeptInput(org string, maxDepth int, deptCount int) CreateWorksDeptInput {
	return CreateWorksDeptInput{
		Org:       org,
		MaxDepth:  maxDepth,
		DeptCount: deptCount,
	}
}
