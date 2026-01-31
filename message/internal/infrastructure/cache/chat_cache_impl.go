package cache

import (
	"context"
	"encoding/json"
	"log"
	"message/internal/consts"
	"message/internal/domain/chat/cache"
	"message/internal/domain/chat/entity"

	"github.com/redis/go-redis/v9"
)

type chatCacheImpl struct {
	cacheClient *redis.ClusterClient
}

func NewChatCache(cacheClient *redis.ClusterClient) cache.ChatCache {
	return &chatCacheImpl{
		cacheClient: cacheClient,
	}
}

func (r *chatCacheImpl) GetFileEntity(ctx context.Context, transactionId string) ([]*entity.ChatFileEntity, error) {

	key := consts.RedisFileUploadEndPrefix + transactionId

	val, err := r.cacheClient.Get(ctx, key).Result()

	if err == redis.Nil || val == "" {
		log.Println("[GetFileEntity] nil.")
		return nil, consts.ErrCacheResultNotFound // 데이터가 없는 경우 (Cache Miss)
	} else if err != nil {
		log.Println("[GetFileEntity] err :", err)
		return nil, err // 기타 에러
	}

	// 2. JSON 문자열을 실제 객체 구조체로 변환 (Unmarshal)
	var entity []*entity.ChatFileEntity
	if err := json.Unmarshal([]byte(val), &entity); err != nil {
		log.Println("[GetFileEntity] json parser error")
		return nil, err
	}

	log.Println("entity : ", entity)
	return entity, nil
}
