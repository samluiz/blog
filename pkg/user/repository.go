package user

import (
	"github.com/jmoiron/sqlx"
	"github.com/samluiz/blog/pkg/types"
)

type Repository interface {
	FindUserById(id int) (*types.GetUserOutput, error)
	FindUserByUsername(username string) (*types.GetUserOutput, error)
	FindUsers() ([]*types.GetUserOutput, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db}
}

func (r *repository) FindUserById(id int) (*types.GetUserOutput, error) {
	var user types.GetUserOutput
	err := r.db.Get(&user, "SELECT id, username, is_admin, avatar, created_at, updated_at FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) FindUserByUsername(username string) (*types.GetUserOutput, error) {
	var user types.GetUserOutput
	err := r.db.Get(&user, "SELECT id, username, is_admin, avatar, created_at, updated_at FROM users WHERE username = $1", username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) FindUsers() ([]*types.GetUserOutput, error) {
	var users []*types.GetUserOutput
	err := r.db.Select(&users, "SELECT id, username, is_admin, avatar, created_at, updated_at FROM users")
	if err != nil {
		return nil, err
	}
	return users, nil
}