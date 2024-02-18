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

func pixivPostRequest(url, payload, token, csrf string) error {
	requestBody := []byte(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
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
	body_s := string(body)
	if !gjson.Valid(body_s) {
		return fmt.Errorf("invalid json: %v", body_s)
	}
	errr := gjson.Get(body_s, "error")

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
	csrf := session.GetCSRFToken(c)

	if token == "" || csrf == "" {
		return c.Redirect("/login")
	}

	id := c.Params("id")
	if id == "" {
		return errors.New("No ID provided.")
	}

	URL := "https://www.pixiv.net/ajax/illusts/bookmarks/add"
	payload := fmt.Sprintf(`{
"illust_id": "%s",
"restrict": 0,
"comment": "",
"tags": []
}`, id)
	if err := pixivPostRequest(URL, payload, token, csrf); err != nil {
		return err
	}

	return c.SendString("Success")
}

func DeleteBookmarkRoute(c *fiber.Ctx) error {
	token := session.CheckToken(c)
	csrf := session.GetCSRFToken(c)

	if token == "" || csrf == "" {
		return c.Redirect("/login")
	}

	id := c.Params("id")
	if id == "" {
		return errors.New("No ID provided.")
	}

	// You can't unlike
	URL := "https://www.pixiv.net/ajax/illusts/bookmarks/delete"
	payload := fmt.Sprintf(`bookmark_id=%s`, id)
	if err := pixivPostRequest(URL, payload, token, csrf); err != nil {
		return err
	}

	return c.SendString("Success")
}

func LikeRoute(c *fiber.Ctx) error {
	token := session.CheckToken(c)
	csrf := session.GetCSRFToken(c)

	if token == "" || csrf == "" {
		return c.Redirect("/login")
	}

	id := c.Params("id")
	if id == "" {
		return errors.New("No ID provided.")
	}

	URL := "https://www.pixiv.net/ajax/illusts/like"
	payload := fmt.Sprintf(`{"illust_id": "%s"}`, id)
	if err := pixivPostRequest(URL, payload, token, csrf); err != nil {
		return err
	}

	return c.SendString("Success")
}
