package pages

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	core "codeberg.org/vnpower/pixivfe/v2/core/webapi"
	"github.com/gofiber/fiber/v2"
)

func ArtworkMultiPage(c *fiber.Ctx) error {
	param_ids := c.Params("ids")
	ids := strings.Split(param_ids, ",")

	artworks := make([]*core.Illust, len(ids))

	wg := sync.WaitGroup{}
	wg.Add(len(ids))
	for i, id := range ids {
		if _, err := strconv.Atoi(id); err != nil {
			return errors.New("invalid id")
		}

		go func(i int, id string) {
			defer wg.Done()

			illust, err := core.GetArtworkByID(c, id, false)
			if err != nil {
				artworks[i] = &core.Illust{
					Title: err.Error(), // this might be flaky
				}
				return
			}

			metaDescription := ""
			for _, i := range illust.Tags {
				metaDescription += "#" + i.Name + ", "
			}

			artworks[i] = illust
		}(i, id)
	}
	wg.Wait()

	return c.Render("pages/artwork-multi", fiber.Map{
		"Artworks": artworks,
		"Title":    fmt.Sprintf("(%d images)", len(artworks)),
	})
}
