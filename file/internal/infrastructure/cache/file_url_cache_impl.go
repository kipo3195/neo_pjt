package cache

import (
	"context"
	"encoding/json"
	"file/internal/consts"
	"file/internal/domain/fileUrl/cache"
	"file/internal/domain/fileUrl/entity"
	"log"
	"time"

	"github.com/go-redis/redis"
)

type fileUrlCacheImpl struct {
	cacheClient *redis.Client
}

func NewFileUrlCache(cacheClient *redis.Client) cache.FileUrlCache {
	return &fileUrlCacheImpl{
		cacheClient: cacheClient,
	}
}

func (r *fileUrlCacheImpl) PutFileUrlInfo(ctx context.Context, transactionId string, entity []entity.CreateFileUrlResultEntity) error {
	// 1. 데이터를 JSON 문자열로 변환
	// 키워드 + 트랜잭션 ID 결합
	cacheKey := consts.RedisFileUrlPrefix + transactionId

	data, err := json.Marshal(entity)
	if err != nil {
		return err
	}

	return r.cacheClient.WithContext(ctx).Set(cacheKey, data, time.Hour).Err()

}

func (r *fileUrlCacheImpl) GetFileUrlInfo(ctx context.Context, transactionId string) ([]entity.CreateFileUrlResultEntity, error) {

	cacheKey := consts.RedisFileUrlPrefix + transactionId

	val, err := r.cacheClient.WithContext(ctx).Get(cacheKey).Result()
	log.Println("## val :", val)
	if err == redis.Nil || val == "" {
		log.Println("## nil")
		return nil, consts.ErrCacheResultNotFound // 데이터가 없는 경우 (Cache Miss)
	} else if err != nil {
		log.Println("## err :", err)
		return nil, err // 기타 에러
	}

	// 2. JSON 문자열을 실제 객체 구조체로 변환 (Unmarshal)
	var entity []entity.CreateFileUrlResultEntity
	if err := json.Unmarshal([]byte(val), &entity); err != nil {
		log.Println("## err 2 :", err)
		return nil, err
	}

	log.Println("entity : ", entity)

	return entity, nil
}

func (r *fileUrlCacheImpl) PutUploadEndFileInfo(ctx context.Context, transactionId string, entity []entity.CreateFileUrlResultEntity) error {
	// 1. 데이터를 JSON 문자열로 변환
	// 키워드 + 트랜잭션 ID 결합
	successCacheKey := consts.RedisFileUploadEndPrefix + transactionId

	data, err := json.Marshal(entity)
	if err != nil {
		return err
	}

	oldCacheKey := consts.RedisFileUrlPrefix + transactionId
	// 세 번째 인자인 expiration에 0을 전달하면 만료 시간이 설정되지 않습니다.
	err = r.cacheClient.WithContext(ctx).Set(successCacheKey, data, 0).Err()
	if err != nil {
		return err
	}

	err = r.cacheClient.WithContext(ctx).Del(oldCacheKey).Err()
	if err != nil {
		return err
	}
	return nil
}
