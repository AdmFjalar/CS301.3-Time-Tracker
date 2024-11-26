package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail = errors.New("a user with that email already exists")
)

type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Password  password  `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	IsActive  int       `json:"is_active"`
	RoleID    int64     `json:"role_id"`
	Role      Role      `json:"role"`
	ManagerID int64     `json:"manager_id"`
}

type password struct {
	text *string
	hash []byte
}

func (p *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.text = &text
	p.hash = hash

	return nil
}

func (p *password) Compare(text string) error {
	return bcrypt.CompareHashAndPassword(p.hash, []byte(text))
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
		INSERT INTO users (passhash, email, role_id) 
		VALUES (?, ?, (SELECT id FROM roles WHERE name = ? LIMIT 1))
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	// Default to 'user' role if no role is provided
	role := user.Role.Name
	if role == "" {
		role = "user"
	}

	// Execute the query and get the result
	result, err := tx.ExecContext(
		ctx,
		query,
		user.Password.hash,
		user.Email,
		role,
	)
	if err != nil {
		// Handle the duplicate email error by checking the error message
		if err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"` {
			return ErrDuplicateEmail
		}
		return err
	}

	// Get the last inserted ID
	userID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// Set the user ID
	user.ID = userID

	// No need to set `CreatedAt` and `UpdatedAt` as they are automatically handled by the database

	return nil
}

func (s *UserStore) GetAll(ctx context.Context) ([]*User, error) {
	query := `
		SELECT users.id, email, first_name, last_name, created_at, roles.*, manager_id
		FROM users
		JOIN roles ON (users.role_id = roles.id)
		WHERE is_active = 1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*User, 0)
	for rows.Next() {
		user := &User{}
		var rawFirstName, rawLastName sql.NullString
		var rawCreatedAt []byte // For scanning the DATETIME field
		var rawManagerID sql.NullInt64

		err := rows.Scan(
			&user.ID,
			&user.Email,
			&rawFirstName, // Use sql.NullString for nullable first_name
			&rawLastName,  // Use sql.NullString for nullable last_name
			&rawCreatedAt, // Scan created_at as raw bytes
			&user.Role.ID,
			&user.Role.Name,
			&user.Role.Level,
			&user.Role.Description,
			&rawManagerID,
		)
		if err != nil {
			return nil, err
		}

		// Assign first_name and last_name only if they are not NULL
		if rawFirstName.Valid {
			user.FirstName = rawFirstName.String
		} else {
			user.FirstName = "" // Set a default or handle as needed
		}

		if rawLastName.Valid {
			user.LastName = rawLastName.String
		} else {
			user.LastName = "" // Set a default or handle as needed
		}

		if rawManagerID.Valid {
			user.ManagerID = rawManagerID.Int64
		} else {
			user.ManagerID = 0 // Set a default or handle as needed
		}

		// Parse rawCreatedAt into a time.Time value
		user.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(rawCreatedAt))
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserStore) GetByID(ctx context.Context, userID int64) (*User, error) {
	query := `
		SELECT users.id, email, first_name, last_name, passhash, created_at, roles.*, manager_id
		FROM users
		JOIN roles ON (users.role_id = roles.id)
		WHERE users.id = ? AND is_active = 1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &User{}
	var rawFirstName, rawLastName sql.NullString
	var rawCreatedAt []byte // For scanning the DATETIME field
	var rawManagerID sql.NullInt64

	err := s.db.QueryRowContext(
		ctx,
		query,
		userID,
	).Scan(
		&user.ID,
		&user.Email,
		&rawFirstName, // Use sql.NullString for nullable first_name
		&rawLastName,  // Use sql.NullString for nullable last_name
		&user.Password.hash,
		&rawCreatedAt, // Scan created_at as raw bytes
		&user.Role.ID,
		&user.Role.Name,
		&user.Role.Level,
		&user.Role.Description,
		&rawManagerID,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	// Assign first_name and last_name only if they are not NULL
	if rawFirstName.Valid {
		user.FirstName = rawFirstName.String
	} else {
		user.FirstName = "" // Set a default or handle as needed
	}

	if rawLastName.Valid {
		user.LastName = rawLastName.String
	} else {
		user.LastName = "" // Set a default or handle as needed
	}

	if rawManagerID.Valid {
		user.ManagerID = rawManagerID.Int64
	} else {
		user.ManagerID = 0 // Set a default or handle as needed
	}

	// Parse rawCreatedAt into a time.Time value
	user.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(rawCreatedAt))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
	SELECT users.id, email, first_name, last_name, passhash, created_at, roles.*, manager_id
	FROM users
	JOIN roles ON (users.role_id = roles.id)
	WHERE users.email = ? AND is_active = 1
`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &User{}
	var rawFirstName, rawLastName sql.NullString
	var rawCreatedAt []byte // For scanning the DATETIME field
	var rawManagerID sql.NullInt64

	err := s.db.QueryRowContext(
		ctx,
		query,
		email,
	).Scan(
		&user.ID,
		&user.Email,
		&rawFirstName, // Use sql.NullString for nullable first_name
		&rawLastName,  // Use sql.NullString for nullable last_name
		&user.Password.hash,
		&rawCreatedAt, // Scan created_at as raw bytes
		&user.Role.ID,
		&user.Role.Name,
		&user.Role.Level,
		&user.Role.Description,
		&rawManagerID,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	// Assign first_name and last_name only if they are not NULL
	if rawFirstName.Valid {
		user.FirstName = rawFirstName.String
	} else {
		user.FirstName = "" // Set a default or handle as needed
	}

	if rawLastName.Valid {
		user.LastName = rawLastName.String
	} else {
		user.LastName = "" // Set a default or handle as needed
	}

	if rawManagerID.Valid {
		user.ManagerID = rawManagerID.Int64
	} else {
		user.ManagerID = 0 // Set a default or handle as needed
	}

	// Parse rawCreatedAt into a time.Time value
	user.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(rawCreatedAt))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserStore) CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.Create(ctx, tx, user); err != nil {
			return err
		}

		if err := s.createUserInvitation(ctx, tx, token, invitationExp, user.ID); err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStore) Activate(ctx context.Context, token string) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		// 1. find the user that this token belongs to
		user, err := s.getUserFromInvitation(ctx, tx, token)
		if err != nil {
			return err
		}

		// 2. update the user
		user.IsActive = 1
		if err := s.update(ctx, tx, user); err != nil {
			return err
		}

		// 3. clean the invitations
		if err := s.deleteUserInvitations(ctx, tx, user.ID); err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStore) RequestPasswordAndEmailReset(ctx context.Context, user *User, token string, exp time.Duration) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.createUserPasswordReset(ctx, tx, token, exp, user); err != nil {
			return err
		}
		return nil
	})
}

func (s *UserStore) Update(ctx context.Context, user *User) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.update(ctx, tx, user); err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStore) ChangePassword(ctx context.Context, user *User) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.updatePassword(ctx, tx, user); err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStore) ResetPassword(ctx context.Context, token string, user *User) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		// 1. update the user
		if err := s.updatePassword(ctx, tx, user); err != nil {
			return err
		}

		// 2. clean the password resets
		if err := s.deletePasswordResets(ctx, tx, user.ID); err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStore) getUserFromPasswordReset(ctx context.Context, tx *sql.Tx, token string) (*User, error) {
	query := `
		SELECT u.id, u.email, u.created_at, u.is_active
		FROM users u
		JOIN password_resets pr ON u.id = pr.user_id
		WHERE pr.token = ? AND pr.expiry > ?
	`

	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &User{}
	var rawCreatedAt []byte // For scanning the DATETIME field

	err := tx.QueryRowContext(ctx, query, hashToken, time.Now()).Scan(
		&user.ID,
		&user.Email,
		&rawCreatedAt, // Scan into rawCreatedAt as []byte
		&user.IsActive,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	// Parse rawCreatedAt into a time.Time value
	user.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(rawCreatedAt))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserStore) createUserPasswordReset(ctx context.Context, tx *sql.Tx, token string, exp time.Duration, user *User) error {
	query := `INSERT INTO password_resets (user_id, token, expiry) VALUES (?, ?, ?)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, user.ID, token, time.Now().Add(exp))
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) getUserFromInvitation(ctx context.Context, tx *sql.Tx, token string) (*User, error) {
	query := `
		SELECT u.id, u.email, u.created_at, u.is_active
		FROM users u
		JOIN user_invitations ui ON u.id = ui.user_id
		WHERE ui.token = ? AND ui.expiry > ?
	`

	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &User{}
	var rawCreatedAt []byte // For scanning the DATETIME field

	err := tx.QueryRowContext(ctx, query, hashToken, time.Now()).Scan(
		&user.ID,
		&user.Email,
		&rawCreatedAt, // Scan into rawCreatedAt as []byte
		&user.IsActive,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	// Parse rawCreatedAt into a time.Time value
	user.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(rawCreatedAt))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserStore) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, exp time.Duration, userID int64) error {
	query := `INSERT INTO user_invitations (token, user_id, expiry) VALUES (?, ?, ?)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, token, userID, time.Now().Add(exp))
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) update(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `UPDATE users SET email = ?, is_active = ?, first_name = ?, last_name = ?, manager_id = ? WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, user.Email, user.IsActive, user.FirstName, user.LastName, user.ManagerID, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) updatePassword(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `UPDATE users SET passhash = ? WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, user.Password.hash, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) deleteUserInvitations(ctx context.Context, tx *sql.Tx, userID int64) error {
	query := `DELETE FROM user_invitations WHERE user_id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) deletePasswordResets(ctx context.Context, tx *sql.Tx, userID int64) error {
	query := `DELETE FROM password_resets WHERE user_id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) Delete(ctx context.Context, userID int64) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.delete(ctx, tx, userID); err != nil {
			return err
		}

		if err := s.deleteUserInvitations(ctx, tx, userID); err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStore) delete(ctx context.Context, tx *sql.Tx, id int64) error {
	query := `DELETE FROM users WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
