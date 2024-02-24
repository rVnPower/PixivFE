// User Settings (Using Browser Cookies)

package core

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type CookieName string

const (
	Cookie_Token      CookieName = "__Host-pixivfe-Token"
	Cookie_CSRF       CookieName = "__Host-pixivfe-CSRF"
	Cookie_ImageProxy CookieName = "__Host-pixivfe-ImageProxy"
)

func GetCookie(c *fiber.Ctx, name CookieName, defaultValue ...string) string {
	return c.Cookies(string(name), defaultValue...)
}

func SetCookie(c *fiber.Ctx, name CookieName, value string) {
	cookie := fiber.Cookie{
		Name:  string(name),
		Value: value,
		Path:    "/",
		// expires in 30 days from now
		Expires: c.Context().Time().Add(30 * (24 * time.Hour)),
		HTTPOnly: true,
        Secure: true,
		SameSite: fiber.CookieSameSiteStrictMode, // bye-bye cross site forgery
	}
	c.Cookie(&cookie)
}

func ClearCookie(c *fiber.Ctx, name CookieName) {
	c.ClearCookie(string(name))
}

func ClearAllCookies(c *fiber.Ctx) {
	c.ClearCookie()
}
