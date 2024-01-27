package pages

import (
	"codeberg.org/vnpower/pixivfe/v2/core/config"
	"github.com/gofiber/fiber/v2"
)

func AboutPage(c *fiber.Ctx) error {
	info := fiber.Map{
		"Time":           core.GlobalServerConfig.StartingTime,
		"Version":        core.GlobalServerConfig.Version,
		"ImageProxy":     core.GlobalServerConfig.ProxyServer,
		"AcceptLanguage": core.GlobalServerConfig.AcceptLanguage,
	}
	return c.Render("pages/about", info)
}
