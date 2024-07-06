package initialize

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func SetDataRedis(client *redis.Client, ctx context.Context, key string, value interface{}, expiration time.Duration) {
	_, err := client.Set(ctx, key, value, expiration).Result()
	if err != nil {
		log.Fatalln("error when setting data redis:", err.Error())
	}
	log.Println("set data redis succedded")
}

func DelDataRedis(client *redis.Client, ctx context.Context, key string) {
	_, err := client.Del(ctx, key).Result()
	if err != nil {
		log.Fatalln("error when deleting data redis:", err.Error())
	}
	log.Println("delete data redis succeded")
}
