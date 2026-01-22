package repository

import (
	"context"
	"log"
	"message/internal/consts"
	"message/internal/domain/chatRoom/entity"
	"message/internal/domain/chatRoom/repository"
	"message/internal/infrastructure/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type chatRoomRepositoryImpl struct {
	db *gorm.DB
}

func NewChatRoomRepository(db *gorm.DB) repository.ChatRoomRepository {
	return &chatRoomRepositoryImpl{
		db: db,
	}
}

func ChatRoomMigrate(db *gorm.DB) {
	db.AutoMigrate(&model.ChatRoom{})
	db.AutoMigrate(&model.ChatRoomDetail{})
	db.AutoMigrate(&model.ChatRoomMember{})
	db.AutoMigrate(&model.ChatRoomOwner{})
}

func (r *chatRoomRepositoryImpl) PutChatRoom(ctx context.Context, en entity.CreateChatRoomEntity) error {
	// 트랜잭션 시작
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 실패 시 롤백 보장
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()

	// 방 정보
	if err := tx.Create(&model.ChatRoom{
		RoomKey:  en.ChatRoomEntity.RoomKey,
		RoomType: en.ChatRoomEntity.RoomType,
	}).Error; err != nil {
		// 감지가 안되므로 일단 주석 처리
		// if errors.Is(err, gorm.ErrDuplicatedKey) {
		// 	log.Println("[PutChatRoom] room - Duplicate key : ", roomEntity.RoomKey)
		// 	return consts.ErrRoomKeyAlreadyExist
		// }
		return err
	}

	// 방 상세 정보
	if err := tx.Create(&model.ChatRoomDetail{
		RoomKey:         en.ChatRoomEntity.RoomKey,
		RoomTitle:       en.ChatRoomEntity.Title,
		RoomSecretFlag:  en.ChatRoomEntity.SecretFlag,
		RoomSecret:      en.ChatRoomEntity.Secret,
		RoomDescription: en.ChatRoomEntity.Description,
		RoomWorksCode:   en.ChatRoomEntity.WorksCode,
		RoomCreateDate:  en.ChatRoomEntity.RegDate,
		RoomUpdateDate:  en.ChatRoomEntity.RegDate,
		RoomCreateUser:  en.ChatRoomEntity.CreateUserHash,
	}).Error; err != nil {
		// 감지가 안되므로 일단 주석 처리
		// if errors.Is(err, gorm.ErrDuplicatedKey) {
		// 	log.Println("[PutChatRoom] room detail - Duplicate key : ", roomEntity.RoomKey)
		// 	return consts.ErrRoomKeyAlreadyExist
		// }
		return err
	}

	// 방 참여자 정보
	if len(en.ChatRoomMemberEntity) > 0 {
		// Entity를 Model로 변환
		memberModels := make([]model.ChatRoomMember, len(en.ChatRoomMemberEntity))
		for i, member := range en.ChatRoomMemberEntity {
			memberModels[i] = model.ChatRoomMember{
				RoomKey:         en.ChatRoomEntity.RoomKey,
				MemberHash:      member.MemberHash,
				MemberFirstDate: en.ChatRoomEntity.RegDate, // 최초 입장 시간
				MemberDate:      en.ChatRoomEntity.RegDate, // 입장 시간
				MemberWorksCode: member.MemberWorksCode,
			}
		}

		// lause.OnConflict를 사용하여 중복 키 발생 시 아무것도 하지 않음 (DoNothing)
		if err := tx.Clauses(clause.OnConflict{
			// 중복을 판단할 컬럼 (RoomKey와 MemberHash의 복합 키)
			Columns: []clause.Column{{Name: "room_key"}, {Name: "member_hash"}},

			// 충돌 발생 시 업데이트할 컬럼과 값 지정
			DoUpdates: clause.Assignments(map[string]interface{}{
				// 컬럼 이름: 업데이트할 값 (함수 인자로 받은 regDate)
				"member_date": en.ChatRoomEntity.RegDate, // 입장 시간
			}),
		}).Create(&memberModels).Error; err != nil {
			// OnConflict를 사용했으므로 DuplicatedKey 에러는 발생하지 않으며,
			// 다른 심각한 DB 에러가 발생했을 때만 여기로 들어옵니다.
			log.Println("[PutChatRoom] member insert/conflict failed : ", err)
			return err
		}
	}

	// 방장
	if err := tx.Create(&model.ChatRoomOwner{
		RoomKey:         en.ChatRoomEntity.RoomKey,
		OwnerHash:       en.ChatRoomEntity.CreateUserHash,
		ActiveFlag:      "Y",
		MemberWorksCode: en.ChatRoomEntity.WorksCode,
		UpdateDate:      en.ChatRoomEntity.RegDate,
	}).Error; err != nil {
		// 감지가 안되므로 일단 주석 처리
		// if errors.Is(err, gorm.ErrDuplicatedKey) {
		// 	log.Println("[PutChatRoom] room detail - Duplicate key : ", roomEntity.RoomKey)
		// 	return consts.ErrRoomKeyAlreadyExist
		// }
		return err
	}

	// 트랜잭션 종료
	if err := tx.Commit().Error; err != nil {
		log.Println("[PutChatRoom] - Commit failed")
		return err
	}
	return nil
}

func (r *chatRoomRepositoryImpl) GetChatRoomDetail(ctx context.Context, en entity.GetChatRoomDetailEntity) ([]entity.ChatRoomDetailEntity, error) {

	var result []entity.ChatRoomDetailEntity

	// 내가 참여중인 방에 한해서 조회 가능하도록 처리함 20251208
	err := r.db.Raw(`
		select 
			detail.*,
			line.line_key, line.event_type, line.cmd, line.contents, line.send_date,
			member_view.member,
			owner_view.owner,
			room.room_type,
			case when title.update_flag is null then 'N' else title.update_flag end as title_update_flag, title.my_room_title, title.update_date as title_update_date,
			unread.member_unread_count as unread_count, unread.member_unread_count_date as unread_count_date, unread.member_read_date as last_read_date
		from chat_room_member as member 
		join chat_room as room 
			on member.room_key = room.room_key and member_hash = ?
		join chat_room_detail as detail 
			on member.room_key = detail.room_key
		left join (select max(line_key) as line_key, room_key, event_type, cmd, contents, send_date from chat_line_event group by room_key) as line
			on member.room_key = line.room_key
		left join chat_room_owner as owner
			on member.room_key = owner.room_key and owner.active_flag = 'Y'
		left join (select group_concat(DISTINCT member_hash separator ',') as member, room_key from chat_room_member group by room_key) as member_view
			on member.room_key = member_view.room_key
		left join (select group_concat(DISTINCT owner_hash separator ',') as owner, room_key from chat_room_owner group by room_key) as owner_view
			on member.room_key = owner_view.room_key
		left join chat_room_title as title 
			on member.room_key = title.room_key and user_hash = ?
		left join (select room_key, member_unread_count, member_read_date, member_unread_count_date from chat_room_member where member_hash = ? ) as unread
			on member.room_key = unread.room_key
		where 
			member.room_key IN (?) and room.room_type = ? order by line.line_key desc, detail.room_create_date desc`,
		en.ReqUserHash, en.ReqUserHash, en.ReqUserHash, en.RoomKey, en.RoomType).Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return result, err
}

func (r *chatRoomRepositoryImpl) GetChatRoomList(ctx context.Context, en entity.GetChatRoomListEntity) ([]entity.ChatRoomDetailEntity, error) {

	var result []entity.ChatRoomDetailEntity

	// 내가 참여중인 방에 한해서 조회 가능하도록 처리함 20251208
	err := r.db.Raw(`
		select 
			detail.*,
			line.line_key, line.event_type, line.cmd, line.contents, line.send_date,
			member_view.member,
			owner_view.owner,
			room.room_type,
			case when title.update_flag is null then 'N' else title.update_flag end as title_update_flag, title.my_room_title, title.update_date as title_update_date,
			unread.member_unread_count as unread_count, unread.member_unread_count_date as unread_count_date, unread.member_read_date as last_read_date
		from chat_room_member as member 
		join chat_room as room 
			on member.room_key = room.room_key and member_hash = ?
		join chat_room_detail as detail 
			on member.room_key = detail.room_key
		left join (select max(line_key) as line_key, room_key, event_type, cmd, contents, send_date from chat_line_event group by room_key) as line
			on member.room_key = line.room_key
		left join chat_room_owner as owner
			on member.room_key = owner.room_key and owner.active_flag = 'Y'
		left join (select group_concat(DISTINCT member_hash separator ',') as member, room_key from chat_room_member group by room_key) as member_view
			on member.room_key = member_view.room_key
		left join (select group_concat(DISTINCT owner_hash separator ',') as owner, room_key from chat_room_owner group by room_key) as owner_view
			on member.room_key = owner_view.room_key
		left join chat_room_title as title 
			on member.room_key = title.room_key and user_hash = ?
		left join (select room_key, member_unread_count, member_read_date, member_unread_count_date from chat_room_member where member_hash = ? ) as unread
			on member.room_key = unread.room_key
		where 
			room.room_type = ? order by line.line_key desc, detail.room_create_date desc limit ?`, en.ReqUserHash, en.ReqUserHash, en.ReqUserHash, en.RoomType, en.ReqCount).Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return result, err

}

func (r *chatRoomRepositoryImpl) GetChatRoomUpdateDate(ctx context.Context, en entity.GetChatRoomUpdateDateEntity) ([]entity.ChatRoomUpdateDateEntity, error) {

	var result []entity.ChatRoomUpdateDateEntity

	var err error

	if en.Type == "A" {
		err = r.db.Raw(
			`select 
				room.room_key, room.room_type,
				detail.room_update_date as detail,
				line.send_date as line,
				member_view.member_date as member, 
				member_read_view.member_unread_count_date as unread,
				owner_view.update_date as owner,
				title.update_date as title
			from chat_room_member as member 
			join chat_room as room 
				on member.room_key = room.room_key and member_hash = ?
			join chat_room_detail as detail 
				on member.room_key = detail.room_key
			left join (select max(line_key), room_key, send_date from chat_line_event group by room_key) as line
				on member.room_key = line.room_key
			left join chat_room_owner as owner
				on member.room_key = owner.room_key and owner.active_flag = 'Y'
			left join (select max(member_date) as member_date, room_key from chat_room_member group by room_key) as member_view
				on member.room_key = member_view.room_key
			left join chat_room_member as member_read_view
				on member.room_key = member_read_view.room_key and member_read_view.member_hash = ?
			left join (select max(update_date) as update_date, room_key from chat_room_owner group by room_key) as owner_view
				on member.room_key = owner_view.room_key
			left join chat_room_title as title 
				on member.room_key = title.room_key and user_hash = ?
			order by line.send_date desc, detail.room_create_date desc `, en.ReqUserHash, en.ReqUserHash, en.ReqUserHash).Scan(&result).Error
	} else if en.Type == "S" {
		err = r.db.Raw(
			`select 
				room.room_key, room.room_type,
				detail.room_update_date as detail,
				line.send_date as line,
				member_view.member_date as member,
				member_read_view.member_unread_count_date as unread,
				owner_view.update_date as owner,
				title.update_date as title
			from chat_room_member as member 
			join chat_room as room 
				on member.room_key = room.room_key and member_hash = ?
			join chat_room_detail as detail 
				on member.room_key = detail.room_key
			left join (select max(line_key), room_key, send_date from chat_line_event group by room_key) as line
				on member.room_key = line.room_key
			left join chat_room_owner as owner
				on member.room_key = owner.room_key and owner.active_flag = 'Y'
			left join (select max(member_date) as member_date, room_key from chat_room_member group by room_key) as member_view
				on member.room_key = member_view.room_key
			left join chat_room_member as member_read_view
				on member.room_key = member_read_view.room_key and member_read_view.member_hash = ?
			left join (select max(update_date) as update_date, room_key from chat_room_owner group by room_key) as owner_view
				on member.room_key = owner_view.room_key
			left join chat_room_title as title 
				on member.room_key = title.room_key and user_hash = ?
			where detail.room_update_date >= ? or line.send_date >= ? or member_view.member_date >= ? or owner_view.update_date >= ? or title.update_date >= ?
			order by line.send_date desc, detail.room_create_date desc`,
			en.ReqUserHash, en.ReqUserHash, en.ReqUserHash,
			en.Date, en.Date, en.Date, en.Date, en.Date,
		).Scan(&result).Error
	} else {
		return nil, consts.ErrRoomUpdateDateTypeError
	}

	if err != nil {
		log.Println("[GetChatRoomUpdateDate] DB error :", err)
		return nil, err
	}

	return result, nil
}

func (r *chatRoomRepositoryImpl) GetChatRoomMemberReadDate(ctx context.Context, en entity.GetChatRoomMemberReadDateEntity) ([]entity.ChatRoomMemberReadDateEntity, error) {

	var result []entity.ChatRoomMemberReadDateEntity

	err := r.db.Raw(
		`select 
			member.member_hash, case when member.member_read_date is null or member.member_read_date = '' then 'N' else  member.member_read_date end as read_date
		from chat_room_member as member 
		join service_users as su 
			on member.member_hash = su.user_hash and member.member_state = '1' and su.use_yn = 'Y'
		where member.room_key = ?`,
		en.RoomKey,
	).Scan(&result).Error

	if err != nil {
		log.Println("[GetChatRoomMemberReadDate] DB error :", err)
		return nil, err
	}

	return result, nil
}

func (r *chatRoomRepositoryImpl) GetChatRoomMy(ctx context.Context, en entity.GetChatRoomMyEntity) ([]entity.ChatRoomMyEntity, error) {

	var result []entity.ChatRoomMyEntity

	err := r.db.Raw(
		`select 
			member.room_key, member.member_read_date as read_date,  member.member_unread_count as unread_count,
			room.room_type
		from chat_room_member as member join service_users as su 
			on member.member_hash = su.user_hash
		left join chat_room_detail as detail
			on member.room_key = detail.room_key
		left join chat_room as room 
			on member.room_key = room.room_key
		where 
			member.member_hash = ? 
		and detail.room_state ='1' 
		and member_state ='1'`,
		en.ReqUserHash,
	).Scan(&result).Error

	if err != nil {
		log.Println("[GetChatRoomMy] DB error :", err)
		return nil, err
	}

	return result, nil
}
