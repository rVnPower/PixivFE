// User Settings (Using Browser Cookies)

package core

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type CookieName string

const ( // the __Host thing force it to be secure and same-origin (no subdomain) >> https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Set-Cookie
	Cookie_Token      CookieName = "__Host-pixivfe-Token"
	Cookie_CSRF       CookieName = "__Host-pixivfe-CSRF"
	Cookie_ImageProxy CookieName = "__Host-pixivfe-ImageProxy"
)

// Go can't make this a const...
var AllCookieNames = []CookieName{Cookie_Token, Cookie_CSRF, Cookie_ImageProxy}

func GetCookie(c *fiber.Ctx, name CookieName, defaultValue ...string) string {
	return c.Cookies(string(name), defaultValue...)
}

func SetCookie(c *fiber.Ctx, name CookieName, value string) {
	cookie := fiber.Cookie{
		Name:  string(name),
		Value: value,
		Path:  "/",
		// expires in 30 days from now
		Expires:  c.Context().Time().Add(30 * (24 * time.Hour)),
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteStrictMode, // bye-bye cross site forgery
	}
	c.Cookie(&cookie)
}

func ClearCookie(c *fiber.Ctx, name CookieName) {
	// c.ClearCookie(string(name)) // gofiber bug
	SetCookie(c, name, "")
}

func ClearAllCookies(c *fiber.Ctx) {
	// c.ClearCookie() // gofiber bug
	for _, name := range AllCookieNames {
		SetCookie(c, name, "")
	}
}
