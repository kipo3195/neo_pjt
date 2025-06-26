package loader

import "common/usecases"

type CommonLoader struct {
	CommonUsecase usecases.CommonUsecase
}

func NewCommonLoader(commonUsecase usecases.CommonUsecase) *CommonLoader {
	return &CommonLoader{
		CommonUsecase: commonUsecase,
	}
}

func (r *CommonLoader) RunAll() error {

	err := r.CommonUsecase.InitConfigHash()
	if err != nil {
		return err
	}

	return nil
}
