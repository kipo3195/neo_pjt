package repository

import (
	"context"
	"log"
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
}

func (r *chatRoomRepositoryImpl) PutChatRoom(ctx context.Context, memberEntity []entity.CreateChatRoomMemberEntity, roomEntity entity.ChatRoomEntity) error {
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
		RoomKey:  roomEntity.RoomKey,
		RoomType: roomEntity.RoomType,
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
		RoomKey:         roomEntity.RoomKey,
		RoomTitle:       roomEntity.Title,
		RoomSecretFlag:  roomEntity.SecretFlag,
		RoomSecret:      roomEntity.Secret,
		RoomDescription: roomEntity.Description,
		RoomWorksCode:   roomEntity.WorksCode,
		RoomCreateDate:  roomEntity.RegDate,
		RoomUpdateDate:  roomEntity.RegDate,
		RoomCreateUser:  roomEntity.CreateUserHash,
	}).Error; err != nil {
		// 감지가 안되므로 일단 주석 처리
		// if errors.Is(err, gorm.ErrDuplicatedKey) {
		// 	log.Println("[PutChatRoom] room detail - Duplicate key : ", roomEntity.RoomKey)
		// 	return consts.ErrRoomKeyAlreadyExist
		// }
		return err
	}

	// 방 참여자 정보
	if len(memberEntity) > 0 {
		// Entity를 Model로 변환
		memberModels := make([]model.ChatRoomMember, len(memberEntity))
		for i, member := range memberEntity {
			memberModels[i] = model.ChatRoomMember{
				RoomKey:         roomEntity.RoomKey,
				MemberHash:      member.MemberHash,
				MemberFirstDate: roomEntity.RegDate, // 최초 입장 시간
				MemberDate:      roomEntity.RegDate, // 입장 시간
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
				"member_date": roomEntity.RegDate, // 입장 시간
			}),
		}).Create(&memberModels).Error; err != nil {
			// OnConflict를 사용했으므로 DuplicatedKey 에러는 발생하지 않으며,
			// 다른 심각한 DB 에러가 발생했을 때만 여기로 들어옵니다.
			log.Println("[PutChatRoom] member insert/conflict failed : ", err)
			return err
		}
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
			AA.*,
			group_concat(DISTINCT BB.member_hash separator ',') as member 
		from (
			select 
				detail.*,
				line.room_hash
			from chat_room_member as member 
			left join chat_room as room 
				on member.room_key = room.room_key
			left join chat_room_detail as detail 
				on room.room_key= detail.room_key
			left join (select max(line_key) as room_hash, room_key from chat_line_event group by room_key) as line 
				on room.room_key = line.room_key
			where 
				member.room_key IN (?) and (member_hash = ? and room.room_type = ?)) as AA 
		join chat_room_member as BB on AA.room_key = BB.room_key
		group by AA.room_key `, en.RoomKey, en.ReqUserHash, en.RoomType).Scan(&result).Error

	if err != nil {
		log.Println("[GetChatRoomDetail] DB error :", err)
		return nil, err
	}

	return result, err
}

func (r *chatRoomRepositoryImpl) GetChatRoomList(ctx context.Context, en entity.GetChatRoomListEntity) ([]entity.ChatRoomDetailEntity, error) {

	var result []entity.ChatRoomDetailEntity

	// 내가 참여중인 방에 한해서 조회 가능하도록 처리함 20251208
	err := r.db.Raw(`
		select 
			AA.*,
			group_concat(DISTINCT BB.member_hash separator ',') as member 
		from (
			select 
				detail.*,
				line.room_hash
			from chat_room_member as member 
			left join chat_room as room 
				on member.room_key = room.room_key
			left join chat_room_detail as detail 
				on room.room_key= detail.room_key
			left join (select max(line_key) as room_hash, room_key from chat_line_event group by room_key) as line 
				on room.room_key = line.room_key
			where 
				member_hash = ? and room.room_type = ? limit ?) as AA 
		join chat_room_member as BB on AA.room_key = BB.room_key
		group by AA.room_key`, en.ReqUserHash, en.RoomType, en.ReqCount).Scan(&result).Error

	if err != nil {
		log.Println("[GetChatRoomList] DB error :", err)
		return nil, err
	}

	return result, err

}
