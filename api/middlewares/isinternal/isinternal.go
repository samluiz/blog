package isinternal

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func New() fiber.Handler {
	return func(c *fiber.Ctx) error {
		domain := string(c.Request().Host())
		host := c.Hostname()

		fmt.Printf("domain: %v, host: %v\n", domain, host)

		if domain != host {
			log.Default().Printf("request to internal route from external domain: %v", domain)
			return c.Redirect("/error/404")
		}

		return c.Next()
	}
}
