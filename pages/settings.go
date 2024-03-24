package pages

import (
	"errors"
	"io"
	"net/http"
	"regexp"

	httpc "codeberg.org/vnpower/pixivfe/v2/core/http"
	session "codeberg.org/vnpower/pixivfe/v2/core/session"
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
	} else {
		session.ClearCookie(c, session.Cookie_ImageProxy)
	}
	return nil
}

func setNovelFontType(c *fiber.Ctx) error {
	fontType := c.FormValue("font-type")
	if fontType != "" {
		session.SetCookie(c, session.Cookie_NovelFontType, fontType)
	}

	return nil
}

func setLogout(c *fiber.Ctx) error {
	session.ClearCookie(c, session.Cookie_Token)
	return nil
}

func resetAll(c *fiber.Ctx) error {
	session.ClearAllCookies(c)
	return nil
}

func SettingsPage(c *fiber.Ctx) error {
	cookies := []fiber.Map{}
	for _, name := range session.AllCookieNames {
		value := session.GetCookie(c, name)
		cookies = append(cookies, fiber.Map{
			"Key":   name,
			"Value": value,
		})
	}
	return c.Render("pages/settings", fiber.Map{
		"CookieList": cookies,
		"ProxyList":  doc.BuiltinProxyList,
	})
}

func SettingsPost(c *fiber.Ctx) error {
	t := c.Params("type")
	var err error

	switch t {
	case "imageServer":
		err = setImageServer(c)
	case "token":
		err = setToken(c)
	case "logout":
		err = setLogout(c)
	case "reset-all":
		err = resetAll(c)
	case "novelFontType":
		err = setNovelFontType(c)
	default:
		err = errors.New("No such setting is available.")
	}

	if err != nil {
		return err
	}

	ret := c.Query("redirect", "/")
	c.Redirect(ret)
	return nil
}
