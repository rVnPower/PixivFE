package pages

import (
	"errors"
	"io"
	"net/http"
	"regexp"

	session "codeberg.org/vnpower/pixivfe/v2/core/user"
	httpc "codeberg.org/vnpower/pixivfe/v2/core/http"
	"codeberg.org/vnpower/pixivfe/v2/doc"
	"github.com/gofiber/fiber/v2"
)

// todo: allow clear proxy
// todo: allow clear all settings

func setToken(c *fiber.Ctx) error {
	// Parse the value from the form
	token := c.FormValue("token")
	if token != "" {
		URL := httpc.GetNewestFromFollowingURL("all", "1")

		_, err := httpc.UnwrapWebAPIRequest(c.Context(), URL, token)
		if err != nil {
			return errors.New("Cannot authorize with supplied token.")
		}

		// Make a test request to verify the token.
		// THE TEST URL IS NSFW!
		req, err := http.NewRequest("GET", "https://www.pixiv.net/en/artworks/115365120", nil)
		if err != nil {
			return err
		}
		req = req.WithContext(c.Context())
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

		// Set the token
		session.SetCookie(c, session.Cookie_Token, token)
		session.SetCookie(c, session.Cookie_CSRF, csrf)

		return nil
	}
	return errors.New("You submitted an empty/invalid form.")
}

func setImageServer(c *fiber.Ctx) error {
	// Parse the value from the form
	token := c.FormValue("image-proxy")
	if token != "" {
		session.SetCookie(c, session.Cookie_ImageProxy, token)
		return nil
	}
	return errors.New("You submitted an empty/invalid form.")
}

func setLogout(c *fiber.Ctx) error {
	session.ClearCookie(c, session.Cookie_Token)
	return nil
}

func SettingsPage(c *fiber.Ctx) error {
	return c.Render("pages/settings", fiber.Map{
		"ProxyList": doc.BuiltinProxyList,
	})
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
		err = errors.New("no such setting available")
	}

	if err != nil {
		return err
	}
	c.Redirect("/")
	return nil
}
