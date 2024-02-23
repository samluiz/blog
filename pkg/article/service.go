package article

import (
	"github.com/samluiz/blog/common/pagination"
	"github.com/samluiz/blog/pkg/types"
)

type Service interface {
	FindArticleById(id int) (*types.GetArticleOutput, error)
	FindArticlesByUserId(userId int, pagination pagination.Pagination) ([]*types.GetArticleOutput, int, error)
	CreateArticle(input *types.CreateArticleInput) (*types.GetArticleOutput, error)
	UpdateArticle(id int, input *types.UpdateArticleInput) (*types.GetArticleOutput, error)
	PublishArticle(id int, input *types.PublishArticleInput) (*types.GetArticleOutput, error)
	DeleteArticle(id int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) FindArticleById(id int) (*types.GetArticleOutput, error) {
	return s.repo.FindArticleById(id)
}

func (s *service) FindArticlesByUserId(userId int, pagination pagination.Pagination) ([]*types.GetArticleOutput, int, error) {
	return s.repo.FindArticlesByUserId(userId, pagination)
}

func (s *service) CreateArticle(input *types.CreateArticleInput) (*types.GetArticleOutput, error) {
	return s.repo.CreateArticle(input)
}

func (s *service) UpdateArticle(id int, input *types.UpdateArticleInput) (*types.GetArticleOutput, error) {
	return s.repo.UpdateArticle(id, input)
}

func (s *service) PublishArticle(id int, input *types.PublishArticleInput) (*types.GetArticleOutput, error) {
	return s.repo.PublishArticle(id, input)
}

func (s *service) DeleteArticle(id int) error {
	return s.repo.DeleteArticle(id)
}
