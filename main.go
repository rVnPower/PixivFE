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
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/gofiber/template/jet/v2"
)

func setupRouter() *fiber.App {
	// HTML templates, automatically loaded
	engine := jet.New("./template", ".jet.html")

	engine.AddFuncMap(handler.GetTemplateFunctions())

	server := fiber.New(fiber.Config{
		AppName:               "PixivFE",
		DisableStartupMessage: false,
		Views:                 engine,
		Prefork:               true,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
		ViewsLayout:           "layout",
	})

	server.Use(logger.New())
	server.Use(cache.New(
		cache.Config{
			KeyGenerator: func(c *fiber.Ctx) string {
				return utils.CopyString(c.OriginalURL())
			},
		},
	))
	server.Use(helmet.New())
	server.Use(csrf.New(csrf.Config{
		KeyLookup:      "header:X-Csrf-Token", // string in the form of '<source>:<key>' that is used to extract token from the request
		CookieName:     "my_csrf_",            // name of the session cookie
		CookieSameSite: "Strict",              // indicates if CSRF cookie is requested by SameSite
		Expiration:     3 * time.Hour,         // expiration is the duration before CSRF token will expire
		KeyGenerator:   utils.UUID,            // creates a new CSRF token
	}))

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

	if strings.Contains(configs.Port, "/") {
		ln, err := net.Listen("unix", configs.Port)
		if err != nil {
			panic("Failed to listen to " + configs.Port)
		}
		r.Listener(ln)
	}
	r.Listen(":" + configs.Port)
}
