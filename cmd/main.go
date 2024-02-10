package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

type PostOutput struct {
	ID          int
	Title       string
	PublishedAt string
}

func main() {

	Posts := map[string][]PostOutput {
		"Posts": {
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
	}}

	engine := html.New("views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/static", "static")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", Posts)
	})

	port := os.Getenv("PORT")

	log.Fatal(app.Listen(":" + port))
}