package pages

import (
	"errors"
	"strconv"

	core "codeberg.org/vnpower/pixivfe/v2/core/webapi"
	"github.com/gofiber/fiber/v2"
)

func ArtworkPage(c *fiber.Ctx) error {
	param_id := c.Params("id")
	if _, err := strconv.Atoi(param_id); err != nil {
		return errors.New("invalid id")
	}

	illust, err := core.GetArtworkByID(c, param_id, true)
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
	param_id := c.Params("id")
	if _, err := strconv.Atoi(param_id); err != nil {
		return errors.New("invalid id")
	}

	illust, err := core.GetArtworkByID(c, param_id, false)
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
