package output

import (
	"org/internal/domain/user/entity"
)

type MyInfoOutput struct {
	UserHash     string           `json:"userHash"`
	UserPhoneNum string           `json:"userPhoneNum"`
	Username     UsernameOutput   `json:"userName"`
	OrgCodes     []string         `json:"orgCodes"`
	ProfileUrl   string           `json:"profileUrl"`
	ProfileMsg   string           `json:"profileMsg"`
	DeptInfo     []DeptInfoOutput `json:"deptInfo"`
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
		UserHash:     entity.UserHash,
		UserPhoneNum: entity.UserPhoneNum,
		Username:     username,
		DeptInfo:     deptInfo,
		ProfileUrl:   entity.ProfileUrl,
		ProfileMsg:   entity.ProfileMsg,
	}
}

func makeDeptInfoOutput(deptInfos []entity.DeptInfoEntity) []DeptInfoOutput {

	var deptInfoOutput []DeptInfoOutput

	for _, deptInfo := range deptInfos {
		deptInfoOutput = append(deptInfoOutput, DeptInfoOutput{
			DeptOrg:  deptInfo.DeptOrg,
			DeptCode: deptInfo.DeptCode,
			DefLang:  deptInfo.DefLang,
			KoLang:   deptInfo.KoLang,
			EnLang:   deptInfo.EnLang,
			JpLang:   deptInfo.JpLang,
			ZhLang:   deptInfo.ZhLang,
			ViLang:   deptInfo.ViLang,
			RuLang:   deptInfo.RuLang,
			Header:   deptInfo.Header,
		})
	}

	return deptInfoOutput
}
