package core

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var Store *session.Store

func SetupStorage() {
	Store = session.New(session.Config{
		Expiration: time.Hour * 24 * 30,
	})
}

func GetImageProxy(c *fiber.Ctx) string {
	sess, err := Store.Get(c)
	if err != nil {
		log.Fatalln("Failed to get current session and its values! Falling back to server default!")
		return GlobalServerConfig.ProxyServer
	}
	value := sess.Get("ImageProxy")
	if value != nil {
		return value.(string)
	}

	return GlobalServerConfig.ProxyServer
}
