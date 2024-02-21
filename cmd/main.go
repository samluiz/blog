package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"
	"github.com/samluiz/blog/api/routes"
	"github.com/samluiz/blog/pkg/config"
	"github.com/samluiz/blog/pkg/user"
)

func main() {

	db, err := config.NewConnection()

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userService := user.NewService(user.NewRepository(db))

	store := session.New()

	engine := html.New("views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
		ViewsLayout: "layout",
	})

	app.Static("/static", "static")

	app.Use(logger.New(logger.Config{
    Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}\n",
	}))

	router := routes.NewRouter(app, store, userService)

	app.Get("/", router.Home)
	app.Post("/auth/login", router.Authenticate)
	app.Get("/admin", router.LoginPage)
	app.Get("/blog/posts/:slug", router.Post)

	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))
}