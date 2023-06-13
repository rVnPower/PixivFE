package main

import (
	"pixivfe/configs"
	"pixivfe/handler"
	"pixivfe/views"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/gofiber/template/jet/v2"
)

func setupRouter() *fiber.App {
	// HTML templates, automatically loaded
	engine := jet.New("./template", ".jet.html")

	handler.GetTemplateFunctions(engine)

	server := fiber.New(fiber.Config{
		Views:       engine,
		Prefork:     true,
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	server.Use(logger.New())
	server.Use(cache.New(
		cache.Config{
			KeyGenerator: func(c *fiber.Ctx) string {
				return utils.CopyString(c.OriginalURL())
			},
		},
	))

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

	if err != nil {
		panic(err)
	}

	r := setupRouter()

	r.Listen(":" + configs.Port)
}
