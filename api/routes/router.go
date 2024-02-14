package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type Router interface {
	Index(c *fiber.Ctx) error
	Post(c *fiber.Ctx) error
}

type router struct {
	app *fiber.App
}

func NewRouter(app *fiber.App) Router {
	return &router{app}
}

type PostPreview struct {
	ID          int
	Title       string
	PublishedAt string
}

type PostOutput struct {
	ID          int
	Title       string
	PublishedAt string
	Content 	 string
	Slug 		 string
}

func (r *router) Index(c *fiber.Ctx) error {

	Posts := []PostPreview {
		{
			ID:        1,
			Title:     "How I've built my blog using Go + Astro + Htmx",
			PublishedAt: time.Now().Format("2006.01.02"),
		},
		{
			ID:        2,
			Title:     "How i became the first millionare in Piauí",
			PublishedAt: time.Now().Format("2006.01.02"),
		},
		{
			ID:        3,
			Title:     "How to build a bazingly fast nutrition API using Java",
			PublishedAt: time.Now().Format("2006.01.02"),
		},
}

	return c.Render("pages/home", fiber.Map{
		"Posts": Posts,
	})
}

func (r *router) Post(c *fiber.Ctx) error {

	var Post = PostOutput{
		ID:          1,
		Title:       "How I've built my blog using Go + Astro + Htmx",
		PublishedAt: time.Now().Format("2006.01.02"),
		Content:     "This is the content of the post",
		Slug:        "how-ive-built-my-blog-using-go-astro-htmx",
	}

	return c.Render("pages/post", Post)
}