package routes

import (
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/samluiz/blog/api/types"
	"github.com/samluiz/blog/pkg/user"
	"golang.org/x/crypto/bcrypt"
)

const IS_LOGGED = "is_logged"
const DASHBOARD_URL = "/blog/admin/dashboard"
const DEV_TO_API_BASE_URL = "https://dev.to/api"
const DEV_TO_USERNAME = "samluiz"
const DATE_LAYOUT = "2006-01-02T15:04:05.999Z"

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

func preprocessHTML(html string) template.HTML {
	return template.HTML(html)
}

func (r *router) HomePage(c *fiber.Ctx) error {
	articles, err := getArticlesFromDevTo()

	if err != nil {
		log.Default().Println(err.Error())
	}

	return c.Render("pages/home", fiber.Map{
		"Articles":  articles,
		"PageTitle": "home",
		"Error":     err,
	})
}

func (r *router) ArticlePage(c *fiber.Ctx) error {
	slug := c.Params("slug")

	article, err := getArticleBySlug(slug)

	if err != nil {
		log.Default().Println(err.Error())
	}

	htmlContent := preprocessHTML(article.BodyHTML)

	return c.Render("pages/article", fiber.Map{
		"Article":   article,
		"HTML":      htmlContent,
		"PageTitle": article.Slug,
		"Error":     err,
	})
}

func getArticleBySlug(slug string) (*types.ArticleResponse, error) {
	var getArticleResponse types.GetArticleByPathResponse
	var articleResponse types.ArticleResponse

	log.Default().Println("Getting article from dev.to")

	request := fiber.Get(DEV_TO_API_BASE_URL + "/articles/" + DEV_TO_USERNAME + "/" + slug)

	status, response, err := request.Bytes()

	log.Default().Printf("Status: %v", status)

	if (status != 200) || (err != nil) {
		log.Default().Printf("Error: %v", err)
		return nil, errors.New("error getting article from dev.to: " + string(response))
	}

	jsonErr := json.Unmarshal(response, &getArticleResponse)
	if jsonErr != nil {
		return nil, jsonErr
	}

	getArticleResponse.PublishedAt = formatDate(getArticleResponse.PublishedAt)
	articleResponse = types.ArticleResponse(getArticleResponse)

	return &articleResponse, nil
}

func getArticlesFromDevTo() ([]types.ArticleResponse, error) {
	var articles []types.GetArticleByPathResponse
	articlesResponse := make([]types.ArticleResponse, len(articles))

	log.Default().Println("Getting articles from dev.to")

	request := fiber.Get(DEV_TO_API_BASE_URL + "/articles/me/published")
	request.Set("api-key", os.Getenv("DEV_TO_API_KEY"))
	request.Request().URI().SetQueryString("page=1&per_page=5")

	status, response, err := request.Bytes()

	log.Default().Printf("Status: %v", status)
	log.Default().Printf("Error: %v", err)

	if (status != 200) || (err != nil) {
		log.Default().Printf("Error: %v", err)
		return nil, errors.New("error getting articles from dev.to: " + string(response))
	}

	jsonErr := json.Unmarshal(response, &articles)
	if jsonErr != nil {
		return nil, jsonErr
	}

	for _, a := range articles {
		a.PublishedAt = formatDate(a.PublishedAt)
		articlesResponse = append(articlesResponse, types.ArticleResponse(a))
	}

	return articlesResponse, nil
}

func formatDate(date string) string {
	parsedDate, err := time.Parse(DATE_LAYOUT, date)
	if err != nil {
		log.Default().Printf("Error parsing time: %v", err)
	}
	return parsedDate.Format("2006.01.02")
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
