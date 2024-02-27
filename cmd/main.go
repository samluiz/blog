package main

import (
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"
	"github.com/samluiz/blog/api/middlewares/isinternal"
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

	// Fiber config
	config := fiber.Config{
		Views:       engine,
		ViewsLayout: "layout",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			if code == fiber.StatusNotFound {
				return c.Redirect("/error/404")
			}

			return c.Redirect("/error?status=" + strconv.Itoa(code))
		},
	}

	// App
	app := fiber.New(config)

	app.Static("/static", "static", fiber.Static{
		CacheDuration: 0,
		MaxAge:        0,
	})

	// Recover middleware
	app.Use(recover.New())

	// Logger middleware
	app.Use(logger.New(logger.Config{
		Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}\n",
	}))

	// Middleware that checks if user is logged in by session
	islogged := islogged.New(islogged.Config{
		Session: store,
	})

	isinternal := isinternal.New()

	// Internal routes
	internal := app.Group("/internal")
	internal.Use(isinternal)

	// Protected routes
	protected := app.Group("/dashboard")
	protected.Use(islogged)

	// Error routes
	errors := app.Group("/error")

	// Router
	router := routes.NewRouter(app, store, userService)

	// App root routes
	app.Get("/", router.HomePage)

	// Blog routes
	app.Get("/auth/login", router.LoginPage)
	app.Get("/articles/:slug", router.ArticlePage)
	app.Get("/articles", router.ArticlesPage)

	// Error routes
	errors.Get("/", router.ErrorPage)
	errors.Get("/404", router.NotFoundPage)

	// Protected routes
	protected.Get("/", router.AdminDashboardPage)

	// Auth routes
	internal.Post("/auth/login", router.Authenticate)

	// Server
	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))
}
