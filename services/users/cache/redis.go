package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

var Ctx = context.Background()
var Rdb *redis.Client

func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "redis:6379", // nombre del contenedor Docker
		Password: "",           // sin password por defecto
		DB:       0,            // base de datos por defecto
	})

	_, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}
	log.Println("✅ Connected to Redis")
}

func Set(key string, value string, ttl time.Duration) {
	err := Rdb.Set(Ctx, key, value, ttl).Err()
	if err != nil {
		log.Printf("Redis SET error: %v", err)
	}
}

func Get(key string) (string, error) {
	val, err := Rdb.Get(Ctx, key).Result()
	if err == redis.Nil {
		return "", nil // no existe
	} else if err != nil {
		return "", err
	}
	return val, nil
}

func Delete(key string) {
	err := Rdb.Del(Ctx, key).Err()
	if err != nil {
		log.Printf("Redis DEL error: %v", err)
	}
}
