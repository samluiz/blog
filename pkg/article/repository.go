package article

import (
	"database/sql"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/samluiz/blog/common/pagination"
	"github.com/samluiz/blog/common/slug"
	"github.com/samluiz/blog/pkg/types"
	"github.com/samluiz/blog/pkg/user"
)

type Repository interface {
	FindArticleById(id int) (*types.GetArticleOutput, error)
	FindArticlesByUserId(userId int, pagination pagination.Pagination) ([]*types.GetArticleOutput, int, error)
	CreateArticle(input *types.CreateArticleInput) (*types.GetArticleOutput, error)
	UpdateArticle(id int, input *types.UpdateArticleInput) (*types.GetArticleOutput, error)
	PublishArticle(id int, input *types.PublishArticleInput) (*types.GetArticleOutput, error)
	DeleteArticle(id int) error
	ArticleExists(id int) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db}
}

func (r *repository) FindArticleById(id int) (*types.GetArticleOutput, error) {
	var article types.GetArticleOutput
	err := r.db.Get(&article, "SELECT * FROM articles WHERE id = ?", id)
	if err != nil {
		return nil, types.ErrArticleNotFound
	}
	return &article, nil
}

func (r *repository) FindArticlesByUserId(userId int, pagination pagination.Pagination) ([]*types.GetArticleOutput, int, error) {

	// Paginated requests will return the total pages to be sent as a response header in the API

	userRepo := user.NewRepository(r.db)

	if err := userRepo.UserExistsById(userId); err != nil {
		return nil, 0, err
	}

	var articles []*types.GetArticleOutput

	var totalItems int

	err := r.db.Get(&totalItems, "SELECT COUNT(*) FROM articles WHERE author_id = ?", userId)

	if err != nil {
		return nil, 0, err
	}

	limit, offset, totalPages, orderBy, sortBy, err := pagination.GetValues(totalItems)

	if err != nil {
		return nil, totalPages, err
	}

	err = r.db.Select(&articles, "SELECT * FROM articles WHERE author_id = ? ORDER BY ? ? LIMIT ? OFFSET ?", userId, orderBy, sortBy, limit, offset)

	if err != nil {
		return nil, 0, err
	}
	return articles, totalPages, nil
}

func (r *repository) CreateArticle(input *types.CreateArticleInput) (*types.GetArticleOutput, error) {

	userRepo := user.NewRepository(r.db)

	if err := userRepo.UserExistsById(input.AuthorID); err != nil {
		return nil, err
	}

	var article types.GetArticleOutput
	tagsString := strings.Join(input.Tags, ",")

	var published_at interface{} = nil
	var isPublishedAtInt int
	visibility := types.PRIVATE

	if input.IsPublished {
		isPublishedAtInt = 1
		published_at = time.Now()
		visibility = types.PUBLIC
	}

	slug_id := slug.GenerateSlugId()
	slug := slug.GenerateSlug(input.Title, slug_id)

	res := r.db.MustExec("INSERT INTO articles (title, slug, slug_id, content, tags, author_id, visibility, is_published, published_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?) RETURNING *", input.Title, slug, slug_id, input.Content, tagsString, input.AuthorID, visibility, isPublishedAtInt, published_at)

	idCreated, err := res.LastInsertId()

	if err != nil {
		return nil, err
	}

	err = r.db.Get(&article, "SELECT * from articles WHERE id = ?", idCreated)

	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *repository) UpdateArticle(id int, input *types.UpdateArticleInput) (*types.GetArticleOutput, error) {
	var article types.GetArticleOutput

	articleToBeUpdated, err := r.FindArticleById(id)

	if err != nil {
		return nil, err
	}

	slug := slug.GenerateSlug(input.Title, articleToBeUpdated.SlugID)

	tagsString := strings.Join(input.Tags, ",")

	_, err = r.db.Exec("UPDATE articles SET title = ?, slug = ?, content = ?, tags = ?, updated_at = ? WHERE id = ?", input.Title, slug, input.Content, tagsString, time.Now(), id)
	if err != nil {
		return nil, err
	}

	err = r.db.Get(&article, "SELECT * FROM articles WHERE id = ?", id)
	if err != nil {
		return nil, types.ErrArticleNotFound
	}
	return &article, nil
}

func (r *repository) PublishArticle(id int, input *types.PublishArticleInput) (*types.GetArticleOutput, error) {
	var article types.GetArticleOutput

	if err := r.ArticleExists(id); err != nil {
		return nil, err
	}

	visibility := types.PUBLIC

	now := time.Now()

	_, err := r.db.Exec("UPDATE articles SET is_published = ?, published_at = ?, visibility = ?, updated_at = ? WHERE id = ?", input.IsPublished, now, visibility, now, id)
	if err != nil {
		return nil, err
	}
	err = r.db.Get(&article, "SELECT * FROM articles WHERE id = ?", id)
	if err != nil {
		return nil, types.ErrArticleNotFound
	}
	return &article, nil
}

func (r *repository) DeleteArticle(id int) error {

	if err := r.ArticleExists(id); err != nil {
		return err
	}

	_, err := r.db.Exec("DELETE FROM articles WHERE id = ? CASCADE", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) ArticleExists(id int) error {
	var count int
	err := r.db.Get(&count, "SELECT COUNT(*) FROM articles WHERE id = ?", id)

	if err != nil {
		if err == sql.ErrNoRows {
			return types.ErrArticleNotFound
		}
		return err
	}

	return nil
}
