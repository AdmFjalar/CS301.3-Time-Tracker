package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Timestamp represents a timestamp entry in the system.
type Timestamp struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	StampType string    `json:"stamp_type"`
	StampTime time.Time `json:"stamp_time"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int       `json:"version"`
}

// Shift represents a work shift with sign-in, sign-out, and break times.
type Shift struct {
	SignIn         time.Time     `json:"SignIn"`
	SignOut        time.Time     `json:"SignOut"`
	Breaks         [][]time.Time `json:"Breaks"`
	TotalBreakTime float64       `json:"TotalBreakTime"` // TotalBreakTime in seconds (float64)
	TotalShiftTime float64       `json:"TotalShiftTime"` // TotalShiftTime in seconds (float64)
	NetWorkTime    float64       `json:"NetWorkTime"`    // NetWorkTime in seconds (float64)
}

// TimestampStore provides methods for managing timestamps in the database.
type TimestampStore struct {
	db *sql.DB
}

// GetUserFeed godoc
//
//	@Summary		Retrieves the user feed
//	@Description	Retrieves the user feed based on the provided query parameters
//	@Tags			timestamps
//	@Produce		json
//	@Param			userID	query		int		true	"User ID"
//	@Param			limit	query		int		false	"Limit"
//	@Param			offset	query		int		false	"Offset"
//	@Param			sort	query		string	false	"Sort"
//	@Success		200		{object}	[]Timestamp
//	@Failure		500		{object}	error
//	@Router			/timestamps/feed [get]
func (s *TimestampStore) GetUserFeed(ctx context.Context, userID int64, fq Query) ([]Timestamp, error) {
	query := `
		SELECT 
			p.id, p.user_id, p.stamp_type, p.time, p.created_at, p.version
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
			&p.StampType,  // Scan the stamp_type directly into the Timestamp struct
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

// GetLatestTimestamp godoc
//
//	@Summary		Retrieves the latest timestamp
//	@Description	Retrieves the most recent timestamp for a specific user
//	@Tags			timestamps
//	@Produce		json
//	@Param			userID	query		int		true	"User ID"
//	@Success		200		{object}	Timestamp
//	@Failure		500		{object}	error
//	@Router			/timestamps/latest [get]
func (s *TimestampStore) GetLatestTimestamp(ctx context.Context, userID int64) (*Timestamp, error) {
	// Define a query to fetch only the latest timestamp for the user
	fq := Query{
		Limit:  1,
		Offset: 0,
		Sort:   "desc", // Sort by created_at in descending order to get the latest timestamp
	}

	// Call GetUserFeed with the modified PaginatedFeedQuery
	timestamps, err := s.GetUserFeed(ctx, userID, fq)
	if err != nil {
		return nil, err
	}

	// Check if any timestamps were retrieved
	if len(timestamps) == 0 {
		return nil, nil
	}

	// Return the first (and only) timestamp in the list
	return &timestamps[0], nil
}

// GetFinishedShifts godoc
//
//	@Summary		Retrieves finished shifts
//	@Description	Retrieves finished shifts for a specific user
//	@Tags			shifts
//	@Produce		json
//	@Param			userID	query		int		true	"User ID"
//	@Success		200		{object}	[]Shift
//	@Failure		500		{object}	error
//	@Router			/shifts [get]
func (s *TimestampStore) GetFinishedShifts(ctx context.Context, userID int64) ([]Shift, error) {
	query := `
		SELECT stamp_type, time
		FROM timestamps
		WHERE user_id = ?
		ORDER BY time ASC
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shifts []Shift
	var currentShift *Shift
	var currentBreakStart time.Time

	for rows.Next() {
		var stampType string
		var stampTimeRaw string
		if err := rows.Scan(&stampType, &stampTimeRaw); err != nil {
			return nil, err
		}

		// Parse the raw timestamp string with the correct format
		stampTime, err := time.Parse("2006-01-02 15:04:05", stampTimeRaw)
		if err != nil {
			return nil, fmt.Errorf("failed to parse timestamp: %v", err)
		}

		switch stampType {
		case "sign-in":
			if currentShift == nil {
				currentShift = &Shift{SignIn: stampTime}
			}

		case "start-break":
			if currentShift != nil && currentShift.SignOut.IsZero() {
				currentBreakStart = stampTime
			}

		case "end-break":
			if currentShift != nil && !currentBreakStart.IsZero() {
				breakEnd := stampTime
				currentShift.Breaks = append(currentShift.Breaks, []time.Time{currentBreakStart, breakEnd})

				// Calculate break duration in seconds (float64)
				breakDuration := breakEnd.Sub(currentBreakStart).Seconds()
				currentShift.TotalBreakTime += breakDuration // Keep TotalBreakTime as seconds (float64)

				currentBreakStart = time.Time{}
			}

		case "sign-out":
			if currentShift != nil && currentShift.SignOut.IsZero() {
				currentShift.SignOut = stampTime
				// Calculate total shift time in seconds (float64)
				currentShift.TotalShiftTime = currentShift.SignOut.Sub(currentShift.SignIn).Seconds() // TotalShiftTime in seconds
				// NetWorkTime is shift time minus break time
				currentShift.NetWorkTime = currentShift.TotalShiftTime - currentShift.TotalBreakTime

				shifts = append(shifts, *currentShift)
				currentShift = nil
			}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return shifts, nil
}

// Create godoc
//
//	@Summary		Creates a timestamp
//	@Description	Inserts a new timestamp into the database
//	@Tags			timestamps
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		Timestamp	true	"Timestamp information"
//	@Success		201		{object}	Timestamp	"Timestamp created"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/timestamps [post]
func (s *TimestampStore) Create(ctx context.Context, timestamp *Timestamp) error {
	// Define allowed previous states for each stamp type
	validTransitions := map[string]string{
		"sign-in":     "sign-out",          // Only allowed if the last stamp is "sign-out"
		"sign-out":    "sign-in,end-break", // Only allowed if the last stamp is "sign-in" or "end-break"
		"start-break": "sign-in,end-break", // Only allowed if the last stamp is "sign-in" or "end-break"
		"end-break":   "start-break",       // Only allowed if the last stamp is "start-break"
	}

	latestTimestamp, err := s.GetLatestTimestamp(ctx, timestamp.UserID)

	// Handle case where no previous timestamps exist (first action should be "sign-in")
	if latestTimestamp == nil && timestamp.StampType != "sign-in" {
		return errors.New("first action must be sign-in")
	}

	// Validate the transition if there is a previous timestamp
	if latestTimestamp != nil {
		previousType := latestTimestamp.StampType

		// Ensure no duplicate consecutive timestamps
		if previousType == timestamp.StampType {
			return errors.New("duplicate timestamp")
		}

		// Get valid previous types for the current stamp type
		validPrevTypes, exists := validTransitions[timestamp.StampType]
		if !exists {
			return errors.New("invalid stamp type")
		}

		// Check if the last stamp type is within the allowed types
		allowedTypes := strings.Split(validPrevTypes, ",")
		if !contains(allowedTypes, previousType) {
			return fmt.Errorf("invalid transition from %s to %s", previousType, timestamp.StampType)
		}
	}

	timestamp.StampTime = time.Now()

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

// GetByID godoc
//
//	@Summary		Retrieves a timestamp by ID
//	@Description	Retrieves a timestamp by its ID
//	@Tags			timestamps
//	@Produce		json
//	@Param			id	path		int	true	"Timestamp ID"
//	@Success		200	{object}	Timestamp
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Router			/timestamps/{id} [get]
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

// Delete godoc
//
//	@Summary		Deletes a timestamp
//	@Description	Removes a timestamp from the database by its ID
//	@Tags			timestamps
//	@Produce		json
//	@Param			id	path		int	true	"Timestamp ID"
//	@Success		204	{string}	string	"Timestamp deleted"
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Router			/timestamps/{id} [delete]
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

// Update godoc
//
//	@Summary		Updates a timestamp
//	@Description	Modifies an existing timestamp in the database
//	@Tags			timestamps
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int			true	"Timestamp ID"
//	@Param			payload	body		Timestamp	true	"Updated timestamp information"
//	@Success		200		{object}	Timestamp
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Router			/timestamps/{id} [patch]
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

// Helper function to check if a value exists in a slice
func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
