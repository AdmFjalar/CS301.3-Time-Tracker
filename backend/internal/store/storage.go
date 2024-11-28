package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("resource not found")
	ErrConflict          = errors.New("resource already exists")
	QueryTimeoutDuration = 5 * time.Second
)

// Storage struct holds the interfaces for interacting with different data stores.
type Storage struct {
	// Timestamps interface provides methods for managing timestamps in the database.
	Timestamps interface {
		GetByID(context.Context, int64) (*Timestamp, error)
		Create(context.Context, *Timestamp) error
		Delete(context.Context, int64) error
		Update(context.Context, *Timestamp) error
		GetUserFeed(context.Context, int64, Query) ([]Timestamp, error)
		GetLatestTimestamp(context.Context, int64) (*Timestamp, error)
		GetFinishedShifts(context.Context, int64) ([]Shift, error)
	}

	// Users interface provides methods for managing users in the database.
	Users interface {
		GetAll(context.Context) ([]*User, error)
		GetByID(context.Context, int64) (*User, error)
		GetByEmail(context.Context, string) (*User, error)
		Create(context.Context, *sql.Tx, *User) error
		CreateAndInvite(ctx context.Context, user *User, token string, exp time.Duration) error
		Activate(context.Context, string) error
		Update(context.Context, *User) error
		ChangePassword(context.Context, *User) error
		ResetPassword(context.Context, string, *User) error
		RequestPasswordAndEmailReset(context.Context, *User, string, time.Duration) error
		Delete(context.Context, int64) error
	}

	// Roles interface provides methods for managing roles in the database.
	Roles interface {
		GetByName(context.Context, string) (*Role, error)
	}
}

// NewStorage creates a new Storage instance with the provided database connection.
func NewStorage(db *sql.DB) Storage {
	return Storage{
		Timestamps: &TimestampStore{db},
		Users:      &UserStore{db},
		Roles:      &RoleStore{db},
	}
}

// withTx executes a function within a database transaction.
func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
