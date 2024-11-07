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

func (s *TimestampStore) GetUserFeed(ctx context.Context, userID int64, fq PaginatedFeedQuery) ([]Timestamp, error) {
	query := `
		SELECT 
			p.id, p.user_id, p.time, p.created_at, p.version
		FROM timestamps p
		LEFT JOIN users u ON p.user_id = u.id
		WHERE 
			p.user_id = ?
		GROUP BY p.id
		ORDER BY p.created_at ` + fq.Sort + `
		LIMIT ? OFFSET ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userID, fq.Limit, fq.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feed []Timestamp
	for rows.Next() {
		var p Timestamp
		var rawStampTime, rawCreatedAt []byte // Temporarily hold time fields as byte slices

		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&rawStampTime, // Scan into rawStampTime as []byte
			&rawCreatedAt, // Scan into rawCreatedAt as []byte
			&p.Version,
		)
		if err != nil {
			return nil, err
		}

		// Parse rawStampTime into a time.Time value
		p.StampTime, err = time.Parse("2006-01-02 15:04:05", string(rawStampTime))
		if err != nil {
			return nil, err
		}

		// Parse rawCreatedAt into a time.Time value
		p.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(rawCreatedAt))
		if err != nil {
			return nil, err
		}

		feed = append(feed, p)
	}

	return feed, nil
}

func (s *TimestampStore) Create(ctx context.Context, timestamp *Timestamp) error {
	query := `INSERT INTO timestamps (user_id, stamp_type, time) VALUES (?, ?, ?)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	result, err := s.db.ExecContext(
		ctx,
		query,
		timestamp.UserID,
		timestamp.StampType,
		timestamp.StampTime,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	timestamp.ID = id

	return nil
}

func (s *TimestampStore) GetByID(ctx context.Context, id int64) (*Timestamp, error) {
	query := `
		SELECT id, user_id, stamp_type, time, created_at, updated_at, version
		FROM timestamps
		WHERE id = ?
		`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var timestamp Timestamp
	var rawStampTime, rawCreatedAt, rawUpdatedAt []byte // Temporarily hold time fields as byte slices

	err := s.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&timestamp.ID,
		&timestamp.UserID,
		&timestamp.StampType,
		&rawStampTime, // Scan into rawStampTime as []byte
		&rawCreatedAt, // Scan into rawCreatedAt as []byte
		&rawUpdatedAt, // Scan into rawUpdatedAt as []byte
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

	// Parse rawStampTime into a time.Time value
	timestamp.StampTime, err = time.Parse("2006-01-02 15:04:05", string(rawStampTime))
	if err != nil {
		return nil, err
	}

	// Parse rawCreatedAt into a time.Time value
	timestamp.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(rawCreatedAt))
	if err != nil {
		return nil, err
	}

	// Parse rawUpdatedAt into a time.Time value
	timestamp.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", string(rawUpdatedAt))
	if err != nil {
		return nil, err
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
	// SQL query to update a timestamp based on its ID and version, and to increment the version
	query := `
		UPDATE timestamps
		SET
			user_id = ?,
			stamp_type = ?,
			time = ?,
			version = version + 1
		WHERE id = ? AND version = ?
		RETURNING version
	`

	// Set up the context with a timeout for the query execution
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	// Execute the query with the provided parameters
	err := s.db.QueryRowContext(
		ctx,
		query,
		timestamp.UserID,
		timestamp.StampType,
		timestamp.StampTime,
		timestamp.ID,
		timestamp.Version,
	).Scan(
		&timestamp.Version, // Scan the returned version into the timestamp struct
	)

	// Handle errors that occur during query execution
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			// If no rows were affected, return a "not found" error
			return ErrNotFound
		default:
			// Return any other errors encountered during the update
			return err
		}
	}

	// Return nil if the update was successful
	return nil
}
