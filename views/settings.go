package views

import (
	"net/http"
	"pixivfe/configs"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func get_storage(c *fiber.Ctx) *session.Session {
	// Get the storage
	sess, err := configs.Store.Get(c)
	if err != nil {
		panic(err)
	}
	return sess
}

func save_storage(sess *session.Session) {
	if err := sess.Save(); err != nil {
		panic(err)
	}
}

func test_request(URL string) error {
	req, _ := http.NewRequest("GET", URL, nil)

	transport := &http.Transport{Proxy: http.ProxyFromEnvironment}
	client := &http.Client{
		Timeout:   time.Duration(5000) * time.Millisecond,
		Transport: transport,
	}

	// Add headers
	// for k, v := range p.Header {
	// 	req.Header.Add(k, v)
	// }
	// for k, v := range p.Cookie {
	// }
	// Make a request
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return err
	}

	return nil
}

func set_token(c *fiber.Ctx) string {
	name := "token"
	sess := get_storage(c)

	// Parse the value from the form
	token := c.FormValue(name)
	if token != "" {
		client := NewPixivClient(3000)
		client.SetSessionID(token)
		client.SetUserAgent(configs.UserAgent)

		// If token is valid, save it
		if _, err := client.PixivRequest("https://www.pixiv.net/ajax/top/illust?mode=r18"); err != nil {
			return "Invalid token"
		}
		sess.Set(name, token)

		save_storage(sess)

		return ""
	}
	return "Empty form"
}

func set_image_server(c *fiber.Ctx) string {
	name := "image-proxy"
	sess := get_storage(c)

	// Parse the value from the form
	token := c.FormValue(name)
	if token != "" {
		if err := test_request("https://" + token + "/img-original/img/2023/06/10/19/51/55/108894158_p0.jpg"); err != nil {
			return "Cannot get image from " + token
		}
		// If token is valid, save it
		sess.Set(name, token)

		save_storage(sess)

		return ""
	}
	return "Empty form"
}
