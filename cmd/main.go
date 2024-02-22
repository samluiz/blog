package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"
	"github.com/samluiz/blog/api/middlewares/islogged"
	"github.com/samluiz/blog/api/routes"
	"github.com/samluiz/blog/pkg/config"
	"github.com/samluiz/blog/pkg/user"
)

func main() {

	// Database
	db, err := config.NewConnection()

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Services
	userService := user.NewService(user.NewRepository(db))

	// Session
	store := session.New()

	// Html template
	engine := html.New("views", ".html")

	// App
	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "layout",
	})

	// Static files
	app.Static("/static", "static")

	// Middlewares

	// Logger middleware
	app.Use(logger.New(logger.Config{
		Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}\n",
	}))

	// Route groups
	blog := app.Group("/blog")
	admin := blog.Group("/admin")
	auth := blog.Group("/auth")

	// Middleware that checks if user is logged in by session
	admin.Use(islogged.New(islogged.Config{
		Session: store,
	}))

	// Router
	router := routes.NewRouter(app, store, userService)

	// App root routes
	app.Get("/", router.HomePage)

	// admin routes
	admin.Get("/dashboard", router.AdminDashboardPage)

	// Auth routes
	auth.Post("/login", router.Authenticate)

	// Blog routes
	blog.Get("/login", router.LoginPage)
	blog.Get("/posts/:slug", router.PostPage)

	// Server
	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))
}
