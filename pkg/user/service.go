package user

import "github.com/samluiz/blog/pkg/types"

type Service interface {
	FindUserById(id int) (*types.GetUserOutput, error)
	FindUserByUsername(username string) (*types.GetUserOutput, error)
	FindExternalUserByUsername(username string, provider string) (*types.GetExternalUserOutput, error)
	FindExternalUserByProviderId(id int, provider string) (*types.GetExternalUserOutput, error)
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

func (s *service) FindExternalUserByUsername(username string, provider string) (*types.GetExternalUserOutput, error) {
	return s.repo.FindExternalUserByUsername(username, provider)
}

func (s *service) FindExternalUserByProviderId(id int, provider string) (*types.GetExternalUserOutput, error) {
	return s.repo.FindExternalUserByProviderId(id, provider)
}

func (s *service) SaveUser(user *types.CreateExternalUserInput) (*types.GetExternalUserOutput, error) {
	return s.repo.SaveUser(user)
}
