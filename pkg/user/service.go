package user

import "github.com/samluiz/blog/pkg/types"

type Service interface {
	FindUserById(id int) (*types.GetUserOutput, error)
	FindUserByUsername(username string) (*types.GetUserOutput, error)
	FindExternalUserByUsername(username string) (*types.GetExternalUserOutput, error)
	FindExternalUserById(id int) (*types.GetExternalUserOutput, error)
	SaveUser(user *types.CreateExternalUserInput) (*types.GetExternalUserOutput, error)
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

func (s *service) FindExternalUserByUsername(username string) (*types.GetExternalUserOutput, error) {
	return s.repo.FindExternalUserByUsername(username)
}

func (s *service) FindExternalUserById(id int) (*types.GetExternalUserOutput, error) {
	return s.repo.FindExternalUserById(id)
}

func (s *service) SaveUser(user *types.CreateExternalUserInput) (*types.GetExternalUserOutput, error) {
	return s.repo.SaveUser(user)
}
