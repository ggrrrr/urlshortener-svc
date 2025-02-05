package pg

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/ggrrrr/urlshortener-svc/be/common/application"
	"github.com/ggrrrr/urlshortener-svc/be/shorturl/internal/repo"
)

type (
	Config struct {
		Host     string `env:"DB_HOST"`
		Port     int    `env:"DB_PORT" envDefault:"5432"`
		Username string `env:"DB_USERNAME"`
		Password string `env:"DB_PASSWORD"`
		Database string `env:"DB_DATABASE"`
		SSLMode  string `env:"DB_SSL_MODE" envDefault:"disable"`
	}

	Repo struct {
		db *sql.DB
	}
)

func NewRepo(db *sql.DB) *Repo {
	return &Repo{
		db: db,
	}
}

// Create implements repo.ShortURLRepo.
func (r *Repo) Create(ctx context.Context, record repo.NewRecord) error {
	sql := fmt.Sprintf("insert into short_url (owner, key, long_url) values($1, $2, $3)")
	_, err := r.db.ExecContext(ctx, sql, record.Owner, record.Key, record.LongURL)
	if err != nil {
		return application.NewSystemError("unable to insert", err)
	}

	return nil
}

// Delete implements repo.ShortURLRepo.
func (r *Repo) Delete(ctx context.Context, key string) error {
	sql := fmt.Sprintf("delete from short_url where key = $1")
	_, err := r.db.ExecContext(ctx, sql, key)
	if err != nil {
		return application.NewSystemError("unable to insert", err)
	}
	return nil
}

// GetByKey implements repo.ShortURLRepo.
func (r *Repo) GetByKey(ctx context.Context, key string) (*repo.URLRecord, error) {
	query := fmt.Sprintf("select key, owner, long_url, created_at, updated_at from short_url where key = $1")
	row := r.db.QueryRowContext(ctx, query, key)
	if row.Err() != nil {
		return nil, row.Err()
	}
	out := repo.URLRecord{}
	err := row.Scan(&out.Key, &out.Owner, &out.LongURL, &out.CreatedAt, &out.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &out, nil
}

// ListByOwner implements repo.ShortURLRepo.
func (r *Repo) ListByOwner(ctx context.Context, owner string) ([]*repo.URLRecord, error) {
	query := fmt.Sprintf("select key, owner, long_url, created_at, updated_at from short_url where owner = $1")
	rows, err := r.db.QueryContext(ctx, query, owner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []*repo.URLRecord{}
	for rows.Next() {
		row := repo.URLRecord{}

		err := rows.Scan(&row.Key, &row.Owner, &row.LongURL, &row.CreatedAt, &row.UpdatedAt)
		if err != nil {
			return nil, err
		}

		out = append(out, &row)

	}

	return out, nil
}

// Update implements repo.ShortURLRepo.
func (r *Repo) Update(ctx context.Context, key string, longURL string) error {
	sql := fmt.Sprintf("update short_url set long_url = $2, updated_at = now() where key = $1")
	_, err := r.db.ExecContext(ctx, sql, key, longURL)
	if err != nil {
		return application.NewSystemError("unable to update", err)
	}

	return nil
}

var _ (repo.ShortURLRepo) = (*Repo)(nil)
