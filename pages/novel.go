package pages

import (
	"fmt"
	"strconv"

	core "codeberg.org/vnpower/pixivfe/v2/core/webapi"
	"github.com/gofiber/fiber/v2"
)

func NovelPage(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := strconv.Atoi(id); err != nil {
		return fmt.Errorf("Invalid ID: %s", id)
	}

	novel, err := core.GetNovelByID(c, id)
	if err != nil {
		return err
	}

	// todo: passing ArtWorkData{} here will not work. maybe lowercase?
	return c.Render("pages/novel", fiber.Map{
		"Novel": novel,
		"Title": novel.Title,
	})
}
