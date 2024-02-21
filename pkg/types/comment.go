package types

import "errors"

type Comment struct {
	ID        int    `db:"id"`
	Content   string `db:"content"`
	PostID    int    `db:"post_id"`
	AuthorID  int    `db:"author_id"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

type CreateCommentInput struct {
	Content  string `db:"content"`
	PostID   int    `db:"post_id"`
	AuthorID int    `db:"author_id"`
}

type UpdateCommentInput struct {
	Content string `db:"content"`
}

var (
	ErrCommentNotFound = errors.New("comment not found")
)
