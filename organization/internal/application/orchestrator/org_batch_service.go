package orchestrator

import (
	"context"
	"org/internal/application/usecase"
	"org/internal/application/usecase/input"
)

type OrgBatchService struct {
	Department usecase.DepartmentUsecase
	Org        usecase.OrgUsecase
	User       usecase.UserUsecase
}

func NewOrgBatchService(department usecase.DepartmentUsecase, org usecase.OrgUsecase, user usecase.UserUsecase) *OrgBatchService {

	return &OrgBatchService{
		Department: department,
		Org:        org,
		User:       user,
	}
}

func (r *OrgBatchService) RegistOrgBatch(ctx context.Context, input input.RegistOrgBatchInput) error {

	r.Org.RegistOrgBatch(ctx, input)

	// 그룹
	// select * from works_dept;
	// select * from works_dept_multi_lang;

	// 사용자
	// select * from user_detail;
	// select * from works_dept_user;
	// select * from works_user_multi_lang;

	return nil
}
