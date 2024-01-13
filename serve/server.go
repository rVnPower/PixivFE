package serve

import (
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"time"
)

func ServerSetup() {
	server := fiber.New(fiber.Config{
		AppName:                 "PixivFE",
		DisableStartupMessage:   true,
		JSONEncoder:             json.Marshal,
		JSONDecoder:             json.Unmarshal,
		EnableTrustedProxyCheck: true,
		TrustedProxies:          []string{"0.0.0.0/0"},
		ProxyHeader:             fiber.HeaderXForwardedFor,
	})

	server.Use(logger.New(logger.Config{
		Format: "[${time} | ${ip}] ${status} ${path}",
	}))

	server.Use(cache.New(
		cache.Config{
			Expiration:   5 * time.Minute,
			CacheControl: true,

			// KeyGenerator: func(c *fiber.Ctx) string {
			// 	return utils.CopyString(c.OriginalURL())
			// },
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
}
