package pages

import (
	"math"
	"strconv"
	"time"

	core "codeberg.org/vnpower/pixivfe/v2/core/webapi"
	"github.com/gofiber/fiber/v2"
)

type userPageData struct {
	user        core.User
	category    core.UserArtCategory
	pageLimit   int
	pageCurrent int
}

func process(c *fiber.Ctx) (userPageData, error) {
	id := c.Params("id")
	if _, err := strconv.Atoi(id); err != nil {
		return userPageData{}, err
	}
	category := core.UserArtCategory(c.Params("category", string(core.UserArt_Any)))
	err := category.Validate()
	if err != nil {
		return userPageData{}, err
	}

	pageCurrentString := c.Query("page", "1")
	pageCurrent, err := strconv.Atoi(pageCurrentString)
	if err != nil {
		return userPageData{}, err
	}

	user, err := core.GetUserArtwork(c, id, category, pageCurrent)
	if err != nil {
		return userPageData{}, err
	}

	var worksCount int
	var worksPerPage float64

	if category == core.UserArt_Bookmarked {
		worksPerPage = 48.0
	} else {
		worksPerPage = 30.0
	}

	worksCount = user.ArtworksCount
	pageLimit := int(math.Ceil(float64(worksCount) / worksPerPage))

	return userPageData{user, category, pageLimit, pageCurrent}, nil
}

func UserPage(c *fiber.Ctx) error {
	data, err := process(c)
	if err != nil {
		return err
	}

	return c.Render("pages/user", fiber.Map{
		"Title":     data.user.Name,
		"User":      data.user,
		"Category":  data.category,
		"PageLimit": data.pageLimit,
		"Page":      data.pageCurrent,
		"MetaImage": data.user.BackgroundImage,
	})
}

func UserAtomFeed(c *fiber.Ctx) error {
	data, err := process(c)
	if err != nil {
		return err
	}

	err = c.Render("pages/user.atom", fiber.Map{
		"URL":       string(c.Request().RequestURI()),
		"Title":     data.user.Name,
		"User":      data.user,
		"Category":  data.category,
		"Updated":   time.Now().Format(time.RFC3339),
		"PageLimit": data.pageLimit,
		"Page":      data.pageCurrent,
		// "MetaImage": data.user.BackgroundImage,
	}, "")
	if err != nil {
		return err
	}

	c.Context().SetContentType("application/atom+xml")

	return nil
}
