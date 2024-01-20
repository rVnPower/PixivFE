package pages

import (
	"strconv"
	"strings"

	session "codeberg.org/vnpower/pixivfe/v2/core/config"
	core "codeberg.org/vnpower/pixivfe/v2/core/webapi"
	"github.com/gofiber/fiber/v2"
)

func LoginUserPage(c *fiber.Ctx) error {
	token := session.GetToken(c)

	if token == "" {
		return c.Redirect("/login")
	}

	// The left part of the token is the member ID
	userId := strings.Split(token, "_")

	c.Redirect("/users/" + userId[0])
	return nil
}

func LoginBookmarkPage(c *fiber.Ctx) error {
	token := session.GetToken(c)
	if token == "" {
		return c.Redirect("/login")
	}

	// The left part of the token is the member ID
	userId := strings.Split(token, "_")

	c.Redirect("/users/" + userId[0] + "/bookmarks#checkpoint")
	return nil
}

func FollowingWorksPage(c *fiber.Ctx) error {
	if token := session.GetToken(c); token == "" {
		return c.Redirect("/login")
	}

	mode := c.Query("mode", "all")
	page := c.Query("page", "1")

	pageInt, _ := strconv.Atoi(page)

	works, err := core.GetNewestFromFollowing(c, mode, page)
	if err != nil {
		return err
	}

	return c.Render("pages/following", fiber.Map{
		"Title":    "Following works",
		"Mode":     mode,
		"Artworks": works,
		"CurPage":  page,
		"Page":     pageInt,
	})
}
