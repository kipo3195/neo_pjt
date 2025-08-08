package services

type AppProfileService struct {
	validator     appValidation.appValidationUsecase
	skin          skin.skinUsecase
	configuration configuration.configurationUsecase
}

func NewAppProfileService(v appValidation.appValidationUsecase, s skin.skinUsecase, c configuration.configurationUsecase) *AppProfileService {
	return &AppProfileService{v, s, c}
}

// api 수정 필요.
func (svc *AppProfileService) GetAppProfile(appID, deviceInfo string) (map[string]interface{}, error) {
	if err := svc.validator.Validate(appID, deviceInfo); err != nil {
		return nil, err
	}
	skinData, _ := svc.skin.GetSkin(appID)
	configData, _ := svc.configuration.GetConfig(appID)

	return map[string]interface{}{
		"appID":  appID,
		"skin":   skinData,
		"config": configData,
	}, nil
}
