package views

import (
	"net/http"
	"time"

	"codeberg.org/vnpower/pixivfe/configs"
	"codeberg.org/vnpower/pixivfe/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

var PC *handler.PixivClient

// func not_found_page(c *fiber.Ctx) {
// 	return c.Render(http.StatusNotFound, "error.html", fiber.Map{
// 		"Title": "Not found",
// 		"Error": "Route " + c.Request.URL.Path + " not found.",
// 	})
// }

func NewPixivClient(timeout int) *handler.PixivClient {
	transport := &http.Transport{Proxy: http.ProxyFromEnvironment}
	client := &http.Client{
		Timeout:   time.Duration(timeout) * time.Millisecond,
		Transport: transport,
	}

	pc := &handler.PixivClient{
		Client: client,
		Header: make(map[string]string),
		Cookie: make(map[string]string),
		Lang:   "en",
	}

	return pc
}

func SetupRoutes(r *fiber.App) {
	PC = NewPixivClient(5000)
	PC.SetSessionID(configs.Token)
	PC.SetUserAgent(configs.UserAgent)
	PC.AddHeader("Accept-Language", configs.AcceptLanguage)

	limit := limiter.New(limiter.Config{
		Max: 10,
	})

	r.Get("/", index_page)
	r.Get("about", about_page)
	r.Get("artworks/:id/", limit, artwork_page).Name("artworks")
	r.Get("users/:id/:category?", user_page)
	r.Get("newest", newest_artworks_page)
	r.Get("ranking", ranking_page)
	r.Get("ranking_log", ranking_log_page)
	r.Get("tags/:name", search_page)
	r.Get("discovery", discovery_page)

	self := r.Group("self")
	self.Get("/", get_logged_in_user)
	self.Get("/following_works", following_works_page)
	self.Get("/bookmarks", your_bookmark_page)
	r.Get("login", login_page)
	r.Post("tags", search)

	settings := r.Group("settings")
	settings.Get("/", settings_page)
	settings.Post("/:type", settings_post)

	// Legacy illust URL
	r.Get("member_illust.php", func(c *fiber.Ctx) error {
		return c.Redirect("artworks/" + c.Query("illust_id"))
	})

	// 404 page
	// r.NoRoute(not_found_page)
}
