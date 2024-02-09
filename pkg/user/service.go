package user

import "github.com/samluiz/blog/pkg/types"

type Service interface {
	FindUserById(id int) (*types.GetUserOutput, error)
	FindUserByUsername(username string) (*types.GetUserOutput, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) FindUserById(id int) (*types.GetUserOutput, error) {
	return s.repo.FindUserById(id)
}

func (s *service) FindUserByUsername(username string) (*types.GetUserOutput, error) {
	return s.repo.FindUserByUsername(username)
}