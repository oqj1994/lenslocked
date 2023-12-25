package M

import (
	"database/sql"
	"time"
)

const (
	DefaultDuration = 1 * time.Hour
)

type PasswordReset struct {
	ID        int
	UserID    int
	Token     string
	TokenHash string
	expiredAt time.Time
}

type PasswordResetService struct {
	BytesPerToken int
	DB            *sql.DB
	Duration      time.Duration
}

func (prs PasswordResetService) Create(email string) (*PasswordReset, error) {
	return nil, nil
}
func (prs PasswordResetService) Consume(token string) (*User, error) {
	return nil, nil
}
