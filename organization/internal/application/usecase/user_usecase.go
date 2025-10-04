package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	mrand "math/rand"
	"org/internal/application/usecase/input"
	"org/internal/application/usecase/output"
	"org/internal/domain/user/entity"
	"org/internal/domain/user/repository"
	"strings"
	"time"
)

type userUsecase struct {
	repository repository.UserRepository
}

type UserUsecase interface {
	GetMyInfo(ctx context.Context, input input.MyInfoInput) (output.MyInfoOutput, error)
	CreateServiceUser(ctx context.Context, input input.CreateServiceUserInput) error
	CreateUserDetail(ctx context.Context, input input.CreateUserDetailInput) error
}

func NewUserUsecase(repository repository.UserRepository) UserUsecase {
	return &userUsecase{
		repository: repository,
	}
}

func (r *userUsecase) GetMyInfo(ctx context.Context, input input.MyInfoInput) (output.MyInfoOutput, error) {

	entity := entity.MakeMyInfoHashEntity(input.MyHash)

	myInfo, err := r.repository.GetMyInfo(ctx, entity)

	if err != nil {
		return output.MyInfoOutput{}, err
	}
	output := output.MakeMyInfoOutput(myInfo)
	return output, nil
}

func (r *userUsecase) CreateServiceUser(ctx context.Context, input input.CreateServiceUserInput) error {

	var entities []entity.ServiceUserEntity
	for i := 0; i < input.UserCount; i++ {

		hash, err := generateUserHash()
		if err != nil {
			return fmt.Errorf("failed to generate user hash: %w", err)
		}

		userID := fmt.Sprintf("%s%04d", input.Keyword, i+1)

		entities = append(entities, entity.ServiceUserEntity{
			UserHash: hash,
			UserId:   userID,
		})
	}

	return r.repository.CreateServiceUser(ctx, entities)
}

func generateUserHash() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (r *userUsecase) CreateUserDetail(ctx context.Context, input input.CreateUserDetailInput) error {

	// like 검색으로 사용자를 조회함.
	entities, err := r.repository.GetServiceUsers(ctx, input.Keyword)
	if err != nil {
		return err
	}

	log.Println("조회된 사용자의 수 : ", len(entities))

	// type을 확인
	if input.Type == "email" {
		for i := 0; i < len(entities); i++ {
			email, _ := generateRandomEmail()
			entities[i].UserEmail = email
		}
	} else if input.Type == "phoneNum" {
		for i := 0; i < len(entities); i++ {
			entities[i].UserPhoneNum = generateRandomPhoneNum()
		}
	} else if input.Type == "all" {
		log.Print("all type is not supported.")
	}

	err = r.repository.CreateUserDetail(ctx, entities)

	if err != nil {
		return err
	}

	return nil
}

func generateRandomEmail() (string, error) {
	// 4바이트(8자리 hex) 랜덤 문자열 생성
	b := make([]byte, 5)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	randomPart := hex.EncodeToString(b)

	domains := []string{
		"example.com",
		"test.com",
		"sample.org",
		"mail.net",
		"demo.co.kr",
		"naver.com",
		"google.com",
	}
	domain := domains[int(b[0])%len(domains)]

	email := fmt.Sprintf("%s@%s", randomPart, domain)
	return strings.ToLower(email), nil
}

func generateRandomPhoneNum() string {
	r := mrand.New(mrand.NewSource(time.Now().UnixNano()))

	prefix := "010"
	mid := r.Intn(10000)  // 0~9999
	last := r.Intn(10000) // 0~9999

	return fmt.Sprintf("%s-%04d-%04d", prefix, mid, last)
}
