package post

import (
	"math"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/samluiz/blog/common"
	"github.com/samluiz/blog/pkg/types"
)

type Repository interface {
	FindPostById(id int) (*types.GetPostOutput, error)
	FindPostsByUserId(userId int, pagination common.Pagination) ([]*types.GetPostOutput, int, error)
	CreatePost(input *types.CreatePostInput) (*types.GetPostOutput, error)
	UpdatePost(id int, input *types.UpdatePostInput) (*types.GetPostOutput, error)
	PublishPost(id int, input *types.PublishPostInput) (*types.GetPostOutput, error)
	DeletePost(id int) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db}
}

func (r *repository) FindPostById(id int) (*types.GetPostOutput, error) {
	var post types.GetPostOutput
	err := r.db.Get(&post, "SELECT * FROM posts WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *repository) FindPostsByUserId(userId int, pagination common.Pagination) ([]*types.GetPostOutput, int, error) {

	// Paginated requests will return the total pages to be sent as a response header in the API

	var posts []*types.GetPostOutput

	var offset int
	var limit int

	if pagination.Page > 0 {
		offset = (pagination.Page - 1) * pagination.Size
	} else {
		offset = 0
	}

	if pagination.Size > 0 {
		limit = pagination.Size
	} else {
		limit = 10
	}

	if pagination.OrderBy == "" {
		pagination.OrderBy = "created_at"
	}
	if pagination.SortBy == "" {
		pagination.SortBy = "DESC"
	}

	var totalItems int

	err := r.db.Get(&totalItems, "SELECT COUNT(*) FROM posts WHERE author_id = $1", userId)

	if err != nil {
		return nil, 0, err
	}

	totalPages := int(math.Round(float64(totalItems) / float64(pagination.Size)))

	err = r.db.Select(&posts, "SELECT * FROM posts WHERE author_id = $1 ORDER BY $2 $3 LIMIT $4 OFFSET $5", userId, pagination.OrderBy, pagination.SortBy, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return posts, totalPages, nil
}

func (r *repository) CreatePost(input *types.CreatePostInput) (*types.GetPostOutput, error) {
	var post types.GetPostOutput
	tagsString := strings.Join(input.Tags, ",")

	var published_at interface{} = nil
	var isPublishedAtInt int
	visibility := types.PRIVATE

	if input.IsPublished {
		isPublishedAtInt = 1
		published_at = time.Now()
		visibility = types.PUBLIC
	}

	slug := strings.ReplaceAll(input.Title, " ", "-")

	res := r.db.MustExec("INSERT INTO posts (title, slug, content, tags, author_id, visibility, is_published, published_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *", input.Title, slug, input.Content, tagsString, input.AuthorID, visibility, isPublishedAtInt, published_at)

	idCreated, err := res.LastInsertId()

	if err != nil {
		return nil, err
	}

	err = r.db.Get(&post, "SELECT * from posts WHERE id = $1", idCreated)

	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *repository) UpdatePost(id int, input *types.UpdatePostInput) (*types.GetPostOutput, error) {
	var post types.GetPostOutput
	tagsString := strings.Join(input.Tags, ",")

	_, err := r.db.Exec("UPDATE posts SET title = $1, content = $2, tags = $3, updated_at = $4 WHERE id = $5", input.Title, input.Content, tagsString, time.Now(), id)
	if err != nil {
		return nil, err
	}
	err = r.db.Get(&post, "SELECT * FROM posts WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *repository) PublishPost(id int, input *types.PublishPostInput) (*types.GetPostOutput, error) {
	var post types.GetPostOutput

	visibility := types.PUBLIC

	now := time.Now()

	_, err := r.db.Exec("UPDATE posts SET is_published = $1, published_at = $2, visibility = $3, updated_at = $4 WHERE id = $5", input.IsPublished, now, visibility, now, id)
	if err != nil {
		return nil, err
	}
	err = r.db.Get(&post, "SELECT * FROM posts WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *repository) DeletePost(id int) error {
	_, err := r.db.Exec("DELETE FROM posts WHERE id = $1 CASCADE", id)
	if err != nil {
		return err
	}
	return nil
}