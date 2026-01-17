package usecase

import (
	"admin/internal/application/usecase/input"
	"admin/internal/application/usecase/output"
	"admin/internal/consts"
	"admin/internal/delivery/adapter"
	"admin/internal/domain/serviceUser/entity"
	"admin/internal/domain/serviceUser/repository"
	"admin/internal/util"
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"

	"github.com/nats-io/nats.go"
)

type serviceUserUsecase struct {
	repo      repository.ServiceUserRepository
	connector *nats.Conn
}

type ServiceUserUsecase interface {
	RegistServiecUser(ctx context.Context, input input.RegistServiceUserInput) (output.RegistServiceUserOutput, error)
	PublishServiceUser(ctx context.Context, input []input.PublishServiceUserInput) error
}

func NewServiceUserUsecase(repo repository.ServiceUserRepository, connector *nats.Conn) ServiceUserUsecase {
	return &serviceUserUsecase{
		repo:      repo,
		connector: connector,
	}
}

func (r *serviceUserUsecase) RegistServiecUser(ctx context.Context, input input.RegistServiceUserInput) (output.RegistServiceUserOutput, error) {

	en := entity.MakeRegistServiceUserEntity(input.Org, input.UserId)

	serviceUsers := make([]entity.ServiceUserEntity, 0)

	for _, e := range en.UserId {

		hash, err := generateUserHash()

		// userHash 생성
		if err != nil {
			log.Println("[RegistServiecUser] generateUserHash error. err : ", err)
			return output.RegistServiceUserOutput{}, err
		}

		// 관리자를 통한 일괄 사용자 인증정보 등록 처리를 위한
		// salt 생성
		salt, err := generateSalt()
		if err != nil {
			log.Println("[RegistServiecUser] generateSalt error. err : ", err)
			return output.RegistServiceUserOutput{}, err
		}

		temp := entity.ServiceUserEntity{
			UserId:   e,
			UserHash: hash,
			Salt:     salt,
			UserAuth: input.UserAuth,
		}

		serviceUsers = append(serviceUsers, temp)
	}

	/* DB 저장 */
	result, err := r.repo.PutServiceUser(ctx, en.Org, serviceUsers)
	if err != nil {
		log.Println("[RegistServiecUser] db save error")
		return output.RegistServiceUserOutput{}, err
	}

	o := adapter.MakeRegistServiceUserOutput(result)

	return o, nil
}

func generateUserHash() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func generateSalt() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (r *serviceUserUsecase) PublishServiceUser(ctx context.Context, input []input.PublishServiceUserInput) error {

	/* service users 등록 사용자 생성*/
	eventUsers := make([]entity.ServiceUsersEventUserEntity, 0)
	for _, v := range input {

		temp := entity.ServiceUsersEventUserEntity{
			Org:      v.Org,
			UserId:   v.UserId,
			UserHash: v.UserHash,
		}

		eventUsers = append(eventUsers, temp)
	}

	publishData := entity.ServiceUsersPublishEntity{
		EventUsers: eventUsers,
	}

	data, err := util.EntityMarshal(publishData)
	if err != nil {
		log.Println(err)
		return err
	}

	/* service users 등록 사용자 전파 */
	// 20261117 모든 서비스에 던지는 구조이므로 모든 서비스의 하나의 인스턴스에 응답확인 필요시 JetStream을 사용해야 할 수 있음.
	err = r.connector.Publish("users.registered", data)
	if err != nil {
		log.Println("[PublishServiceUser] NATS Sending error")
		return consts.ErrPublishToMessageBrokerError
	}

	log.Println("[PublishServiceUser] NATS Publish end.")

	return nil
}
