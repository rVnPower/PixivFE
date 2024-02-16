package pages

import (
	"fmt"
	"io"
	"net/http"
	config "codeberg.org/vnpower/pixivfe/v2/core/config"

	"github.com/gofiber/fiber/v2"
)

func SPximgProxy(c *fiber.Ctx) error {
	URL := fmt.Sprintf("https://s.pximg.net/%s", c.Params("*"))
	req, _ := http.NewRequest("GET", URL, nil)

	// Make the request
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	c.Set("Content-Type", resp.Header.Get("Content-Type"))

	return c.Send([]byte(body))
}

func IPximgProxy(c *fiber.Ctx) error {
	proxy_authority := config.GetImageProxyAuthority(c)
	URL := fmt.Sprintf("https://%s/%s", proxy_authority, c.Params("*"))
	req, _ := http.NewRequest("GET", URL, nil)

	// Make the request
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	c.Set("Content-Type", resp.Header.Get("Content-Type"))

	return c.Send([]byte(body))
}

func UgoiraProxy(c *fiber.Ctx) error {
	URL := fmt.Sprintf("https://ugoira.com/api/mp4/%s", c.Params("*"))
	req, _ := http.NewRequest("GET", URL, nil)

	// Make the request
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	c.Set("Content-Type", resp.Header.Get("Content-Type"))

	return c.Send([]byte(body))
}
