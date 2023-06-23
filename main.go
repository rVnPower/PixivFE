package main

import (
	"net"
	"pixivfe/configs"
	"pixivfe/handler"
	"pixivfe/views"
	"strings"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/limiter"
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
		DisableStartupMessage:   false,
		Views:                   engine,
		Prefork:                 false,
		JSONEncoder:             json.Marshal,
		JSONDecoder:             json.Unmarshal,
		ViewsLayout:             "layout",
		EnableTrustedProxyCheck: true,
		TrustedProxies:          []string{"0.0.0.0/0"},
		ProxyHeader:             fiber.HeaderXForwardedFor,
	})

	server.Use(logger.New())
	server.Use(cache.New(
		cache.Config{
			Expiration: 5 * time.Minute,

			KeyGenerator: func(c *fiber.Ctx) string {
				return utils.CopyString(c.OriginalURL())
			},
		},
	))
	server.Use(recover.New())

	server.Use(limiter.New(limiter.Config{
		Max:        30,
		Expiration: 5 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendString("Rate limited")
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
	// server.Use(func(c *fiber.Ctx) error {
	// 	sess, err := configs.Store.Get(c)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	var token_string, image_string string

	// 	token := sess.Get("token")
	// 	if token != nil {
	// 		token_string = token.(string)
	// 	} else {
	// 		token_string = configs.ProxyServer
	// 	}

	// 	image := sess.Get("image-proxy")
	// 	if image != nil {
	// 		image_string = image.(string)
	// 	}

	// 	c.Bind(fiber.Map{
	// 		"Token":      token_string,
	// 		"ImageProxy": image_string,
	// 	})
	// 	return c.Next()
	// })

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
	r.Listen(":" + configs.Port)
}
