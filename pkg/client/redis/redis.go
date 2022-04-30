package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
	"url-shortener-alt/internal/config"
)

type Client interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Exists(ctx context.Context, keys ...string) *redis.IntCmd
}

func NewClient(ctx context.Context, config config.StorageConfig) (client *redis.Client, err error) {
	db, err := strconv.Atoi(config.Database)
	if err != nil {
		return
	}
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Host, config.Port),
		Password: config.Password,
		DB:       db,
	})
	return
}
