package pages

import (
	"errors"
	"math"
	"strconv"

	core "codeberg.org/vnpower/pixivfe/v2/core/webapi"
	"github.com/gofiber/fiber/v2"
)

func UserPage(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := strconv.Atoi(id); err != nil {
		return err
	}
	category := core.UserArtCategory(c.Params("category", string(core.UserArt_Any)))
	err := category.Validate()
	if err != nil {
		return err
	}

	page := c.Query("page", "1")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return err
	}

	user, err := core.GetUserArtwork(c, id, category, pageInt)
	if err != nil {
		return err
	}

	var worksCount int
	var worksPerPage float64

	if category == core.UserArt_Bookmarked {
		worksPerPage = 48.0
	} else {
		worksPerPage = 30.0
	}

	worksCount = user.ArtworksCount
	pageLimit := math.Ceil(float64(worksCount) / worksPerPage)

	return c.Render("pages/user", fiber.Map{
		"Title":     user.Name,
		"User":      user,
		"Category":  category,
		"PageLimit": int(pageLimit),
		"Page":      pageInt,
		"MetaImage": user.BackgroundImage,
	})
}

func UserFeed(c *fiber.Ctx) error {
	return errors.New("not implemented")
}
