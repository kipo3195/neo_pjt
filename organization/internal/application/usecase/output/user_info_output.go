package output

import (
	"org/internal/domain/user/entity"
)

// usecase에서 도메인 호출 당연히 가능.
func MakeUserInfoOutput(en []entity.MyInfoEntity) []MyInfoOutput {
	var output []MyInfoOutput
	for i := 0; i < len(en); i++ {
		username := UsernameOutput{
			Def: en[i].Username.Ko, // 수정 필요
			Ko:  en[i].Username.Ko,
			En:  en[i].Username.En,
			Jp:  en[i].Username.Jp,
			Zh:  en[i].Username.Zh,
			Ru:  en[i].Username.Ru,
			Vi:  en[i].Username.Vi,
		}

		deptInfo := makeDeptInfoOutput(en[i].DeptInfo)

		output = append(output, MyInfoOutput{
			UserHash:     en[i].UserHash,
			UserPhoneNum: en[i].UserPhoneNum,
			Username:     username,
			UserEmail:    en[i].UserEmail,
			DeptInfo:     deptInfo,
			ProfileUrl:   en[i].ProfileUrl,
			ProfileMsg:   en[i].ProfileMsg,
		})
	}
	return output
}
