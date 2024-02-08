package pages

import (
	"errors"
	"io"
	"net/http"
	"regexp"

	session "codeberg.org/vnpower/pixivfe/v2/core/config"
	httpc "codeberg.org/vnpower/pixivfe/v2/core/http"
	"github.com/gofiber/fiber/v2"
)

func setToken(c *fiber.Ctx) error {
	// Parse the value from the form
	token := c.FormValue("token")
	if token != "" {
		_, err := httpc.UnwrapWebAPIRequest(httpc.GetNewestFromFollowingURL("all", "1"), token)
		if err != nil {
			return errors.New("Cannot authorize with supplied token.")
		}

		// Make a test request to verify the token.
		// THE TEST URL IS NSFW!
		req, _ := http.NewRequest("GET", "https://www.pixiv.net/en/artworks/115365120", nil)
		req.Header.Add("User-Agent", "Mozilla/5.0")
		req.AddCookie(&http.Cookie{
			Name:  "PHPSESSID",
			Value: token,
		})

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return errors.New("Cannot authorize with supplied token.")
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return errors.New("Cannot parse the response from Pixiv. Please report this issue.")
		}

		// CSRF token
		r := regexp.MustCompile(`"token":"([0-9a-f]+)"`)
		csrf := r.FindStringSubmatch(string(body))[1]

		if csrf == "" {
			return errors.New("Cannot authorize with supplied token.")
		}

		// Set the tokens
		if err := session.SetSessionValue(c, "Token", token); err != nil {
			return err
		}
		if err := session.SetSessionValue(c, "CSRF", csrf); err != nil {
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
	c.Redirect("/")
	return nil
}
