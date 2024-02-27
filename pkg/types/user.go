package types

import (
	"errors"
	"time"
)

type User struct {
	ID        int       `db:"id"`
	Username  string    `db:"username"`
	Password  string    `db:"password"`
	IsAdmin   bool      `db:"is_admin"`
	Avatar    string    `db:"avatar"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type GetUserOutput struct {
	ID        int       `db:"id"`
	Username  string    `db:"username"`
	Password  string    `db:"password"`
	IsAdmin   bool      `db:"is_admin"`
	Avatar    string    `db:"avatar"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type CreateExternalUserInput struct {
	ProviderId int    `json:"provider_id"`
	Name       string `json:"name"`
	Username   string `json:"username"`
	Provider   string `json:"provider"`
	Avatar     string `json:"avatar"`
}

type GetExternalUserOutput struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	Username  string    `db:"username"`
	Provider  string    `db:"provider"`
	Avatar    string    `db:"avatar"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrUserUnauthorized = errors.New("user is not authorized")
)
