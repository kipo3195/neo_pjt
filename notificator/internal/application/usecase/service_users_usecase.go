package usecase

import (
	"context"
	"log"
	"notificator/internal/application/usecase/input"
	"notificator/internal/domain/serviceUsers/entity"
	"notificator/internal/domain/serviceUsers/repository"
)

type serviceUsersUsecase struct {
	repo repository.ServiceUsersRepository
}

type ServiceUsersUsecase interface {
	RecvRegistServiceUsersMessage(ctx context.Context, input input.ServiceUsersMessageInput) error
}

func NewServiceUsersUsecase(repo repository.ServiceUsersRepository) ServiceUsersUsecase {
	return &serviceUsersUsecase{
		repo: repo,
	}
}

func (r *serviceUsersUsecase) RecvRegistServiceUsersMessage(ctx context.Context, in input.ServiceUsersMessageInput) error {

	log.Println("[RecvRegistServiceUsersMessage] users : ", in.EventUsers)

	en := make([]entity.RegisterServiceUsersEntity, 0)

	for _, i := range in.EventUsers {

		temp := entity.RegisterServiceUsersEntity{
			Org:      i.Org,
			UserHash: i.UserHash,
			UserId:   i.UserId,
		}
		en = append(en, temp)
	}

	err := r.repo.PutServiceUser(ctx, en)

	if err != nil {
		log.Println("[RecvRegistServiceUsersMessage] db save error")
		return err
	}

	return nil
}
