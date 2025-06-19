package repositories

import (
	"context"
	"org/entities"
	"org/models"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

type UserRepository interface {
	GetMyInfo(ctx context.Context, entity entities.GetMyInfoEntity) (entities.MyInfoEntity, error)
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetMyInfo(ctx context.Context, entity entities.GetMyInfoEntity) (entities.MyInfoEntity, error) {
	var myInfo models.MyInfo

	err := r.db.Raw(
		`SELECT 
			su.user_hash,
			ud.user_phone_num,
			wuml.kr_lang,
			wuml.en_lang,
			wuml.cn_lang,
			wuml.jp_lang
		FROM service_users AS su
		JOIN user_detail AS ud 
			ON su.user_hash = ud.user_hash
		JOIN works_user_multi_lang AS wuml 
			ON su.user_hash = wuml.user_hash
		WHERE su.user_hash = ? AND su.use_yn ='Y'`,

		entity.MyHash).Scan(&myInfo).Error

	if err != nil {
		return entities.MyInfoEntity{}, err
	}

	return toMyInfoEntity(myInfo), nil
}

func toMyInfoEntity(myInfo models.MyInfo) entities.MyInfoEntity {

	// 사용자 명 다국어 처리
	userName := entities.UsernameEntity{
		Kr: myInfo.KrLang,
		En: myInfo.EnLang,
		Cn: myInfo.CnLang,
		Jp: myInfo.JpLang,
	}

	return entities.MyInfoEntity{
		UserHash:     myInfo.UserHash,
		UserPhoneNum: myInfo.UserPhoneNum,
		Username:     userName,
	}
}
