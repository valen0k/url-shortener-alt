package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"log"
	"strings"
	"url-shortener-alt/internal/url"
	"url-shortener-alt/pkg/client/postgresql"
)

type repository struct {
	client postgresql.Client
}

func (r repository) CheckHash(ctx context.Context, hashUrl string) bool {
	query := `
			SELECT EXISTS(
			    SELECT * 
			    FROM public.shortener 
			    WHERE hash_url = $1)
			`
	log.Println(formatQuery(query))
	var res bool
	if err := r.client.QueryRow(ctx, query, hashUrl).Scan(&res); err != nil {
		return false
	}
	return res
}

func (r repository) Create(ctx context.Context, url url.Url) (string, error) {
	query := `
			INSERT INTO public.shortener (hash_url, original_url)
			VALUES ($1, $2)
			RETURNING id
			`
	log.Println(formatQuery(query))
	if err := r.client.QueryRow(ctx, query, url.HashUrl, url.OriginalUrl).Scan(&url.ID); err != nil {
		return "", err
	}
	return url.ID, nil
}

func (r repository) FindByHash(ctx context.Context, hashUrl string) (url.Url, error) {
	query := `
			SELECT id, hash_url, original_url, created_at 
			FROM public.shortener 
			WHERE hash_url = $1
			`
	log.Println(formatQuery(query))
	var u url.Url
	if err := r.client.QueryRow(ctx, query, hashUrl).Scan(&u.ID, &u.HashUrl, &u.OriginalUrl, &u.CreatedAt); err != nil {
		return url.Url{}, err
	}
	return u, nil
}

func (r repository) FindByOriginal(ctx context.Context, originalUrl string) (url.Url, error) {
	query := `
			SELECT id, hash_url, original_url, created_at 
			FROM public.shortener 
			WHERE original_url = $1
			`
	log.Println(formatQuery(query))
	var u url.Url
	if err := r.client.QueryRow(ctx, query, originalUrl).Scan(&u.ID, &u.HashUrl, &u.OriginalUrl, &u.CreatedAt); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			log.Println(newErr)
			return url.Url{ID: "nil"}, newErr
		}
		return url.Url{}, err
	}
	return u, nil
}

func formatQuery(query string) string {
	return strings.ReplaceAll(strings.ReplaceAll(query, "\t", ""), "\n", "")
}

func NewRepository(client postgresql.Client) url.Repository {
	return &repository{client: client}
}
