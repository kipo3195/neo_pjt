package entity

type OrgEntity struct {
	// RootDept any `json:"rootDept"`
	// OrgTree  any `json:"orgTree"`
	// 20251214 여기 수정함 조직도 json
	RootDept []OrgInfo     `json:"rootDept"`
	OrgTree  []OrgTreeInfo `json:"orgTree"`
}
