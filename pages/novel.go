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

	related, err := core.GetNovelRelated(c, id)
	if err != nil {
		return err
	}

	user, err := core.GetUserBasicInformation(c, novel.UserID)
	if err != nil {
		return err
	}

	return c.Render("pages/novel", fiber.Map{
		"Novel":        novel,
		"NovelRelated": related,
		"User":         user,
		"Title":        novel.Title,
	})
}
