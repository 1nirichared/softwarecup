package redis

import (
	"backend/config"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client

func InitRedis() error {
	cfg := config.GlobalConfig.Redis
	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// 测试连接
	ctx := context.Background()
	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		return err
	}

	log.Println("Redis connected successfully")
	return nil
}

func GetRedis() *redis.Client {
	return RDB
}

// 设置缓存
func SetCache(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return RDB.Set(ctx, key, value, expiration).Err()
}

// 获取缓存
func GetCache(ctx context.Context, key string) (string, error) {
	return RDB.Get(ctx, key).Result()
}

// 删除缓存
func DeleteCache(ctx context.Context, key string) error {
	return RDB.Del(ctx, key).Err()
}

// 检查键是否存在
func Exists(ctx context.Context, key string) (bool, error) {
	result, err := RDB.Exists(ctx, key).Result()
	return result > 0, err
}
