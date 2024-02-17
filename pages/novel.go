package pages

import (
	"errors"
	"strconv"

	core "codeberg.org/vnpower/pixivfe/v2/core/webapi"
	"github.com/gofiber/fiber/v2"
)

func NovelPage(c *fiber.Ctx) error {
	param_id := c.Params("id")
	if _, err := strconv.Atoi(param_id); err != nil {
		return errors.New("invalid id")
	}

	novel, err := core.GetNovelByID(c, param_id)
	if err != nil {
		return err
	}

	// todo: passing ArtWorkData{} here will not work. maybe lowercase?
	return c.Render("pages/novel", fiber.Map{
		"Novel": novel,
		"Title": novel.Title,
	})
}
