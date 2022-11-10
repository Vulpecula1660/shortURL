package redis

import (
	"context"

	db "shortURL/db/sqlc"
)

type RedisQuerier interface {
	SetBloom(ctx context.Context, shortUrl string) (bool, error)
	ExistBloom(ctx context.Context, shortUrl string) (bool, error)
	SetData(ctx context.Context, shortUrl string, url db.Url) error
	GetData(ctx context.Context, shortUrl string) (db.Url, bool, error)
}

var _ RedisQuerier = (*RedisQueries)(nil)
