package dto

type CreateOrgFileRequest struct {
	OrgCode string `json:"orgCode" validate:"required"`
}
