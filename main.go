package main

import (
	"pixivfe/configs"
	"pixivfe/handler"
	"pixivfe/views"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func setupRouter() *fiber.App {
	// HTML templates, automatically loaded
	engine := html.New("./template", ".html")

	handler.GetTemplateFunctions(engine)

	server := fiber.New(fiber.Config{
		Views: engine,
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

	if err != nil {
		panic(err)
	}

	r := setupRouter()

	r.Listen(":" + configs.Port)
}
