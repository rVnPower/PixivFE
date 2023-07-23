package configs

import (
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"
)

var Store *session.Store

func SetupStorage() {
	Store = session.New(session.Config{
		Expiration: time.Hour * 24 * 30,
	})
	Store.RegisterType("")
}
