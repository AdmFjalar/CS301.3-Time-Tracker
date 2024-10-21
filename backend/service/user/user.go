package user

import (
	"database/sql"
	"fmt"
	"time"

	"net/http"

	"github.com/AdmFjalar/CS301.3-Time-Tracker/service/auth"
	"github.com/AdmFjalar/CS301.3-Time-Tracker/types"
	"github.com/AdmFjalar/CS301.3-Time-Tracker/utils"
	"github.com/gorilla/mux"
)

// Handler struct represents the HTTP handler for user-related operations.
// It holds a reference to the UserStore for database operations.
type Handler struct {
	store types.UserStore
}

// NewHandler creates a new Handler instance with the given UserStore.
func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

// RegisterRoutes registers the HTTP routes for user-related operations.
func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/logout", h.handleLogout).Methods("POST")
	router.HandleFunc("/timestamps", h.handleCreateTimestamp).Methods("POST")
}

// handleLogin handles the user login request.
// It validates the user credentials and generates a JWT token if successful.
func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	user, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid email or password"))
		return
	}

	if !auth.CheckPasswordHash(payload.Password, user.Password) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid email or password"))
		return
	}

	token, err := auth.GenerateJWT(user.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

// handleLogout handles the user logout request.
// It clears the JWT token or invalidates the session.
func (h *Handler) handleLogout(w http.ResponseWriter, r *http.Request) {
	// Implement logout logic here
	w.Write([]byte("logout"))
}

// handleCreateTimestamp handles the request to create a new timestamp.
// It parses the request payload, sets the current timestamp values, and stores it in the database.
func (h *Handler) handleCreateTimestamp(w http.ResponseWriter, r *http.Request) {
	var timestamp types.TimeStamp
	if err := utils.ParseJSON(r, &timestamp); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	timestamp.Year = int16(time.Now().Year())
	timestamp.Month = uint8(time.Now().Month())
	timestamp.Day = uint8(time.Now().Day())
	timestamp.Hour = uint8(time.Now().Hour())
	timestamp.Minute = uint8(time.Now().Minute())
	timestamp.Second = uint8(time.Now().Second())

	if err := h.store.CreateTimestamp(timestamp); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, timestamp)
}

// Store struct represents the storage for user-related operations.
// It holds a reference to the database connection.
type Store struct {
	db *sql.DB
}

// NewStore creates a new Store instance with the given database connection.
func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// GetUserByEmail retrieves a user by their email from the database.
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

// CreateTimestamp inserts a new timestamp into the database.
func (s *Store) CreateTimestamp(timestamp types.TimeStamp) error {
	_, err := s.db.Exec("INSERT INTO timestamps (stamp_type, user_id, year, month, day, hour, minute, second) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		timestamp.StampType, timestamp.UserID, timestamp.Year, timestamp.Month, timestamp.Day, timestamp.Hour, timestamp.Minute, timestamp.Second)
	return err
}

// ScanRowIntoUser scans a database row into a User struct.
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
