package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/samluiz/blog/api/routes"
)

func main() {

	engine := html.New("views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
		ViewsLayout: "layout",
	})

	app.Static("/static", "static")

	router := routes.NewRouter(app)

	app.Get("/", router.Home)
	app.Get("/posts/:slug", router.Post)

	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))
}