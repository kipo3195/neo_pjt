package output

import (
	"org/internal/domain/user/entity"
)

type MyInfoOutput struct {
	UserHash  string           `json:"userHash"`
	Username  UsernameOutput   `json:"userName"`
	MyOrgCode string           `json:"myOrgCode"`
	DeptInfo  []DeptInfoOutput `json:"deptInfo"`
}

func MakeMyInfoOutput(entity entity.MyInfoEntity) MyInfoOutput {

	username := UsernameOutput{
		Def: entity.Username.Ko, // 수정 필요
		Ko:  entity.Username.Ko,
		En:  entity.Username.En,
		Jp:  entity.Username.Jp,
		Zh:  entity.Username.Zh,
		Ru:  entity.Username.Ru,
		Vi:  entity.Username.Vi,
	}

	deptInfo := makeDeptInfoOutput(entity.DeptInfo)

	return MyInfoOutput{
		UserHash:  entity.UserHash,
		Username:  username,
		DeptInfo:  deptInfo,
		MyOrgCode: entity.MyOrg,
	}
}

func makeDeptInfoOutput(deptInfos []entity.DeptInfoEntity) []DeptInfoOutput {

	var deptInfoOutput []DeptInfoOutput

	for _, deptInfo := range deptInfos {
		roleNameOutput := RoleNameOutput{
			KoLang: deptInfo.RoleName.KoLang,
			EnLang: deptInfo.RoleName.EnLang,
			JpLang: deptInfo.RoleName.JpLang,
			ZhLang: deptInfo.RoleName.ZhLang,
		}

		positionNameOutput := PositionNameOutput{
			KoLang: deptInfo.PositionName.KoLang,
			EnLang: deptInfo.PositionName.EnLang,
			JpLang: deptInfo.PositionName.JpLang,
			ZhLang: deptInfo.PositionName.ZhLang,
		}

		deptNameOutput := DeptNameOutput{
			DefLang: deptInfo.DeptName.DefLang,
			KoLang:  deptInfo.DeptName.KoLang,
			EnLang:  deptInfo.DeptName.EnLang,
			JpLang:  deptInfo.DeptName.JpLang,
			ZhLang:  deptInfo.DeptName.ZhLang,
			ViLang:  deptInfo.DeptName.ViLang,
			RuLang:  deptInfo.DeptName.RuLang,
		}

		deptInfoOutput = append(deptInfoOutput, DeptInfoOutput{
			DeptOrg:      deptInfo.DeptOrg,
			DeptCode:     deptInfo.DeptCode,
			DeptName:     deptNameOutput,
			Header:       deptInfo.Header,
			Description:  deptInfo.Description,
			RoleName:     roleNameOutput,
			PositionName: positionNameOutput,
		})

	}

	return deptInfoOutput
}
