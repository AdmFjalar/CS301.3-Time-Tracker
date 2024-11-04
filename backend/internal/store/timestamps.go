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
	StampTime time.Time `json:"stamp_time"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int       `json:"version"`
}

type TimestampStore struct {
	db *sql.DB
}

func (s *TimestampStore) Create(ctx context.Context, timestamp *Timestamp) error {
	query := `INSERT INTO timestamps (user_id, stamp_type, stamp_time) VALUES (?, ?, ?)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		timestamp.UserID,
		timestamp.StampType,
	).Scan(
		&timestamp.ID,
		&timestamp.StampTime,
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
		SELECT id, user_id, stamp_type, stamp_time, created_at, updated_at, version
		FROM timestamps
		WHERE id = ?
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
		&timestamp.StampTime,
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
	query := `DELETE FROM timestamps WHERE id = ?`

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
			user_id = ?,
			stamp_type = ?,
			stamp_time = ?,
			second = ?,
			version + 1
		WHERE id = ? AND version = ?
		RETURNING version
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		timestamp.UserID,
		timestamp.StampType,
		timestamp.StampTime,
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
