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
	//CreateUserDetail(ctx context.Context, input input.CreateUserDetailInput) error
	CreateUserMultiLang(ctx context.Context, input input.CreateUserMultiLangInput) error
	GetServiceUsers(ctx context.Context, org string) ([]output.ServiceUsersOutput, error)
	GetUserInfo(ctx context.Context, input input.GetUserInfoInput) ([]output.MyInfoOutput, error)
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

func (r *userUsecase) GetUserInfo(ctx context.Context, input input.GetUserInfoInput) ([]output.MyInfoOutput, error) {

	entity := entity.MakeUserInfoEntity(input.UserHashs)
	userInfo, err := r.repository.GetUserInfo(ctx, entity)

	if err != nil {
		return nil, err
	}

	// 같은 usecase의 output 호출 가능
	return output.MakeUserInfoOutput(userInfo), nil
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

func (r *userUsecase) GetServiceUsers(ctx context.Context, org string) ([]output.ServiceUsersOutput, error) {

	// like 검색으로 사용자를 조회함.
	entities, err := r.repository.GetServiceUsers(ctx, org)

	if err != nil {
		return nil, err
	}

	temp := make([]output.ServiceUsersOutput, len(entities))

	for i := 0; i < len(entities); i++ {
		temp[i] = output.ServiceUsersOutput{
			UserHash: entities[i].UserHash,
		}
	}

	return temp, nil

}

// func (r *userUsecase) CreateUserDetail(ctx context.Context, input input.CreateUserDetailInput) error {

// 	// like 검색으로 사용자를 조회함.
// 	entities, err := r.repository.GetServiceUsers(ctx, input.Keyword)
// 	if err != nil {
// 		return err
// 	}

// 	log.Println("조회된 사용자의 수 : ", len(entities))

// 	// type을 확인
// 	for i := 0; i < len(entities); i++ {
// 		email, _ := generateRandomEmail()
// 		entities[i].UserEmail = email
// 		entities[i].UserPhoneNum = generateRandomPhoneNum()
// 	}

// 	err = r.repository.CreateUserDetail(ctx, entities)

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

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

func (r *userUsecase) CreateUserMultiLang(ctx context.Context, input input.CreateUserMultiLangInput) error {

	// 🔹 한글 성/이름 파츠 (약 300개)
	koreanParts := []string{
		"김", "이", "박", "최", "정", "강", "조", "윤", "장", "임", "한", "오", "서", "신", "권", "황", "안", "송", "류", "홍",
		"전", "고", "문", "손", "양", "배", "백", "허", "남", "심", "노", "하", "주", "구", "성", "우", "나", "진", "민", "유",
		"임", "채", "곽", "엄", "변", "염", "여", "방", "조", "위", "표", "명", "기", "라", "마", "차", "길", "표", "지", "선",
		"하", "도", "연", "동", "현", "성", "재", "민", "수", "훈", "석", "영", "진", "태", "지", "준", "희", "윤", "원", "호",
		"예", "우", "은", "나", "해", "소", "유", "현", "슬", "하", "연", "진", "민", "경", "예", "지", "하", "영", "보", "라",
		"수", "아", "윤", "채", "별", "은", "율", "담", "온", "나", "리", "서", "정", "희", "해", "솔", "선", "하", "연", "유",
		"승", "동", "정", "철", "명", "지", "우", "태", "규", "건", "혁", "상", "재", "성", "범", "환", "민", "기", "형", "도",
		"진", "욱", "현", "찬", "선", "연", "호", "빈", "훈", "완", "슬", "은", "주", "미", "진", "혜", "나", "유", "연", "하",
		"현", "서", "아", "라", "린", "솔", "별", "윤", "유", "은", "선", "혜", "라", "민", "지", "서", "윤", "나", "채", "리",
	}

	// 🔹 영어 이름 파츠
	englishParts := []string{
		"Kim", "Lee", "Park", "Choi", "Jung", "Kang", "Yoon", "Han", "Lim", "Cho", "Aiden", "Noah", "Ethan", "Liam", "Mason",
		"Jacob", "William", "James", "Benjamin", "Lucas", "Logan", "Alexander", "Henry", "Owen", "Sebastian", "Daniel", "Matthew",
		"Samuel", "David", "Joseph", "John", "Carter", "Wyatt", "Jack", "Luke", "Jayden", "Gabriel", "Isaac", "Dylan", "Anthony",
		"Emily", "Olivia", "Sophia", "Isabella", "Mia", "Charlotte", "Amelia", "Harper", "Ella", "Avery", "Evelyn", "Abigail",
		"Scarlett", "Grace", "Lily", "Chloe", "Zoey", "Hannah", "Nora", "Layla", "Addison", "Ellie", "Lillian", "Stella", "Natalie",
		"Leah", "Audrey", "Bella", "Skylar", "Hazel", "Violet", "Aurora", "Claire", "Lucy", "Anna", "Samantha", "Sadie", "Caroline",
	}

	// 🔹 베트남 이름 파츠
	vietnameseParts := []string{
		"Nguyen", "Tran", "Le", "Pham", "Huynh", "Phan", "Vu", "Dang", "Bui", "Do", "Hoang", "Ngo", "Duong", "Ly", "Vo",
		"Thao", "Anh", "Bao", "Linh", "My", "Trang", "Huong", "Minh", "Tuan", "Son", "Nam", "Hien", "Phuong", "Quang",
		"Thanh", "Khanh", "Duy", "Lan", "Ha", "Loan", "Diep", "Nga", "Hanh", "Tam", "Phuc", "Thien", "Long", "An", "Thu",
	}

	// 🔹 러시아 이름 파츠
	russianParts := []string{
		"Иван", "Алексей", "Сергей", "Дмитрий", "Михаил", "Николай", "Андрей", "Владимир", "Павел", "Антон",
		"Екатерина", "Ольга", "Татьяна", "Наталья", "Мария", "Елена", "Анна", "Ирина", "Виктория", "Алиса",
	}

	// 🔹 일본 이름 파츠
	japaneseParts := []string{
		"キム", "イ", "パク", "チェ", "チョン", "カン", "ユン", "チョ", "ハン", "イム",
		"さくら", "はるか", "みさき", "りん", "ひな", "ゆい", "あい", "まな", "りこ", "かな",
		"たけし", "けん", "しんじ", "ゆうた", "だいすけ", "たかし", "けい", "そうた", "ゆうき", "まこと",
	}

	// 🔹 중국어(간체) 이름 파츠
	chineseParts := []string{
		"金", "李", "王", "张", "刘", "陈", "杨", "赵", "黄", "周", "吴", "徐", "孙", "胡", "朱",
		"何", "郭", "林", "罗", "郑", "冯", "邓", "曹", "潘", "曾", "丁", "姜", "程", "傅", "魏",
		"伟", "芳", "娜", "敏", "静", "秀英", "丽", "强", "磊", "军", "洋", "勇", "艳", "杰", "娟",
		"涛", "明", "超", "霞", "平", "刚", "桂英", "丹", "萍", "玲", "鑫", "桂兰", "亮", "俊", "鹏",
	}

	serviceUsers, err := r.repository.GetServiceUsers(ctx, input.Keyword)

	if err != nil {
		log.Println("[CreateUserMultiLang] service user 조회 에러")
		return err
	}

	count := len(serviceUsers)

	koreanNames := generateNames(koreanParts, count)
	englishNames := generateNames(englishParts, count)
	vietnameseNames := generateNames(vietnameseParts, count)
	russianNames := generateNames(russianParts, count)
	japaneseNames := generateNames(japaneseParts, count)
	chineseNames := generateNames(chineseParts, count)

	for i := 0; i < count; i++ {

		e := entity.MakeUserMultilangEntity(serviceUsers[i].UserHash, koreanNames[i], englishNames[i], vietnameseNames[i], russianNames[i], japaneseNames[i], chineseNames[i])

		err = r.repository.CreateUserMultiLang(ctx, e)
		if err != nil {
			fmt.Printf("serviceUsers : %s insert error.", serviceUsers[i].UserHash)
		}

	}

	return nil
}

// 🔹 이름 생성 함수
func generateNames(parts []string, count int) []string {

	r := mrand.New(mrand.NewSource(time.Now().UnixNano()))

	names := make([]string, count)
	for i := 0; i < count; i++ {
		nameLen := r.Intn(3) + 2 // 2~4글자
		var sb strings.Builder
		for j := 0; j < nameLen; j++ {
			sb.WriteString(parts[r.Intn(len(parts))])
		}
		names[i] = sb.String()
	}
	return names
}
