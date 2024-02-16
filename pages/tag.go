package pages

import (
	"net/url"
	"strconv"

	core "codeberg.org/vnpower/pixivfe/v2/core/webapi"
	"github.com/gofiber/fiber/v2"
)

func TagPage(c *fiber.Ctx) error {
	queries := make(map[string]string, 3)
	queries["Mode"] = c.Query("mode", "safe")
	queries["Category"] = c.Query("category", "artworks")
	queries["Order"] = c.Query("order", "date_d")

	name, err := url.PathUnescape(c.Params("name"))
	if err != nil {
		return err
	}

	page := c.Query("page", "1")
	pageInt, _ := strconv.Atoi(page)

	tag, err := core.GetTagData(c, name)
	if err != nil {
		return err
	}
	result, err := core.GetSearch(c, queries["Category"], name, queries["Order"], queries["Mode"], page)
	if err != nil {
		return err
	}

	return c.Render("pages/tag", fiber.Map{"Title": "Results for " + tag.Name, "Tag": tag, "Data": result, "Queries": queries, "Page": pageInt})
}
