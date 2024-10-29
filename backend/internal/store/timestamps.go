package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Timestamp struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	StampType string    `json:"stamp_type"`
	Year      int       `json:"year"`
	Month     int       `json:"month"`
	Day       int       `json:"day"`
	Hour      int       `json:"hour"`
	Minute    int       `json:"minute"`
	Second    int       `json:"second"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int       `json:"version"`
}

type TimestampStore struct {
	db *sql.DB
}

func (s *TimestampStore) Create(ctx context.Context, timestamp *Timestamp) error {
	query := `INSERT INTO timestamps (user_id, stamp_type, year, month, day, hour, minute, second) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		timestamp.UserID,
		timestamp.StampType,
		timestamp.Year,
		timestamp.Month,
		timestamp.Day,
		timestamp.Hour,
		timestamp.Minute,
		timestamp.Second,
	).Scan(
		&timestamp.ID,
		&timestamp.CreatedAt,
		&timestamp.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *TimestampStore) GetByID(ctx context.Context, id int64) (*Timestamp, error) {
	query := `
		SELECT id, user_id, stamp_type, year, month, day, hour, minute, second, created_at, updated_at, version
		FROM timestamps
		WHERE id = $1
		`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var timestamp Timestamp
	err := s.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&timestamp.ID,
		&timestamp.UserID,
		&timestamp.StampType,
		&timestamp.Year,
		&timestamp.Month,
		&timestamp.Day,
		&timestamp.Hour,
		&timestamp.Minute,
		&timestamp.Second,
		&timestamp.CreatedAt,
		&timestamp.UpdatedAt,
		&timestamp.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &timestamp, nil
}

func (s *TimestampStore) Delete(ctx context.Context, timestampID int64) error {
	query := `DELETE FROM timestamps WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, timestampID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *TimestampStore) Update(ctx context.Context, timestamp *Timestamp) error {
	query := `
		UPDATE timestamps
		SET
			user_id = $1,
			stamp_type = $2,
			year = $3,
			month = $4,
			day = $5,
			hour = $6,
			minute = $7,
			second = $8,
			version + 1
		WHERE id = $9 AND version = $10
		RETURNING version
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		timestamp.UserID,
		timestamp.StampType,
		timestamp.Year,
		timestamp.Month,
		timestamp.Day,
		timestamp.Hour,
		timestamp.Minute,
		timestamp.Second,
		timestamp.ID,
		timestamp.Version,
	).Scan(
		&timestamp.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}
	return nil
}
