package pages

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func SPximgProxy(c *fiber.Ctx) error {
	URL := fmt.Sprintf("https://s.pximg.net/%s", c.Params("*"))
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(c.Context())

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
	URL := fmt.Sprintf("https://i.pximg.net/%s", c.Params("*"))
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(c.Context())
	req.Header.Add("Referer", "https://www.pixiv.net/")

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
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(c.Context())

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
