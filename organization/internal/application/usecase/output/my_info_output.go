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

		deptInfoOutput = append(deptInfoOutput, DeptInfoOutput{
			DeptOrg:     deptInfo.DeptOrg,
			DeptCode:    deptInfo.DeptCode,
			DefLang:     deptInfo.DefLang,
			KoLang:      deptInfo.KoLang,
			EnLang:      deptInfo.EnLang,
			JpLang:      deptInfo.JpLang,
			ZhLang:      deptInfo.ZhLang,
			ViLang:      deptInfo.ViLang,
			RuLang:      deptInfo.RuLang,
			Header:      deptInfo.Header,
			Description: deptInfo.Description,
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
