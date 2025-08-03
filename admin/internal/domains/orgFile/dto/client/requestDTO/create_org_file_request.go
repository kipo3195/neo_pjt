package requestDTO

type CreateOrgFileRequestDTO struct {
	Body   CreateOrgFileRequestBody
	Header CreateOrgFileRequestHeader
}

type CreateOrgFileRequestBody struct {
	OrgCode []string `json:"orgCode" validate:"required"`
}

type CreateOrgFileRequestHeader struct {
}
