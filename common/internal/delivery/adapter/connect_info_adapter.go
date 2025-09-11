package adapter

import (
	"common/internal/application/usecase/input"
	"common/internal/delivery/dto/device"
	"common/internal/delivery/dto/worksInfo"
)

type ConnectInput interface {
	ToConnectInfoInput() *input.ConnectInfoInput
}

// worksInfo
type WorksInfoWrapper struct {
	worksInfo.ConnectInfoRequest
}

func (w WorksInfoWrapper) ToConnectInfoInput() *input.ConnectInfoInput {
	return &input.ConnectInfoInput{
		WorksCode: w.WorksCode,
		Uuid:      w.Uuid,
		Device:    w.Device,
	}
}

// device
type DeviceWrapper struct {
	device.DeviceRequestBody
}

func (d DeviceWrapper) ToConnectInfoInput() *input.ConnectInfoInput {
	return &input.ConnectInfoInput{
		WorksCode: d.WorksCode,
		Uuid:      d.Uuid,
		Device:    d.Device,
	}
}

// 공통 함수
func MakeConnectInfoInput(ci ConnectInput) *input.ConnectInfoInput {
	return ci.ToConnectInfoInput()
}
