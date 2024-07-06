package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepository interface {
	Set(client *redis.Client, ctx context.Context, key string, value interface{}, expiration time.Duration) (result string, err error)
	Del(client *redis.Client, ctx context.Context, key string) (result int64, err error)
}

type RedisRepositoryImplementation struct {
}

func NewRedisRepository() RedisRepository {
	return &RedisRepositoryImplementation{}
}

func (repository *RedisRepositoryImplementation) Set(client *redis.Client, ctx context.Context, key string, value interface{}, expiration time.Duration) (result string, err error) {
	return client.Set(ctx, key, value, expiration).Result()
}

func (repository *RedisRepositoryImplementation) Del(client *redis.Client, ctx context.Context, key string) (result int64, err error) {
	return client.Del(ctx, key).Result()
}
