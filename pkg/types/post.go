package types

import "time"

const (
	PRIVATE = "PRIVATE"
	PUBLIC  = "PUBLIC"
)

type Post struct {
	ID          int 			`db:"id"`
	Title       string 		`db:"title"`
	Slug        string 		`db:"slug"`
	Content     string 		`db:"content"`
	Tags        []string	`db:"tags"`
	AuthorID    int 			`db:"author_id"`
	Visibility  string 		`db:"visibility"`
	IsPublished bool 			`db:"is_published"`
	PublishedAt time.Time `db:"published_at"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type GetPostOutput struct {
	ID          int 			`db:"id"`
	Title       string 		`db:"title"`
	Slug        string 		`db:"slug"`
	Content     string 		`db:"content"`
	Tags        string		`db:"tags"`
	AuthorID    int 			`db:"author_id"`
	Visibility  string 		`db:"visibility"`
	IsPublished bool 			`db:"is_published"`
	PublishedAt time.Time `db:"published_at"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}


type CreatePostInput struct {
	Title       string 		`db:"title" validate:"required"`
	Content     string 		`db:"content" validate:"required"`
	IsPublished bool 			`db:"is_published" validate:"required"`
	AuthorID    int 			`db:"author_id" validate:"required"`
	Tags        []string 	`db:"tags"`
}

type UpdatePostInput struct {
	Title   string 				`db:"title"`
	Content string 				`db:"content"`
	Tags    []string 			`db:"tags"`
}

type PublishPostInput struct {
	IsPublished bool `db:"is_published"`
}