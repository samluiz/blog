package comment

import (
	"github.com/jmoiron/sqlx"
	"github.com/samluiz/blog/pkg/types"
)

type Repository interface {
	FindCommentsByPostId(postId int) ([]*types.Comment, error)
	FindCommentById(id int) (*types.Comment, error)
	FindCommentsByUserId(userId int) ([]*types.Comment, error)
	CreateComment(input *types.CreateCommentInput) (*types.Comment, error)
	UpdateComment(id int, input *types.UpdateCommentInput) (*types.Comment, error)
	DeleteComment(id int) error
	DeleteCommentsByPostId(postId int) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db}
}

func (r *repository) FindCommentsByPostId(postId int) ([]*types.Comment, error) {
	var comments []*types.Comment
	err := r.db.Select(&comments, "SELECT * FROM comments WHERE post_id = $1", postId)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *repository) FindCommentById(id int) (*types.Comment, error) {
	var comment types.Comment
	err := r.db.Get(&comment, "SELECT * FROM comments WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *repository) FindCommentsByUserId(userId int) ([]*types.Comment, error) {
	var comments []*types.Comment
	err := r.db.Select(&comments, "SELECT * FROM comments WHERE author_id = $1", userId)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *repository) CreateComment(input *types.CreateCommentInput) (*types.Comment, error) {
	var comment types.Comment
	err := r.db.Get(&comment, "INSERT INTO comments (author_id, post_id, content) VALUES ($1, $2, $3) RETURNING *", input.AuthorID, input.PostID, input.Content)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *repository) UpdateComment(id int, input *types.UpdateCommentInput) (*types.Comment, error) {
	var comment types.Comment
	err := r.db.Get(&comment, "UPDATE comments SET content = $1 WHERE id = $2 RETURNING *", input.Content, id)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *repository) DeleteComment(id int) error {
	_, err := r.db.Exec("DELETE FROM comments WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) DeleteCommentsByPostId(postId int) error {
	_, err := r.db.Exec("DELETE FROM comments WHERE post_id = $1", postId)
	if err != nil {
		return err
	}
	return nil
}