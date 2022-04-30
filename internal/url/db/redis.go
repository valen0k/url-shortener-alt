package db

import (
	"context"
	"time"
	"url-shortener-alt/internal/url"
	"url-shortener-alt/pkg/client/redis"
)

type repository struct {
	client redis.Client
}

func (r *repository) Create(ctx context.Context, url string, hash string, expiration time.Duration) (string, error) {
	return r.client.Set(ctx, hash, url, expiration).Result()
}

func (r *repository) FindByHash(ctx context.Context, hash string) (string, error) {
	return r.client.Get(ctx, hash).Result()
}

func NewRepository(client redis.Client) url.Repository {
	return &repository{client: client}
}
