package timestamp

import (
	"database/sql"

	"github.com/AdmFjalar/CS301.3-Time-Tracker/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateTimestamp(timestamp types.TimeStamp) error {
	_, err := s.db.Exec("INSERT INTO timestamps (stamp_type, user_id, year, month, day, hour, minute, second) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		timestamp.StampType, timestamp.UserID, timestamp.Year, timestamp.Month, timestamp.Day, timestamp.Hour, timestamp.Minute, timestamp.Second)
	return err
}

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

func (s *Store) DeleteTimestamp(timestampID uint32) error {
	_, err := s.db.Exec("DELETE FROM timestamps WHERE timestamp_id = ?", timestampID)
	return err
}

func (s *Store) UpdateTimestamp(timestamp types.TimeStamp) error {
	_, err := s.db.Exec("UPDATE timestamps SET stamp_type = ?, user_id = ?, year = ?, month = ?, day = ?, hour = ?, minute, second = ? WHERE timestamp_id = ?",
		timestamp.StampType, timestamp.UserID, timestamp.Year, timestamp.Month, timestamp.Day, timestamp.Hour, timestamp.Minute, timestamp.Second, timestamp.TimeStampID)
	return err
}
