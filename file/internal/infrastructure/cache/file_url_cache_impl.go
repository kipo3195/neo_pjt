package cache

import (
	"context"
	"encoding/json"
	"file/internal/consts"
	"file/internal/domain/fileUrl/cache"
	"file/internal/domain/fileUrl/entity"
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
