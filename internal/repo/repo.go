package repo

import (
	"context"
	"time"
)

type (
	NewRecord struct {
		Owner   string
		Key     string
		LongURL string
	}

	URLRecord struct {
		Key       string
		Owner     string
		LongURL   string
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	ShortURLRepo interface {
		Create(ctx context.Context, record NewRecord) error
		GetByKey(ctx context.Context, key string) (string, error)
		GetLongURL(ctx context.Context, longURL string) (string, error)
		Delete(ctx context.Context, key string) error
		Update(ctx context.Context, key string, longURL string) error
		List(ctx context.Context, owner string) ([]*URLRecord, error)
	}
)
