package db

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// New godoc
//
//	@Summary		Creates a new database connection
//	@Description	Initializes a new database connection with the given parameters
//	@Tags			database
//	@Produce		json
//	@Param			addr			query		string	true	"Database address"
//	@Param			maxOpenConns	query		int		true	"Maximum open connections"
//	@Param			maxIdleConns	query		int		true	"Maximum idle connections"
//	@Param			maxIdleTime		query		string	true	"Maximum idle time"
//	@Success		200				{object}	sql.DB	"Database connection"
//	@Failure		500				{object}	error	"Internal server error"
//	@Router			/database/new [get]
func New(addr string, maxOpenConns, maxIdleConns int, maxIdleTime string) (*sql.DB, error) {
	db, err := sql.Open("mysql", addr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)

	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
