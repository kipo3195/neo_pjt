package orchestrator

import "org/internal/application/usecase"

type OrgUserService struct {
	Org  usecase.OrgUsecase
	User usecase.UserUsecase
}

func NewOrgUserService(org usecase.OrgUsecase, user usecase.UserUsecase) *OrgUserService {
	return &OrgUserService{
		Org:  org,
		User: user,
	}
}
