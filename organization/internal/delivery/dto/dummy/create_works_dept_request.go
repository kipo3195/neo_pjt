package dummy

type CreateWorksDeptRequest struct {
	Org       string `json:"org"`
	MaxDepth  int    `json:"maxDepth"`
	DeptCount int    `json:"deptCount"`
}
