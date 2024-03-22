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
	page int
}

func fetchData(c *fiber.Ctx, getTags bool) (userPageData, error) {
	id := c.Params("id")
	if _, err := strconv.Atoi(id); err != nil {
		return userPageData{}, err
	}
	category := core.UserArtCategory(c.Params("category", string(core.UserArt_Any)))
	err := category.Validate()
	if err != nil {
		return userPageData{}, err
	}

	page_param := c.Query("page", "1")
	page, err := strconv.Atoi(page_param)
	if err != nil {
		return userPageData{}, err
	}

	user, err := core.GetUserArtwork(c, id, category, page, getTags)
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

	return userPageData{user, category, pageLimit, page}, nil
}

func UserPage(c *fiber.Ctx) error {
	data, err := fetchData(c, true)
	if err != nil {
		return err
	}

	return c.Render("pages/user", fiber.Map{
		"Title":     data.user.Name,
		"User":      data.user,
		"Category":  data.category,
		"PageLimit": data.pageLimit,
		"Page":      data.page,
		"MetaImage": data.user.BackgroundImage,
	})
}

func UserAtomFeed(c *fiber.Ctx) error {
	data, err := fetchData(c, false)
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
		"Page":      data.page,
		// "MetaImage": data.user.BackgroundImage,
	}, "")
	if err != nil {
		return err
	}

	c.Context().SetContentType("application/atom+xml")

	return nil
}
