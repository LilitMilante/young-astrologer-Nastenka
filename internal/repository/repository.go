package repository

import (
	"context"
	"database/sql"
	"errors"

	"young-astrologer-Nastenka/internal/entity"
	"young-astrologer-Nastenka/internal/service"
)

type rowsScanner interface {
	Scan(...any) error
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateAPOD(ctx context.Context, apod entity.APOD) error {
	q := `INSERT INTO apods (title, explanation, date, url, image_b64) VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.ExecContext(ctx, q, apod.Title, apod.Explanation, apod.Date, apod.Url, apod.ImageB64)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) APODs(ctx context.Context) ([]entity.APOD, error) {
	q := `SELECT title, explanation, date, image_b64 FROM apods`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apods []entity.APOD

	for rows.Next() {
		apod, err := scanAPOD(rows)
		if err != nil {
			return nil, err
		}

		apods = append(apods, apod)
	}

	return apods, nil
}

func (r *Repository) APODByDate(ctx context.Context, date string) (entity.APOD, error) {
	q := `SELECT title, explanation, date, image_b64 FROM apods WHERE date = $1`
	row := r.db.QueryRowContext(ctx, q, date)

	return scanAPOD(row)
}

func scanAPOD(rs rowsScanner) (entity.APOD, error) {
	var apod entity.APOD

	err := rs.Scan(&apod.Title, &apod.Explanation, &apod.Date, &apod.ImageB64)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.APOD{}, service.ErrNotFound
		}
		return entity.APOD{}, err
	}

	return apod, nil
}
