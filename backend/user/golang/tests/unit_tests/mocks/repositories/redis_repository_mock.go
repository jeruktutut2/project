package mockrepositories

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
)

type RedisRepositoryMock struct {
	Mock mock.Mock
}

func (repository *RedisRepositoryMock) Set(client *redis.Client, ctx context.Context, key string, value interface{}, expiration time.Duration) (result string, err error) {
	arguments := repository.Mock.Called(client, ctx, key, value, expiration)
	return arguments.Get(0).(string), arguments.Error(1)
}

func (repository *RedisRepositoryMock) Del(client *redis.Client, ctx context.Context, key string) (result int64, err error) {
	arguments := repository.Mock.Called(client, ctx, key)
	return arguments.Get(0).(int64), arguments.Error(1)
}
