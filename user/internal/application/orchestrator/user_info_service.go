package orchestrator

import (
	"context"
	"user/internal/application/usecase"
	"user/internal/application/usecase/input"
	"user/internal/application/usecase/output"
)

type UserInfoService struct {
	Profile    usecase.ProfileUsecase
	UserDetail usecase.UserDetailUsecase
}

func NewUserInfoService(p usecase.ProfileUsecase, u usecase.UserDetailUsecase) *UserInfoService {
	return &UserInfoService{
		Profile:    p,
		UserDetail: u,
	}
}

func (r *UserInfoService) GetMyInfo(ctx context.Context, userHash []string) (output.UserInfoOutput, error) {

	detailOutput, err := r.UserDetail.GetMyDetailInfo(ctx, userHash)
	if err != nil {
		return output.UserInfoOutput{}, err
	}

	profileOutput, err := r.Profile.GetMyProfileInfo(ctx, userHash)
	if err != nil {
		return output.UserInfoOutput{}, err
	}

	output := output.UserInfoOutput{
		UserProfile: profileOutput,
		UserDetail:  detailOutput,
	}

	return output, nil

}

func (r *UserInfoService) GetUserInfo(ctx context.Context, detailInput []input.GetUserDetailInfoInput, profileInput []input.GetUserProfileInfoInput) (output.UserInfoOutput, error) {

	detailResult, err := r.UserDetail.GetUserDetailInfo(ctx, detailInput)

	if err != nil {
		return output.UserInfoOutput{}, err
	}

	// 사용자 정보

	detailOutput := make([]output.UserDetailOutput, 0)

	for _, r := range detailResult {

		d := output.UserDetailOutput{
			UserHash:     r.UserHash,
			UserEmail:    r.UserEmail,
			UserPhoneNum: r.UserPhoneNum,
		}

		detailOutput = append(detailOutput, d)
	}

	// 프로필 정보

	profileResult, err := r.Profile.GetProfileInfo(ctx, profileInput)

	if err != nil {
		return output.UserInfoOutput{}, err
	}

	profileOutput := make([]output.UserProfileOutput, 0)

	for _, r := range profileResult {

		p := output.UserProfileOutput{
			UserHash:    r.UserHash,
			ProfileHash: r.ProfileHash,
			ProfileMsg:  r.ProfileMsg,
		}

		profileOutput = append(profileOutput, p)
	}

	output := output.UserInfoOutput{
		UserProfile: profileOutput,
		UserDetail:  detailOutput,
	}

	return output, nil

}
