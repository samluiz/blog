package routes

import (
	"errors"
	"html/template"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/samluiz/blog/api/integrations"
	"github.com/samluiz/blog/api/parsers"
	apiTypes "github.com/samluiz/blog/api/types"
	"github.com/samluiz/blog/common/logger"
	"github.com/samluiz/blog/common/providers"
	"github.com/samluiz/blog/pkg/types"
	"github.com/samluiz/blog/pkg/user"
	"golang.org/x/crypto/bcrypt"
)

const IS_LOGGED = "is_logged"
const DASHBOARD_URL = "/dashboard"

var LOGGER = logger.New(os.Stdout, logger.DebugLevel, "[ROUTER]")

type Router interface {
	HomePage(c *fiber.Ctx) error
	ArticlePage(c *fiber.Ctx) error
	ArticlesPage(c *fiber.Ctx) error
	LoginPage(c *fiber.Ctx) error
	AdminDashboardPage(c *fiber.Ctx) error
	AdminArticlesPartial(c *fiber.Ctx) error
	Authenticate(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	GithubCallback(c *fiber.Ctx) error
	NotFoundPage(c *fiber.Ctx) error
	ErrorPage(c *fiber.Ctx) error
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
	articles, err := integrations.GetArticlesFromDevTo(1, 3)

	if err != nil {
		LOGGER.Error(err.Error())
	}

	session, err := r.store.Get(c)

	if err != nil {
		LOGGER.Error("error getting session: %v", err)
	}

	isLogged := session.Get(IS_LOGGED)
	user := session.Get("user")

	return c.Render("pages/home", fiber.Map{
		"Articles":    articles,
		"IsLogged":    isLogged,
		"User":        user,
		"PageTitle":   "home",
		"Description": "My personal portfolio, but also a blog about software development, programming, and technology. Articles about web development, backend, frontend, and whatever i wanna share.",
		"Route":       "",
		"Error":       err,
	})
}

func (r *router) ArticlePage(c *fiber.Ctx) error {
	slug := c.Params("slug")

	article, err := integrations.GetArticleBySlugDevTo(slug)

	if err != nil {
		LOGGER.Error(err.Error())
	}

	markdownContent := template.HTML(parsers.MarkdownToHTML([]byte(article.BodyMarkdown)))

	session, err := r.store.Get(c)

	if err != nil {
		LOGGER.Error("error getting session: %v", err)
	}

	isLogged := session.Get(IS_LOGGED)
	user := session.Get("user")

	return c.Render("pages/article", fiber.Map{
		"Article":     article,
		"IsLogged":    isLogged,
		"User":        user,
		"Markdown":    markdownContent,
		"PageTitle":   article.Slug,
		"Description": article.Description,
		"Route":       "articles/" + article.Slug,
		"Error":       err,
	})
}

func (r *router) ArticlesPage(c *fiber.Ctx) error {
	articles, err := integrations.GetArticlesFromDevTo(1, 10)

	if err != nil {
		LOGGER.Error(err.Error())
	}

	session, err := r.store.Get(c)

	if err != nil {
		LOGGER.Error("error getting session: %v", err)
	}

	isLogged := session.Get(IS_LOGGED)
	user := session.Get("user")

	return c.Render("pages/articles", fiber.Map{
		"Articles":    articles,
		"IsLogged":    isLogged,
		"User":        user,
		"PageTitle":   "articles",
		"Description": "Articles about web development, backend, frontend, and whatever i wanna share.",
		"Route":       "articles",
		"Error":       err,
	})
}

func (r *router) LoginPage(c *fiber.Ctx) error {

	redirect := c.Query("redirect")

	if redirect == "" || redirect == "/auth/login" {
		redirect = "/"
	}

	session, err := r.store.Get(c)

	if err != nil {
		LOGGER.Error("error getting session: %v", err)
		return fiber.ErrInternalServerError
	}

	isLogged := session.Get(IS_LOGGED)

	if isLogged != nil && isLogged == true {
		LOGGER.Info("user is already logged in. redirecting to dashboard.")
		return c.Redirect(DASHBOARD_URL)
	}

	githubUrl := integrations.GetGithubAuthURL()

	return c.Render("pages/login", fiber.Map{
		"PageTitle": "login",
		"GithubURL": githubUrl,
		"Redirect":  redirect,
	})
}

func (r *router) AdminDashboardPage(c *fiber.Ctx) error {
	session, err := r.store.Get(c)

	if err != nil {
		LOGGER.Error("error getting session: %v", err)
	}

	isLogged := session.Get(IS_LOGGED)
	user := session.Get("user")

	return c.Render("pages/dashboard", fiber.Map{
		"IsLogged": isLogged,
		"User":     user,
	})
}

func (r *router) AdminArticlesPartial(c *fiber.Ctx) error {

	return c.SendFile("views/partials/articles.html")
}

func (r *router) Authenticate(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	const WRONG_CREDENTIALS = "Wrong credentials. Please try again."
	const UNKNOWN_ERROR = "Something went wrong. Please try again."

	user, err := r.userService.FindUserByUsername(username)

	if err != nil {
		LOGGER.Error(err.Error())
		if errors.Is(err, types.ErrUserNotFound) {
			LOGGER.Debug("user: %s not found", username)
			return c.SendString(WRONG_CREDENTIALS)
		}
		return c.SendString(UNKNOWN_ERROR)
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return c.SendString(WRONG_CREDENTIALS)
	}

	session, err := r.store.Get(c)

	if err != nil {
		LOGGER.Error("error getting session: %v", err)
		return c.SendString(UNKNOWN_ERROR)
	}

	sessionUser := apiTypes.SessionUser{
		ID:       user.ID,
		Username: username,
		IsAdmin:  user.IsAdmin,
		Avatar:   user.Avatar,
	}
	session.Set("user", sessionUser)
	session.Set(IS_LOGGED, true)

	err = session.Save()

	if err != nil {
		LOGGER.Error("error saving session: %v", err)
		return c.SendString(UNKNOWN_ERROR)
	}

	res := c.Response()
	res.Header.Add("HX-Redirect", c.Get("X-Redirect"))

	return c.SendStatus(fiber.StatusOK)
}

func (r *router) Logout(c *fiber.Ctx) error {
	session, err := r.store.Get(c)

	if err != nil {
		LOGGER.Error("error retrieving session: %v", err)
		return c.Redirect("/")
	}

	session.Destroy()

	return c.Redirect("/")
}

func (r *router) GithubCallback(c *fiber.Ctx) error {
	code := c.Query("code")

	if code == "" {
		return c.Redirect("/auth/login")
	}

	githubResponse, err := integrations.ExchangeGithubToken(code)

	if err != nil {
		LOGGER.Error(err.Error())
		return c.Redirect("/auth/login")
	}

	session, err := r.store.Get(c)

	if err != nil {
		LOGGER.Error("error getting session: %v", err)
		return fiber.ErrInternalServerError
	}

	userInfo, err := integrations.GetGithubUserInfo(githubResponse.AccessToken)

	if err != nil {
		LOGGER.Error(err.Error())
		return c.Redirect("/auth/login")
	}

	user, err := r.userService.FindExternalUserByProviderId(userInfo.ID, providers.GITHUB)

	if err != nil {
		if errors.Is(err, types.ErrUserNotFound) {
			LOGGER.Info("user not found. creating user...")
			newUser := &types.CreateExternalUserInput{
				ProviderId: userInfo.ID,
				Name:       userInfo.Name,
				Username:   userInfo.Login,
				Provider:   providers.GITHUB,
				Avatar:     userInfo.AvatarURL,
			}
			user, err = r.userService.SaveUser(newUser)
			if err != nil {
				LOGGER.Error(err.Error())
				return c.Redirect("/auth/login")
			}
		} else {
			LOGGER.Error(err.Error())
			return c.Redirect("/auth/login")
		}
	}

	sessionUser := apiTypes.SessionUser{
		ID:       user.ID,
		Username: user.Username,
		IsAdmin:  false,
		Provider: user.Provider,
		Avatar:   user.Avatar,
	}
	session.Set("user", sessionUser)
	session.Set("github_token", githubResponse.AccessToken)
	session.Set(IS_LOGGED, true)

	err = session.Save()

	if err != nil {
		LOGGER.Error("error saving session: %v", err)
		return c.Redirect("/auth/login")
	}

	return c.Redirect("/")
}

func (r *router) NotFoundPage(c *fiber.Ctx) error {
	return c.Render("pages/not-found", fiber.Map{
		"PageTitle": "404",
	})
}

func (r *router) ErrorPage(c *fiber.Ctx) error {
	httpStatus := c.Params("status")

	return c.Render("pages/error", fiber.Map{
		"PageTitle":  "error",
		"HttpStatus": httpStatus,
	})
}
