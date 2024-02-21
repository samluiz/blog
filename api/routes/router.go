package routes

import (
	"errors"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/samluiz/blog/common/slug"
	"github.com/samluiz/blog/pkg/user"
	"golang.org/x/crypto/bcrypt"
)

type Router interface {
	Home(c *fiber.Ctx) error
	Post(c *fiber.Ctx) error
	LoginPage(c *fiber.Ctx) error
	Authenticate(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
}

type router struct {
	app *fiber.App
	store *session.Store
	userService user.Service
}

func NewRouter(app *fiber.App, store *session.Store, userService user.Service) Router {
	return &router{app, store, userService}
}

type PostOutput struct {
	ID          int
	Title       string
	PublishedAt string
	Content 	 string
	Slug 		 string
}

var Posts = []PostOutput{
	{
		ID:          1,
		Title:       "How I've built my blog using Go + Htmx + TailwindCSS",
		PublishedAt: time.Now().Format("2006.01.02"),
		Content:     "This is the content of the post",
		Slug:        slug.GenerateSlug("How I've built my blog using Go + Htmx + TailwindCSS", slug.GenerateSlugId()),
	},

	{
		ID:          2,
		Title:       "How to build an API using Go and Fiber",
		PublishedAt: time.Now().Format("2006.01.02"),
		Content:     "This is the content of the post",
		Slug:        slug.GenerateSlug("How to build an API using Go and Fiber", slug.GenerateSlugId()),
	},

	{
		ID:          3,
		Title:       "Como eu penso o JPA?",
		PublishedAt: time.Now().Format("2006.01.02"),
		Content:     `O JPA é uma especificação que define uma interface comum para frameworks de mapeamento objeto-relacional.`,
		Slug:	slug.GenerateSlug("My math roadmap: from zero to hero", slug.GenerateSlugId()),
	},
}

func (r *router) Home(c *fiber.Ctx) error {

	return c.Render("pages/home", fiber.Map{
		"Posts": Posts,
		"PageTitle": "home",
	})
}

func (r *router) Post(c *fiber.Ctx) error {
	var post PostOutput

	for _, p := range Posts {
		if p.Slug == c.Params("slug") {
			post = p
			break
		}
	}

	return c.Render("pages/post", fiber.Map{
		"Post": post,
		"PageTitle": post.Slug,
	})
}

func (r *router) LoginPage(c *fiber.Ctx) error {

	session, err := r.store.Get(c)

	if err != nil {
		log.Default().Printf("Error getting session: %v", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	log.Default().Printf("Session: %v", session.Get("is_logged_in"))

	if session.Get("is_logged_in") != nil {
		log.Default().Println("User is already logged in")
		return c.Redirect("/admin")
	}	

	return c.Render("pages/login", fiber.Map{
		"PageTitle": "login",
	})
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
			"error": err.Error(),
		})
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"error": errors.New("wrong password").Error(),
		})
	}

	session, err := r.store.Get(c)

	if err != nil {
		log.Default().Printf("Error getting session: %v", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	session.Set("username", username)
	session.Set("is_admin", user.IsAdmin)
	session.Set("is_logged_in", true)

	err = session.Save()

	if err != nil {
		log.Default().Printf("Error saving session: %v", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Redirect("/admin/posts")
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