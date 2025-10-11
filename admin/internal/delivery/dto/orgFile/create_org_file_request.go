package orgFile

type CreateOrgFileRequest struct {
	OrgCode []string `json:"orgCode" validate:"required"`
}
