package islogged

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/samluiz/blog/api/routes"
)

func New(config Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Default().Printf("validating session...")

		session, err := config.Session.Get(c)

		if err != nil {
			log.Default().Printf("error retrieving the session: %v", err)
			return c.Redirect("/")
		}

		isLogged := session.Get(routes.IS_LOGGED)

		if isLogged == nil || isLogged == false {
			log.Default().Printf("user is not logged in. redirecting to login...")
			return c.Redirect("/blog/login")
		}

		log.Default().Println("user is logged in")

		return c.Next()
	}
}
