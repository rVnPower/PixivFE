package pages

import (
	"errors"

	session "codeberg.org/vnpower/pixivfe/v2/core/config"
	http "codeberg.org/vnpower/pixivfe/v2/core/http"
	"github.com/gofiber/fiber/v2"
)

func setToken(c *fiber.Ctx) error {
	// Parse the value from the form
	token := c.FormValue("token")
	if token != "" {
		if _, err := http.UnwrapWebAPIRequest("https://www.pixiv.net/ajax/user/extra", token); err != nil {
			return errors.New("Cannot authorize with supplied token.")
		}
		if err := session.SetSessionValue(c, "Token", token); err != nil {
			return err
		}

		return nil
	}
	return errors.New("You submitted an empty/invalid form.")
}

func setImageServer(c *fiber.Ctx) error {
	// Parse the value from the form
	token := c.FormValue("image-proxy")
	if token != "" {
		if err := session.SetSessionValue(c, "ImageProxy", token); err != nil {
			return err
		}

		return nil
	}
	return errors.New("You submitted an empty/invalid form.")
}

func setLogout(c *fiber.Ctx) error {
	session.RemoveSessionValue(c, "Token")
	return nil
}

func LoginPage(c *fiber.Ctx) error {
	return c.Render("pages/login", fiber.Map{})
}

func SettingsPage(c *fiber.Ctx) error {
	return c.Render("pages/settings", fiber.Map{})
}

func SettingsPost(c *fiber.Ctx) error {
	t := c.Params("type")
	var err error

	switch t {
	case "image_server":
		err = setImageServer(c)
	case "token":
		err = setToken(c)
	case "logout":
		err = setLogout(c)
	default:
		err = errors.New("No such methods available.")
	}

	if err != nil {
		return err
	}
	c.Redirect("/settings")
	return nil
}
