package types

// LoginUserPayload represents the payload for user login requests.
// It includes the user's email and password.
type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// TimeStamp represents a timestamp entry.
// It includes the type of timestamp, user ID, timestamp ID, and the date and time components.
type TimeStamp struct {
	StampType   string `json:"stamp_type"`
	UserID      uint32 `json:"user_id"`
	TimeStampID uint32 `json:"timestamp_id"`
	Year        int16  `json:"year"`
	Month       uint8  `json:"month"`
	Day         uint8  `json:"day"`
	Hour        uint8  `json:"hour"`
	Minute      uint8  `json:"minute"`
	Second      uint8  `json:"second"`
}

// User represents a user in the system.
// It includes the user's ID, first name, last name, email, password, and the creation timestamp.
type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
}

// UserStore defines the interface for user-related database operations.
// It includes methods for retrieving a user by email and creating a timestamp.
type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	CreateTimestamp(timestamp TimeStamp) error
}
