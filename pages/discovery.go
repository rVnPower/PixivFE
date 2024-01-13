package pages

import (
	core "codeberg.org/vnpower/pixivfe/v2/core/webapi"
	"github.com/gofiber/fiber/v2"
)

func DiscoveryPage(c *fiber.Ctx) error {
	mode := c.Query("mode", "safe")

	works, err := core.GetDiscoveryArtwork(c, mode)
	if err != nil {
		return err
	}

	return c.Render("pages/discovery", fiber.Map{
		"Artworks": works,
		"Title":    "Discovery",
	})
}
