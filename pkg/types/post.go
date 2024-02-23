package types

import (
	"errors"
	"time"
)

const (
	PRIVATE = "PRIVATE"
	PUBLIC  = "PUBLIC"
)

type Article struct {
	ID          int       `db:"id"`
	Title       string    `db:"title"`
	Slug        string    `db:"slug"`
	SlugID      string    `db:"slug_id"`
	Content     string    `db:"content"`
	Tags        []string  `db:"tags"`
	AuthorID    int       `db:"author_id"`
	Visibility  string    `db:"visibility"`
	IsPublished bool      `db:"is_published"`
	PublishedAt time.Time `db:"published_at"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type GetArticleOutput struct {
	ID          int       `db:"id"`
	Title       string    `db:"title"`
	Slug        string    `db:"slug"`
	SlugID      string    `db:"slug_id"`
	Content     string    `db:"content"`
	Tags        string    `db:"tags"`
	AuthorID    int       `db:"author_id"`
	Visibility  string    `db:"visibility"`
	IsPublished bool      `db:"is_published"`
	PublishedAt time.Time `db:"published_at"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type CreateArticleInput struct {
	Title       string   `db:"title"`
	Content     string   `db:"content"`
	IsPublished bool     `db:"is_published"`
	AuthorID    int      `db:"author_id"`
	Tags        []string `db:"tags"`
}

type UpdateArticleInput struct {
	Title   string   `db:"title"`
	Content string   `db:"content"`
	Tags    []string `db:"tags"`
}

type PublishArticleInput struct {
	IsPublished bool `db:"is_published"`
}

var (
	ErrArticleNotFound = errors.New("article not found")
)
