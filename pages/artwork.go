package pages

import (
	"fmt"
	"strconv"

	core "codeberg.org/vnpower/pixivfe/v2/core/webapi"
	"github.com/gofiber/fiber/v2"
)

func ArtworkPage(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := strconv.Atoi(id); err != nil {
		return fmt.Errorf("Invalid ID: %s", id)
	}

	illust, err := core.GetArtworkByID(c, id, true)
	if err != nil {
		return err
	}

	metaDescription := ""
	for _, i := range illust.Tags {
		metaDescription += "#" + i.Name + ", "
	}

	// todo: passing ArtWorkData{} here will not work. maybe lowercase?
	return c.Render("pages/artwork", fiber.Map{
		"Illust":          illust,
		"Title":           illust.Title,
		"MetaDescription": metaDescription,
		"MetaImage":       illust.Images[0].Original,
	})
}

func ArtworkEmbedPage(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := strconv.Atoi(id); err != nil {
		return fmt.Errorf("Invalid ID: %s", id)
	}

	illust, err := core.GetArtworkByID(c, id, false)
	if err != nil {
		return err
	}

	metaDescription := ""
	for _, i := range illust.Tags {
		metaDescription += "#" + i.Name + ", "
	}

	// todo: passing ArtWorkData{} here will not work. maybe lowercase?
	return c.Render("embed", fiber.Map{
		"Illust":          illust,
		"Title":           illust.Title,
		"MetaDescription": metaDescription,
		"MetaImage":       illust.Images[0].Original,
	}, "embed")
}
