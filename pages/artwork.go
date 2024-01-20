package pages

import (
	"errors"
	"strconv"

	core "codeberg.org/vnpower/pixivfe/v2/core/webapi"
	"github.com/gofiber/fiber/v2"
)

func ArtworkPage(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := strconv.Atoi(id); err != nil {
		return errors.New("Invalid ID.")
	}

	illust, err := core.GetArtworkByID(c, id)
	if err != nil {
		return err
	}

	// Optimize this
	return c.Render("pages/artwork", fiber.Map{
		"Illust":   illust,
		"Title":    illust.Title,
		"PageType": "artwork",
	})
}
