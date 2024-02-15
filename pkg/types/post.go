package types

import (
	"errors"
	"time"
)

const (
	PRIVATE = "PRIVATE"
	PUBLIC  = "PUBLIC"
)

type Post struct {
	ID          int 			`db:"id"`
	Title       string 		`db:"title"`
	Slug        string 		`db:"slug"`
	SlugID 			string 		`db:"slug_id"`
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
	SlugID 			string 		`db:"slug_id"`
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
	Title       string 		`db:"title"`
	Content     string 		`db:"content"`
	IsPublished bool 			`db:"is_published"`
	AuthorID    int 			`db:"author_id"`
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

var (
	ErrPostNotFound = errors.New("Post not found.")
)