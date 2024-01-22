package pages

import (
	"fmt"
	"io"
	"net/http"

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

	c.Set("Content-Type", "image/jpg, image/png, image/jpeg")

	return c.Send([]byte(body))
}
