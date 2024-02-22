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

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrUserUnauthorized = errors.New("user is not authorized")
)
