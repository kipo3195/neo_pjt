package orchestrator

import "admin/internal/application/usecase"

type ServiceUserAuthRegisterService struct {
	ServiceUser      usecase.ServiceUserUsecase
	UserAuthRegister usecase.UserAuthRegisterUsecase
}

func NewServiceUserAuthRegisterService(serviceUser usecase.ServiceUserUsecase, userAuthRegister usecase.UserAuthRegisterUsecase) *ServiceUserAuthRegisterService {

	return &ServiceUserAuthRegisterService{
		ServiceUser:      serviceUser,
		UserAuthRegister: userAuthRegister,
	}

}
