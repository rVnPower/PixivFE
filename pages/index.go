package pages

import (
	session "codeberg.org/vnpower/pixivfe/v2/core/user"
	core "codeberg.org/vnpower/pixivfe/v2/core/webapi"
	"github.com/gofiber/fiber/v2"
)

func IndexPage(c *fiber.Ctx) error {

	// If token is set, do the landing request...
	if token := session.GetToken(c); token != "" {
		mode := c.Query("mode", "all")

		works, err := core.GetLanding(c, mode)

		if err != nil {
			return err
		}

		return c.Render("pages/index", fiber.Map{
			"Title": "Landing", "Data": works,
		})
	}

	// ...otherwise, default to today's illustration ranking
	works, err := core.GetRanking(c, "daily", "illust", "", "1")
	if err != nil {
		return err
	}
	return c.Render("pages/index", fiber.Map{
		"Title": "Landing", "NoTokenData": works,
	})
}
