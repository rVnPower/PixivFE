package main

import (
	"net"
	"strings"
	"time"

	"codeberg.org/vnpower/pixivfe/configs"
	"codeberg.org/vnpower/pixivfe/handler"
	"codeberg.org/vnpower/pixivfe/views"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/gofiber/template/jet/v2"
)

func setup_router() *fiber.App {
	// HTML templates, automatically loaded
	engine := jet.New("./template", ".jet.html")

	engine.AddFuncMap(handler.GetTemplateFunctions())

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

	server.Use(logger.New())
	server.Use(cache.New(
		cache.Config{
			Next: func(c *fiber.Ctx) bool {
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
		var baseURL string
		if configs.BaseURL != "localhost" {
			baseURL = "https://" + configs.BaseURL
		}
		c.Bind(fiber.Map{"FullURL": baseURL + c.OriginalURL(), "BaseURL": baseURL})
		return c.Next()
	})

	// Static files
	server.Static("/favicon.ico", "./template/favicon.ico")
	server.Static("css/", "./template/css")
	server.Static("assets/", "./template/assets")

	// Routes/Views
	views.SetupRoutes(server)

	// Disable trusted proxies since we do not use any for now
	// server.SetTrustedProxies(nil)

	return server
}

func main() {
	err := configs.ParseConfig()
	configs.SetupStorage()

	if err != nil {
		panic(err)
	}

	r := setup_router()

	if strings.Contains(configs.Port, "/") {
		ln, err := net.Listen("unix", configs.Port)
		if err != nil {
			panic("Failed to listen to " + configs.Port)
		}
		r.Listener(ln)
	}
	println("PixivFE is up and running on port " + configs.Port + "!")
	r.Listen(":" + configs.Port)
}
