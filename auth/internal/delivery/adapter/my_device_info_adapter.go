package adapter

import (
	"auth/internal/application/usecase/input"
	"auth/internal/application/usecase/output"
	"auth/internal/domain/device/entity"
)

func MakeMyDeviceInfoInput(userhash string) input.MyDeviceInfoInput {

	return input.MyDeviceInfoInput{
		UserHash: userhash,
	}
}

func MakeGetMyDeviceInfoOutput(myInfo []entity.MyDeviceInfo) []output.MyDevcieInfoOutput {

	result := make([]output.MyDevcieInfoOutput, 0)
	for _, temp := range myInfo {

		output := output.MyDevcieInfoOutput{
			Uuid:      temp.Uuid,
			Version:   temp.Version,
			ModelName: temp.ModelName,
			CreateAt:  temp.CreateAt,
		}
		result = append(result, output)
	}

	return result
}
