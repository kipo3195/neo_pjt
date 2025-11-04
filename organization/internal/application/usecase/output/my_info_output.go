package output

import (
	"org/internal/domain/user/entity"
)

type MyInfoOutput struct {
	UserHash string           `json:"userHash"`
	Username UsernameOutput   `json:"userName"`
	OrgCodes []string         `json:"orgCodes"`
	DeptInfo []DeptInfoOutput `json:"deptInfo"`
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

	deptInfo, orgCodes := makeDeptInfoOutput(entity.DeptInfo)

	return MyInfoOutput{
		UserHash: entity.UserHash,
		Username: username,
		DeptInfo: deptInfo,
		OrgCodes: orgCodes,
	}
}

func makeDeptInfoOutput(deptInfos []entity.DeptInfoEntity) ([]DeptInfoOutput, []string) {

	var deptInfoOutput []DeptInfoOutput

	orgCodes := make(map[string]struct{})
	var uniqueOrgs []string

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

		// DeptOrg가 orgCodes에 이미 존재하는지 체크
		if _, exists := orgCodes[deptInfo.DeptOrg]; !exists {
			// 신규 DeptOrg면 put
			orgCodes[deptInfo.DeptOrg] = struct{}{}
			uniqueOrgs = append(uniqueOrgs, deptInfo.DeptOrg)
		}
	}

	return deptInfoOutput, uniqueOrgs
}
