package pages

import (
	"github.com/gofiber/fiber/v2"
)

func IndexPage(c *fiber.Ctx) error {
	return c.Render("pages/index", fiber.Map{"Title": "Landing", "Token": false})
}
