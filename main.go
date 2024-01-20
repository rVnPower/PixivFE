package main

import (
	"log"
	"net"
	"os"
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
		},
	))

	server.Use(cache.New(
		cache.Config{
			// Next: func(c *fiber.Ctx) bool {
			// 	// Disable cache for settings page
			// 	return strings.Contains(c.Path(), "/settings") || c.Path() == "/"
			// },
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

	// Global headers (from GotHub)
	server.Use(func(c *fiber.Ctx) error {
		c.Set("X-Frame-Options", "SAMEORIGIN")
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("Referrer-Policy", "no-referrer")
		c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")

		return c.Next()
	})

	server.Static("/favicon.ico", "./views/assets/favicon.ico")
	server.Static("css/", "./views/css")
	server.Static("assets/", "./views/assets")
	server.Static("/robots.txt", "./views/assets/robots.txt")

	// Routes

	server.Get("about", pages.AboutPage)
	server.Get("newest", pages.NewestPage)
	server.Get("discovery", pages.DiscoveryPage)
	server.Get("ranking", pages.RankingPage)
	server.Get("rankingCalendar", pages.RankingCalendarPage)
	server.Get("users/:id/:category?", pages.UserPage)

	limit := limiter.New(limiter.Config{
		Max: 10,
	})
	server.Get("artworks/:id/", limit, pages.ArtworkPage).Name("artworks")

	// Settings group
	server.Get("login", pages.LoginPage)
	settings := server.Group("settings")
	settings.Get("/", pages.SettingsPage)
	settings.Post("/:type", pages.SettingsPost)

	// Personal group
	self := server.Group("self")
	self.Get("/", pages.LoginUserPage)
	self.Get("/followingWorks", pages.FollowingWorksPage)
	self.Get("/bookmarks", pages.LoginBookmarkPage)

	// Listen
	if config.GlobalServerConfig.UnixSocket != "" {
		ln, err := net.Listen("unix", config.GlobalServerConfig.UnixSocket)
		if err != nil {
			log.Fatalf("Failed to run on Unix socket. %s", err)
			os.Exit(1)
		}
		log.Println("PixivFE is running on the configured Unix socket.")
		server.Listener(ln)
	} else {
		log.Println("PixivFE is running on the configured port.")
		server.Listen(":" + config.GlobalServerConfig.Port)
	}
}
