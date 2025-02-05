package pg

import (
	"database/sql"
	"fmt"
	"log/slog"
)

const createTable string = `
CREATE TABLE IF NOT EXISTS short_url (
	key TEXT  not null,
	owner TEXT  not null,
	long_url TEXT  not null,
	created_at TIMESTAMP DEFAULT NOW(),
	updated_at TIMESTAMP DEFAULT NOW(),
	UNIQUE(key)
);

CREATE INDEX  IF NOT EXISTS short_url_owner_idx ON short_url (owner) ;
`
const dropTable string = `
DROP TABLE IF EXISTS short_url;

DROP TABLE IF EXISTS short_url_owner_idx ;
`

func Up(db *sql.DB) error {

	_, err := db.Exec(createTable)
	return err
}

func Down(db *sql.DB) error {

	_, err := db.Exec(dropTable)
	return err
}

func Connect(cfg Config) (*sql.DB, error) {
	if cfg.SSLMode == "" {
		cfg.SSLMode = "disable"
	}
	psqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database, cfg.SSLMode)

	db, err := sql.Open("postgres", psqlConn)
	if err != nil {
		slog.Error("Connect.Open", "error", err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		slog.Error("Connect.Ping", "error", err)
		return nil, err
	}
	slog.Info(
		"Connect.Open",
		slog.String("host", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)),
		slog.String("database", cfg.Database),
	)

	return db, nil
}
