package url

import "context"

type Repository interface {
	Create(ctx context.Context, url Url) (string, error)
	FindByHash(ctx context.Context, hashUrl string) (Url, error)
	FindByOriginal(ctx context.Context, originalUrl string) (Url, error)
	CheckHash(ctx context.Context, hashUrl string) bool
}
