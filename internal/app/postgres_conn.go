package app

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"young-astrologer-Nastenka/migrations"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

const (
	pingInterval = 500 * time.Millisecond
	pingCount    = 6
)

func ConnectToPostgres(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("connect to postgres: %w", err)
	}

	// Ping and reconnect if needed
	for i := 0; i < pingCount; i++ {
		err = db.Ping()
		if err != nil {
			time.Sleep(pingInterval)
			continue
		}

		break
	}

	if err != nil {
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	err = upMigrations(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func upMigrations(db *sql.DB) error {
	fs := migrations.FS
	goose.SetBaseFS(fs)
	goose.SetLogger(goose.NopLogger())

	err := goose.SetDialect("postgres")
	if err != nil {
		return err
	}

	err = goose.Up(db, ".")
	if err != nil && !errors.Is(err, goose.ErrNoNextVersion) {
		return err
	}

	return nil
}
