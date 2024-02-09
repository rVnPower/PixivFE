package pages

import (
	"errors"
	"log"
	"strconv"
	"strings"

	core "codeberg.org/vnpower/pixivfe/v2/core/webapi"
	"github.com/gofiber/fiber/v2"
)

func ArtworkMultiPage(c *fiber.Ctx) error {
	ids := c.Params("ids")

	artworks := []fiber.Map{}

	// optimize: sequential access might be slow
	for _, id := range strings.Split(ids, ",") {
		if _, err := strconv.Atoi(id); err != nil {
			return errors.New("Invalid ID.")
		}

		illust, err := core.GetArtworkByID(c, id)
		if err != nil {
			return err
		}

		metaDescription := ""
		for _, i := range illust.Tags {
			metaDescription += "#" + i.Name + ", "
		}

		artworks = append(artworks, fiber.Map{
			"Illust":          illust,
			"Title":           illust.Title,
			"PageType":        "artwork",
			"MetaDescription": metaDescription,
			"MetaImage":       illust.Images[0].Original,
		})
	}

	log.Println("artworks:", artworks)

	return c.Render("pages/artwork-multi", fiber.Map{
		"Artworks": artworks,
	})
}
