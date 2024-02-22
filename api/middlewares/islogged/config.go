package islogged

import (
	"github.com/gofiber/fiber/v2/middleware/session"
)

type Config struct {
	Session *session.Store
}
