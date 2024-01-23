package pages

import (
	core "codeberg.org/vnpower/pixivfe/v2/core/webapi"
	"github.com/gofiber/fiber/v2"
)

func IndexPage(c *fiber.Ctx) error {
	works, err := core.GetRanking(c, "daily", "illust", "", "1")
	if err != nil {
		return err
	}

	return c.Render("pages/index", fiber.Map{
		"Title": "Landing", "Token": false, "Data": works,
	})
}
