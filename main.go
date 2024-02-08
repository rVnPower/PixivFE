package main

import (
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	config "codeberg.org/vnpower/pixivfe/v2/core/config"
	"codeberg.org/vnpower/pixivfe/v2/pages"
	"codeberg.org/vnpower/pixivfe/v2/serve"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/gofiber/template/jet/v2"
)

func main() {
	config.SetupStorage()
	config.GlobalServerConfig.InitializeConfig()

	engine := jet.New("./views", ".jet.html")

	engine.AddFuncMap(serve.GetTemplateFunctions())

	server := fiber.New(fiber.Config{
		AppName:                 "PixivFE",
		DisableStartupMessage:   true,
		Views:                   engine,
		Prefork:                 false,
		JSONEncoder:             json.Marshal,
		JSONDecoder:             json.Unmarshal,
		ViewsLayout:             "layout",
		EnableTrustedProxyCheck: true,
		TrustedProxies:          []string{"0.0.0.0/0"},
		ProxyHeader:             fiber.HeaderXForwardedFor,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Status code defaults to 500
			code := fiber.StatusInternalServerError

			// // Retrieve the custom status code if it's a *fiber.Error
			// var e *fiber.Error
			// if errors.As(err, &e) {
			// 	code = e.Code
			// }

			// Send custom error page
			err = c.Status(code).Render("pages/error", fiber.Map{"Title": "Error", "Error": err})
			if err != nil {
				// In case the SendFile fails
				return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
			}

			// Return from handler
			return nil
		},
	})

	server.Use(logger.New(
		logger.Config{
			Format: "${time} ${ip} | ${path}\n",
			Next: func(c *fiber.Ctx) bool {
				path := c.Path()
				return strings.Contains(path, "/assets") || strings.Contains(path, "/s.pximg.net") || strings.Contains(path, "/css")
			},
		},
	))

	server.Use(cache.New(
		cache.Config{
			Next: func(c *fiber.Ctx) bool {
				resp_code := c.Response().StatusCode()
				if resp_code < 200 || resp_code >= 300 {
					return true
				}

				// Disable cache for settings page
				return strings.Contains(c.Path(), "/settings") || c.Path() == "/"
			},
			Expiration:   5 * time.Minute,
			CacheControl: true,

			KeyGenerator: func(c *fiber.Ctx) string {
				return utils.CopyString(c.OriginalURL())
			},
		},
	))
	server.Use(recover.New())

	server.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	}))

	server.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			path := c.Path()
			return strings.Contains(path, "/assets") || strings.Contains(path, "/s.pximg.net") || strings.Contains(path, "/css")
		},
		Expiration:        30 * time.Second,
		Max:               config.GlobalServerConfig.RequestLimit,
		LimiterMiddleware: limiter.SlidingWindow{},
		LimitReached: func(c *fiber.Ctx) error {
			log.Println("Hit!")
			return errors.New("Woah! You are going too fast! I'll have to keep an eye on you.")
		},
	}))

	// Global headers (from GotHub)
	server.Use(func(c *fiber.Ctx) error {
		c.Set("X-Frame-Options", "SAMEORIGIN")
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("Referrer-Policy", "no-referrer")
		c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")

		return c.Next()
	})

	server.Use(func(c *fiber.Ctx) error {
		baseURL := c.BaseURL() + c.OriginalURL()
		c.Bind(fiber.Map{"BaseURL": baseURL})
		return c.Next()
	})

	server.Static("/favicon.ico", "./views/assets/favicon.ico")
	server.Static("/css/", "./views/css")
	server.Static("/assets/", "./views/assets")
	server.Static("/robots.txt", "./views/assets/robots.txt")

	// Routes

	server.Get("/", pages.IndexPage)
	server.Get("/about", pages.AboutPage)
	server.Get("/newest", pages.NewestPage)
	server.Get("/discovery", pages.DiscoveryPage)
	server.Get("/ranking", pages.RankingPage)
	server.Get("/rankingCalendar", pages.RankingCalendarPage)
	server.Post("/rankingCalendar", pages.RankingCalendarPicker)
	server.Get("/users/:id/:category?", pages.UserPage)
	server.Get("/artworks/:id/", pages.ArtworkPage).Name("artworks")

	// Settings group
	settings := server.Group("/settings")
	settings.Get("/", pages.SettingsPage)
	settings.Post("/:type", pages.SettingsPost)

	// Personal group
	self := server.Group("/self")
	self.Get("/", pages.LoginUserPage)
	self.Get("/followingWorks", pages.FollowingWorksPage)
	self.Get("/bookmarks", pages.LoginBookmarkPage)
	self.Post("/addBookmark/:id", pages.AddBookmarkRoute)
	self.Post("/deleteBookmark/:id", pages.DeleteBookmarkRoute)
	self.Post("/like/:id", pages.LikeRoute)

	server.Get("/tags/:name", pages.TagPage)
	server.Post("/tags",
		func(c *fiber.Ctx) error {
			name := c.FormValue("name")

			return c.Redirect("/tags/"+name, http.StatusFound)
		})

	// Legacy illust URL
	server.Get("/member_illust.php", func(c *fiber.Ctx) error {
		return c.Redirect("/artworks/" + c.Query("illust_id"))
	})

	// Proxy routes
	proxy := server.Group("/proxy")
	proxy.Get("/s.pximg.net/*", pages.SPximgProxy)
	proxy.Get("/ugoira.com/*", pages.UgoiraProxy)

	// Listen
	if config.GlobalServerConfig.UnixSocket != "" {
		ln, err := net.Listen("unix", config.GlobalServerConfig.UnixSocket)
		if err != nil {
			log.Fatalf("Failed to run on Unix socket. %s", err)
			os.Exit(1)
		}
		log.Printf("PixivFE is running on %v\n", config.GlobalServerConfig.UnixSocket)
		server.Listener(ln)
	} else {
		addr := config.GlobalServerConfig.Host + ":" + config.GlobalServerConfig.Port
		log.Printf("PixivFE is running on %v\n", addr)

		// note: string concatenation is very flaky
		server.Listen(addr)
	}
}
