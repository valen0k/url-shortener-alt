package url

import (
	"context"
	"time"
)

type Repository interface {
	Create(ctx context.Context, url string, hash string, expiration time.Duration) (string, error)
	FindByHash(ctx context.Context, hash string) (string, error)
}
