package user

import (
	"database/sql"
	"fmt"

	"github.com/AdmFjalar/CS301.3-Time-Tracker/backend/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	rows, err := s.db.Query(("SELECT * FROM users WHERE email= ?"), email)
	if err != nil {
		return nil, err
	}

	u := new(types.User)
	{
		for rows.Next() {
			u, err = ScanRowIntoUser(rows)
			if err != nil {
				return nil, err
			}
		}

		if u.ID == 0 {
			return nil, fmt.Errorf("user not found")
		}
	}

	return u, nil
}

func ScanRowIntoUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)

	err := rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}
