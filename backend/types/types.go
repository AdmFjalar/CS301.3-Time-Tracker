package types

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

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

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
}

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	CreateTimestamp(timestamp TimeStamp) error
}
