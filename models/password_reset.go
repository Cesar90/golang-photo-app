package models

import (
	"database/sql"
	"fmt"
	"time"
)

const (
	DefaultResetDuration = 1 * time.Hour
)

type PasswordReset struct {
	ID     int
	UserID int
	// Token is only set when a PasswordReset is being created.
	Token     string
	TokenHash string
	ExpiresAt time.Time
}

type PasswordResetSevice struct {
	DB *sql.DB
	// BytesPerToken is used to determine how many bytes to use when generating
	// each password reset token. If this values is not set or is less than hte
	// MinBytesPerToken const it will be ignored and MinBytestPerToken will be
	// used
	BytesPerToken int
	// Duration is the amount of time that a PasswordReset is valid for.
	// Defaults to DefaultResetDuration
	Duration time.Duration
}

func (service *PasswordResetSevice) Create(email string) (*PasswordReset, error) {
	return nil, fmt.Errorf("TODO: Implement PasswordResetService.Create")
}

func (service *PasswordResetSevice) Consume(token string) (*User, error) {
	return nil, fmt.Errorf("TODO: Implement PasswordResetService.Consume")
}
