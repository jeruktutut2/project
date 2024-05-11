package util

import (
	"context"
	"log"
	"strconv"
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

func NewRedisConnection(host string, port int, db int) RedisUtil {
	println(time.Now().String()+" redis: connecting to", host)
	rdb := redis.NewClient(&redis.Options{
		Addr:     host + ":" + strconv.Itoa(port),
		Password: "", // no password set
		DB:       db, // use default DB
	})
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalln("redis connection error:", err)
	}
	println(time.Now().String()+" redis: connected to", host)
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
}
