package routes

import (
	"errors"
	"html/template"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/samluiz/blog/api/integrations"
	"github.com/samluiz/blog/api/parsers"
	"github.com/samluiz/blog/pkg/user"
	"golang.org/x/crypto/bcrypt"
)

const IS_LOGGED = "is_logged"
const DASHBOARD_URL = "/blog/admin/dashboard"

type Router interface {
	HomePage(c *fiber.Ctx) error
	ArticlePage(c *fiber.Ctx) error
	LoginPage(c *fiber.Ctx) error
	AdminDashboardPage(c *fiber.Ctx) error
	AdminArticlesPartial(c *fiber.Ctx) error
	Authenticate(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
}

type router struct {
	app         *fiber.App
	store       *session.Store
	userService user.Service
}

func NewRouter(app *fiber.App, store *session.Store, userService user.Service) Router {
	return &router{app, store, userService}
}

func (r *router) HomePage(c *fiber.Ctx) error {
	articles, err := integrations.GetArticlesFromDevTo()

	if err != nil {
		log.Default().Println(err.Error())
	}

	return c.Render("pages/home", fiber.Map{
		"Articles":    articles,
		"PageTitle":   "home",
		"Description": "My personal portfolio, but also a blog about software development, programming, and technology. Articles about web development, backend, frontend, and whatever i wanna share.",
		"Error":       err,
	})
}

func (r *router) ArticlePage(c *fiber.Ctx) error {
	slug := c.Params("slug")

	article, err := integrations.GetArticleBySlugDevTo(slug)

	if err != nil {
		log.Default().Println(err.Error())
	}

	htmlContent := template.HTML(article.BodyHTML)
	markdownContent := template.HTML(parsers.MarkdownToHTML([]byte(article.BodyMarkdown)))

	return c.Render("pages/article", fiber.Map{
		"Article":     article,
		"HTML":        htmlContent,
		"Markdown":    markdownContent,
		"PageTitle":   article.Slug,
		"Description": article.Description,
		"Error":       err,
	})
}

func (r *router) LoginPage(c *fiber.Ctx) error {

	session, err := r.store.Get(c)

	if err != nil {
		log.Default().Printf("Error getting session: %v", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	isLogged := session.Get(IS_LOGGED)

	if isLogged != nil && isLogged == true {
		log.Default().Println("User is already logged in")
		return c.Redirect(DASHBOARD_URL)
	}

	return c.Render("pages/login", fiber.Map{
		"PageTitle": "login",
	})
}

func (r *router) AdminDashboardPage(c *fiber.Ctx) error {

	return c.Render("pages/dashboard", nil)
}

func (r *router) AdminArticlesPartial(c *fiber.Ctx) error {

	return c.SendFile("views/partials/articles.html")
}

func (r *router) Authenticate(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	user, err := r.userService.FindUserByUsername(username)

	// TODO: pass this errors to the htmx view
	if err != nil {
		log.Default().Printf("Error finding user: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"error":   err.Error(),
		})
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"error":   errors.New("wrong password").Error(),
		})
	}

	session, err := r.store.Get(c)

	if err != nil {
		log.Default().Printf("Error getting session: %v", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	session.Set("username", username)
	session.Set("is_admin", user.IsAdmin)
	session.Set(IS_LOGGED, true)

	err = session.Save()

	if err != nil {
		log.Default().Printf("Error saving session: %v", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Redirect(DASHBOARD_URL)
}

func (r *router) Logout(c *fiber.Ctx) error {
	session, err := r.store.Get(c)

	if err != nil {
		log.Default().Printf("Error getting session: %v", err)
		return c.Redirect("/")
	}

	session.Destroy()

	return c.Redirect("/")
}
