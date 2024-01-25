package pages

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"

	session "codeberg.org/vnpower/pixivfe/v2/core/config"
	"github.com/gofiber/fiber/v2"
	"github.com/tidwall/gjson"
)

func pixivPostRequest(id, token, csrf string) error {
	URL := "https://www.pixiv.net/ajax/illusts/bookmarks/add"

	requestBody := []byte(fmt.Sprintf(`{
"illust_id": "%s",
"restrict": 0,
"comment": "",
"tags": []
}`, id))

	req, _ := http.NewRequest("POST", URL, bytes.NewBuffer(requestBody))
	req.Header.Add("User-Agent", "Mozilla/5.0")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Cookie", "PHPSESSID="+token)
	req.Header.Add("x-csrf-token", csrf)
	// req.AddCookie(&http.Cookie{
	// 	Name:  "PHPSESSID",
	// 	Value: token,
	// })

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.New("Failed to do this action.")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Cannot parse the response from Pixiv. Please report this issue.")
	}

	errr := gjson.Get(string(body), "error")

	if !errr.Exists() {
		return errors.New("Incompatible request body.")
	}

	if errr.Bool() {
		return errors.New("Pixiv: Invalid request.")
	}
	return nil
}

func AddBookmarkRoute(c *fiber.Ctx) error {
	token := session.CheckToken(c)
	if token == "" {
		return c.Redirect("/login")
	}

	csrf := session.GetCSRFToken(c)
	if csrf == "" {
		return c.Redirect("/login")
	}

	id := c.Params("id")
	if id == "" {
		return errors.New("No ID provided.")
	}

	if err := pixivPostRequest(id, token, csrf); err != nil {
		return err
	}

	return c.SendString("Success")
}
