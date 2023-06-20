package configs

import "github.com/gofiber/fiber/v2/middleware/session"

var Store *session.Store

func SetupStorage() {
	Store = session.New()
	Store.RegisterType("")
}
