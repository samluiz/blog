package comment

import (
	"github.com/samluiz/blog/pkg/types"
)

type Service interface {
	FindCommentsByPostId(postId int) ([]*types.Comment, error)
	FindCommentById(id int) (*types.Comment, error)
	FindCommentsByUserId(userId int) ([]*types.Comment, error)
	CreateComment(input *types.CreateCommentInput) (*types.Comment, error)
	UpdateComment(id int, input *types.UpdateCommentInput) (*types.Comment, error)
	DeleteComment(id int) error
	DeleteCommentsByPostId(postId int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) FindCommentById(id int) (*types.Comment, error) {
	return s.repo.FindCommentById(id)
}

func (s *service) FindCommentsByPostId(postId int) ([]*types.Comment, error) {
	return s.repo.FindCommentsByPostId(postId)
}

func (s *service) FindCommentsByUserId(userId int) ([]*types.Comment, error) {
	return s.repo.FindCommentsByUserId(userId)
}

func (s *service) CreateComment(input *types.CreateCommentInput) (*types.Comment, error) {
	return s.repo.CreateComment(input)
}

func (s *service) UpdateComment(id int, input *types.UpdateCommentInput) (*types.Comment, error) {
	return s.repo.UpdateComment(id, input)
}

func (s *service) DeleteComment(id int) error {
	return s.repo.DeleteComment(id)
}

func (s *service) DeleteCommentsByPostId(postId int) error {
	return s.repo.DeleteCommentsByPostId(postId)
}