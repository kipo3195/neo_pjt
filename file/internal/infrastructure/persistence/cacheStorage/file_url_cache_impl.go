package cacheStorage

import (
	"context"
	"encoding/json"
	"file/internal/consts"
	"file/internal/domain/fileUrl/cache"
	"file/internal/domain/fileUrl/entity"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type fileUrlCacheImpl struct {
	cacheClient *redis.ClusterClient
}

func NewFileUrlCache(cacheClient *redis.ClusterClient) cache.FileUrlCache {
	return &fileUrlCacheImpl{
		cacheClient: cacheClient,
	}
}

func (r *fileUrlCacheImpl) PutFileUrlInfo(ctx context.Context, transactionId string, entity []entity.CreateFileUrlResultEntity) error {
	// 1. 데이터를 JSON 문자열로 변환
	// 키워드 + 트랜잭션 ID 결합
	key := consts.RedisFileUrlPrefix + transactionId

	data, err := json.Marshal(entity)
	if err != nil {
		return err
	}

	log.Println("[PutFileUrlInfo] redis save cacheKey :", key)
	log.Println("[PutFileUrlInfo] redis save tid :", transactionId)
	log.Println("[PutFileUrlInfo] redis save entity :", data)

	err = r.cacheClient.Set(ctx, key, data, time.Hour).Err()

	if err != nil {
		log.Println("err : ", err)
	}

	return err

}

func (r *fileUrlCacheImpl) GetFileUrlInfo(ctx context.Context, transactionId string) ([]entity.CreateFileUrlResultEntity, error) {

	key := consts.RedisFileUrlPrefix + transactionId

	log.Println("[GetFileUrlInfo] redis tid :", transactionId)
	log.Println("[GetFileUrlInfo] redis cacheKey :", key)
	val, err := r.cacheClient.Get(ctx, key).Result()
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
	err = r.cacheClient.Set(ctx, successCacheKey, data, 0).Err()
	if err != nil {
		return err
	}

	err = r.cacheClient.Del(ctx, oldCacheKey).Err()
	if err != nil {
		return err
	}
	return nil
}
