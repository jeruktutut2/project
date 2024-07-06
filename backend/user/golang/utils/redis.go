package utils

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisUtil interface {
	GetClient() *redis.Client
	Close()
}

type RedisUtilImplementation struct {
	Client *redis.Client
}

func NewRedisConnection(host string, port string, db int) RedisUtil {
	println(time.Now().String()+" redis: connecting to", host+":"+port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: "", // no password set
		DB:       db, // use default DB
	})
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalln("redis connection error:", err)
	}
	println(time.Now().String()+" redis: connected to", host+":"+port)
	return &RedisUtilImplementation{
		Client: rdb,
	}
}

func (util *RedisUtilImplementation) GetClient() *redis.Client {
	return util.Client
}

func (util *RedisUtilImplementation) Close() {
	err := util.Client.Close()
	if err != nil {
		panic("redis close connection error: " + err.Error())
	}
	println(time.Now().String(), "redis: closed properly")
}
