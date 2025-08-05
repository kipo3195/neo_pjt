package client

type configurationUsecase struct {
	repository repositories.configurationRepository
}

type ConfigurationUsecase interface {
}

func NewConfigurationUsecase(repository repositories.skinRepository) ConfigurationUsecase {
	return configurationUsecase{
		repository: repository,
	}
}
