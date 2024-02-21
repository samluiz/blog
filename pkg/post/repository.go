package post

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
	FindPostById(id int) (*types.GetPostOutput, error)
	FindPostsByUserId(userId int, pagination pagination.Pagination) ([]*types.GetPostOutput, int, error)
	CreatePost(input *types.CreatePostInput) (*types.GetPostOutput, error)
	UpdatePost(id int, input *types.UpdatePostInput) (*types.GetPostOutput, error)
	PublishPost(id int, input *types.PublishPostInput) (*types.GetPostOutput, error)
	DeletePost(id int) error
	PostExists(id int) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db}
}

func (r *repository) FindPostById(id int) (*types.GetPostOutput, error) {
	var post types.GetPostOutput
	err := r.db.Get(&post, "SELECT * FROM posts WHERE id = ?", id)
	if err != nil {
		return nil, types.ErrPostNotFound
	}
	return &post, nil
}

func (r *repository) FindPostsByUserId(userId int, pagination pagination.Pagination) ([]*types.GetPostOutput, int, error) {

	// Paginated requests will return the total pages to be sent as a response header in the API

	userRepo := user.NewRepository(r.db)

	if err := userRepo.UserExistsById(userId); err != nil {
		return nil, 0, err
	}

	var posts []*types.GetPostOutput

	var totalItems int

	err := r.db.Get(&totalItems, "SELECT COUNT(*) FROM posts WHERE author_id = ?", userId)

	if err != nil {
		return nil, 0, err
	}

	limit, offset, totalPages, orderBy, sortBy, err := pagination.GetValues(totalItems)

	if err != nil {
		return nil, totalPages, err
	}

	err = r.db.Select(&posts, "SELECT * FROM posts WHERE author_id = ? ORDER BY ? ? LIMIT ? OFFSET ?", userId, orderBy, sortBy, limit, offset)

	if err != nil {
		return nil, 0, err
	}
	return posts, totalPages, nil
}

func (r *repository) CreatePost(input *types.CreatePostInput) (*types.GetPostOutput, error) {

	userRepo := user.NewRepository(r.db)

	if err := userRepo.UserExistsById(input.AuthorID); err != nil {
		return nil, err
	}

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

	slug_id := slug.GenerateSlugId()
	slug := slug.GenerateSlug(input.Title, slug_id)

	res := r.db.MustExec("INSERT INTO posts (title, slug, slug_id, content, tags, author_id, visibility, is_published, published_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?) RETURNING *", input.Title, slug, slug_id, input.Content, tagsString, input.AuthorID, visibility, isPublishedAtInt, published_at)

	idCreated, err := res.LastInsertId()

	if err != nil {
		return nil, err
	}

	err = r.db.Get(&post, "SELECT * from posts WHERE id = ?", idCreated)

	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *repository) UpdatePost(id int, input *types.UpdatePostInput) (*types.GetPostOutput, error) {
	var post types.GetPostOutput

	postToBeUpdated, err := r.FindPostById(id)

	if err != nil {
		return nil, err
	}

	slug := slug.GenerateSlug(input.Title, postToBeUpdated.SlugID)

	tagsString := strings.Join(input.Tags, ",")

	_, err = r.db.Exec("UPDATE posts SET title = ?, slug = ?, content = ?, tags = ?, updated_at = ? WHERE id = ?", input.Title, slug, input.Content, tagsString, time.Now(), id)
	if err != nil {
		return nil, err
	}

	err = r.db.Get(&post, "SELECT * FROM posts WHERE id = ?", id)
	if err != nil {
		return nil, types.ErrPostNotFound
	}
	return &post, nil
}

func (r *repository) PublishPost(id int, input *types.PublishPostInput) (*types.GetPostOutput, error) {
	var post types.GetPostOutput

	if err := r.PostExists(id); err != nil {
		return nil, err
	}

	visibility := types.PUBLIC

	now := time.Now()

	_, err := r.db.Exec("UPDATE posts SET is_published = ?, published_at = ?, visibility = ?, updated_at = ? WHERE id = ?", input.IsPublished, now, visibility, now, id)
	if err != nil {
		return nil, err
	}
	err = r.db.Get(&post, "SELECT * FROM posts WHERE id = ?", id)
	if err != nil {
		return nil, types.ErrPostNotFound
	}
	return &post, nil
}

func (r *repository) DeletePost(id int) error {

	if err := r.PostExists(id); err != nil {
		return err
	}

	_, err := r.db.Exec("DELETE FROM posts WHERE id = ? CASCADE", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) PostExists(id int) error {
	var count int
	err := r.db.Get(&count, "SELECT COUNT(*) FROM posts WHERE id = ?", id)

	if err != nil {
		if err == sql.ErrNoRows {
			return types.ErrPostNotFound
		}
		return err
	}

	return nil
}