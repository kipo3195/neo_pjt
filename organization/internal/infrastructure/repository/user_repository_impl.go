package repository

import (
	"context"
	"org/internal/domain/user/entity"
	"org/internal/domain/user/repository"
	"org/internal/infrastructure/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepositoryImpl{
		db: db,
	}
}

func UserMigrate(db *gorm.DB) {
	db.AutoMigrate(&model.UserDetail{})
	db.AutoMigrate(&model.UserGrade{})
	db.AutoMigrate(&model.UserProfile{})
}

func (r *userRepositoryImpl) GetMyInfo(ctx context.Context, en entity.MyInfoHashEntity) (entity.MyInfoEntity, error) {
	var myDetailInfo model.MyDetailInfo
	var myDeptInfo []model.DeptInfo

	// 트랜잭션 시작
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return entity.MyInfoEntity{}, tx.Error
	}

	// 첫 번째 쿼리: 사용자 상세 정보
	err := tx.Raw(
		`SELECT 
			su.user_hash,
			ud.user_phone_num,
			ud.user_email,
			wuml.ko_lang,
			wuml.en_lang,
			wuml.zh_lang,
			wuml.jp_lang,
			wuml.ru_lang,
			wuml.vi_lang,
			up.profile_url,
			up.profile_msg
		FROM service_users AS su
		JOIN user_detail AS ud 
			ON su.user_hash = ud.user_hash
		JOIN works_user_multi_lang AS wuml 
			ON su.user_hash = wuml.user_hash
		LEFT JOIN user_profile AS up
			ON su.user_hash = up.user_hash	
		WHERE su.user_id = ? AND su.use_yn = 'Y'`,
		en.MyHash).Scan(&myDetailInfo).Error
	if err != nil {
		tx.Rollback()
		return entity.MyInfoEntity{}, err
	}

	// 두 번째 쿼리: 부서 정보
	err = tx.Raw(
		`SELECT 
			wdml.dept_org,
			wdml.dept_code,
			wdml.def_lang,
			wdml.ko_lang,
			wdml.en_lang,
			wdml.jp_lang,
			wdml.zh_lang,
			wdml.ru_lang,
			wdml.vi_lang,
			wd.header
		FROM works_dept AS wd 
		JOIN works_dept_multi_lang AS wdml 
			ON wd.dept_code = wdml.dept_code 
		JOIN (
			SELECT wdu.dept_code FROM service_users AS su 
			JOIN works_dept_user AS wdu 
				ON su.user_hash = wdu.user_hash 
			WHERE su.use_yn = 'Y' AND su.user_id = ?) AS a 
			ON wdml.dept_code = a.dept_code`,
		en.MyHash).Scan(&myDeptInfo).Error
	if err != nil {
		tx.Rollback()
		return entity.MyInfoEntity{}, err
	}

	// 트랜잭션 커밋
	if err := tx.Commit().Error; err != nil {
		return entity.MyInfoEntity{}, err
	}

	// 매핑 및 반환
	return toMyInfoEntity(myDetailInfo, myDeptInfo), nil
}

func toMyInfoEntity(myDetailInfo model.MyDetailInfo, myDeptInfo []model.DeptInfo) entity.MyInfoEntity {

	// 사용자 명 다국어 처리
	userName := entity.UserNameEntity{
		Def: myDetailInfo.KoLang, // 수정 필요
		Ko:  myDetailInfo.KoLang,
		En:  myDetailInfo.EnLang,
		Zh:  myDetailInfo.ZhLang,
		Jp:  myDetailInfo.JpLang,
		Ru:  myDetailInfo.RuLang,
		Vi:  myDetailInfo.ViLang,
	}

	// 부서 정보 파싱
	deptInfoEntity := toDeptInfoEntity(myDeptInfo)

	// 내 정보
	return entity.MyInfoEntity{
		UserHash:     myDetailInfo.UserHash,
		UserPhoneNum: myDetailInfo.UserPhoneNum,
		Username:     userName,
		ProfileUrl:   myDetailInfo.ProfileUrl,
		ProfileMsg:   myDetailInfo.ProfileMsg,
		DeptInfo:     deptInfoEntity,
	}
}

func toDeptInfoEntity(myDeptInfo []model.DeptInfo) []entity.DeptInfoEntity {

	var deptEntity []entity.DeptInfoEntity

	for _, dept := range myDeptInfo {
		deptEntity = append(deptEntity, entity.DeptInfoEntity{
			DeptOrg:  dept.DeptOrg,
			DeptCode: dept.DeptOrg,
			DefLang:  dept.DefLang,
			KoLang:   dept.KoLang,
			EnLang:   dept.EnLang,
			JpLang:   dept.JpLang,
			ZhLang:   dept.ZhLang,
			ViLang:   dept.ViLang,
			RuLang:   dept.RuLang,
			Header:   dept.Header,
		})
	}
	return deptEntity
}

func (r *userRepositoryImpl) CreateServiceUser(ctx context.Context, entities []entity.ServiceUserEntity) error {

	// entity → model 변환
	var models []model.ServiceUsers
	for _, e := range entities {
		models = append(models, model.ServiceUsers{
			UserHash: e.UserHash,
			UserId:   e.UserId,
			UseYn:    "Y", // 기본값 설정
		})
	}

	tx := r.db.WithContext(ctx)
	if err := tx.Create(&models).Error; err != nil {
		return err
	}

	return nil
}

func (r *userRepositoryImpl) GetServiceUsers(ctx context.Context, keyword string) ([]entity.UserDetailEntity, error) {

	var models []model.ServiceUsers

	tx := r.db.WithContext(ctx)
	if err := tx.
		Where("user_id LIKE ? AND use_yn = 'Y'", keyword+"%").
		Find(&models).Error; err != nil {
		return nil, err
	}

	// model → entity 변환
	var entities []entity.UserDetailEntity
	for _, m := range models {
		entities = append(entities, entity.UserDetailEntity{
			UserHash: m.UserHash,
		})
	}

	return entities, nil
}

func (r *userRepositoryImpl) CreateUserDetail(ctx context.Context, entities []entity.UserDetailEntity) error {

	// entity → model 변환
	var models []model.UserDetail
	for _, e := range entities {
		models = append(models, model.UserDetail{
			UserHash:     e.UserHash,
			UserPhoneNum: e.UserPhoneNum,
			UserEmail:    e.UserEmail, // entity에 이메일이 있다면
		})
	}

	tx := r.db.WithContext(ctx)

	// Upsert: 이미 존재하면 update, 없으면 insert
	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_hash"}},                               // PK 기준
		DoUpdates: clause.AssignmentColumns([]string{"user_phone_num", "user_email"}), // update할 컬럼
	}).Create(&models).Error; err != nil {
		return err
	}

	return nil
}
