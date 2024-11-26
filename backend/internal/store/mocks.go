package store

import (
	"context"
	"database/sql"
	"time"
)

func NewMockStore() Storage {
	return Storage{
		Users: &MockUserStore{},
	}
}

type MockUserStore struct{}

func (m *MockUserStore) Create(ctx context.Context, tx *sql.Tx, u *User) error {
	return nil
}

func (m *MockUserStore) GetAll(ctx context.Context) ([]*User, error) {
	return []*User{}, nil
}

func (m *MockUserStore) GetByID(ctx context.Context, userID int64) (*User, error) {
	return &User{ID: userID}, nil
}

func (m *MockUserStore) GetByEmail(context.Context, string) (*User, error) {
	return &User{}, nil
}

func (m *MockUserStore) CreateAndInvite(ctx context.Context, user *User, token string, exp time.Duration) error {
	return nil
}

func (m *MockUserStore) Activate(ctx context.Context, t string) error {
	return nil
}

func (m *MockUserStore) Update(ctx context.Context, u *User) error {
	return nil
}

func (m *MockUserStore) ChangePassword(ctx context.Context, u *User) error {
	return nil
}

func (m *MockUserStore) ResetPassword(ctx context.Context, t string) error {
	return nil
}

func (m *MockUserStore) RequestPasswordAndEmailReset(ctx context.Context, user *User, token string, exp time.Duration) error {
	return nil
}

func (m *MockUserStore) Delete(ctx context.Context, id int64) error {
	return nil
}
