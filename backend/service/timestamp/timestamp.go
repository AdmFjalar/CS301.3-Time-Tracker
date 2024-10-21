package timestamp

import (
	"database/sql"

	"github.com/AdmFjalar/CS301.3-Time-Tracker/types"
)

// Store struct represents the storage for timestamps.
// It holds a reference to the database connection.
type Store struct {
	db *sql.DB
}

// NewStore creates a new Store instance with the given database connection.
func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// CreateTimestamp inserts a new timestamp into the database.
func (s *Store) CreateTimestamp(timestamp types.TimeStamp) error {
	_, err := s.db.Exec("INSERT INTO timestamps (stamp_type, user_id, year, month, day, hour, minute, second) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		timestamp.StampType, timestamp.UserID, timestamp.Year, timestamp.Month, timestamp.Day, timestamp.Hour, timestamp.Minute, timestamp.Second)
	return err
}

// GetUserTimestamps retrieves all timestamps for a given user from the database.
func (s *Store) GetUserTimestamps(userID uint32) ([]types.TimeStamp, error) {
	rows, err := s.db.Query("SELECT * FROM timestamps WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var timestamps []types.TimeStamp
	for rows.Next() {
		var timestamp types.TimeStamp
		err := rows.Scan(&timestamp.StampType, &timestamp.UserID, &timestamp.TimeStampID, &timestamp.Year, &timestamp.Month, &timestamp.Day, &timestamp.Hour, &timestamp.Minute, &timestamp.Second)
		if err != nil {
			return nil, err
		}
		timestamps = append(timestamps, timestamp)
	}

	return timestamps, nil
}

// GetTimestampByID retrieves a timestamp by its ID from the database.
func (s *Store) GetTimestampByID(timestampID uint32) (*types.TimeStamp, error) {
	row := s.db.QueryRow("SELECT * FROM timestamps WHERE timestamp_id = ?", timestampID)

	var timestamp types.TimeStamp
	err := row.Scan(&timestamp.StampType, &timestamp.UserID, &timestamp.TimeStampID, &timestamp.Year, &timestamp.Month, &timestamp.Day, &timestamp.Hour, &timestamp.Minute, &timestamp.Second)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &timestamp, nil
}

// DeleteTimestamp deletes a timestamp by its ID from the database.
func (s *Store) DeleteTimestamp(timestampID uint32) error {
	_, err := s.db.Exec("DELETE FROM timestamps WHERE timestamp_id = ?", timestampID)
	return err
}

// UpdateTimestamp updates an existing timestamp in the database.
func (s *Store) UpdateTimestamp(timestamp types.TimeStamp) error {
	_, err := s.db.Exec("UPDATE timestamps SET stamp_type = ?, user_id = ?, year = ?, month = ?, day, hour, minute, second = ? WHERE timestamp_id = ?",
		timestamp.StampType, timestamp.UserID, timestamp.Year, timestamp.Month, timestamp.Day, timestamp.Hour, timestamp.Minute, timestamp.Second, timestamp.TimeStampID)
	return err
}
