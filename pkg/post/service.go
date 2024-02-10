package post

import (
	"github.com/samluiz/blog/common"
	"github.com/samluiz/blog/pkg/types"
)

type Service interface {
	FindPostById(id int) (*types.GetPostOutput, error)
	FindPostsByUserId(userId int, pagination common.Pagination) ([]*types.GetPostOutput, int, error)
	CreatePost(input *types.CreatePostInput) (*types.GetPostOutput, error)
	UpdatePost(id int, input *types.UpdatePostInput) (*types.GetPostOutput, error)
	PublishPost(id int, input *types.PublishPostInput) (*types.GetPostOutput, error)
	DeletePost(id int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) FindPostById(id int) (*types.GetPostOutput, error) {
	return s.repo.FindPostById(id)
}

func (s *service) FindPostsByUserId(userId int, pagination common.Pagination) ([]*types.GetPostOutput, int, error) {
	return s.repo.FindPostsByUserId(userId, pagination)
}

func (s *service) CreatePost(input *types.CreatePostInput) (*types.GetPostOutput, error) {
	return s.repo.CreatePost(input)
}

func (s *service) UpdatePost(id int, input *types.UpdatePostInput) (*types.GetPostOutput, error) {
	return s.repo.UpdatePost(id, input)
}

func (s *service) PublishPost(id int, input *types.PublishPostInput) (*types.GetPostOutput, error) {
	return s.repo.PublishPost(id, input)
}

func (s *service) DeletePost(id int) error {
	return s.repo.DeletePost(id)
}