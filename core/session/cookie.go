// User Settings (Using Browser Cookies)

package core

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type CookieName string

const ( // the __Host thing force it to be secure and same-origin (no subdomain) >> https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Set-Cookie
	Cookie_Token         CookieName = "__Host-pixivfe-Token"
	Cookie_CSRF          CookieName = "__Host-pixivfe-CSRF"
	Cookie_ImageProxy    CookieName = "__Host-pixivfe-ImageProxy"
	Cookie_NovelFontType CookieName = "__Host-pixivfe-NovelFontType"
	Cookie_ShowArtR18    CookieName = "__Host-pixivfe-ShowArtR18"
	Cookie_ShowArtR18G   CookieName = "__Host-pixivfe-ShowArtR18G"
	Cookie_ShowArtAI     CookieName = "__Host-pixivfe-ShowArtAI"
)

// Go can't make this a const...
var AllCookieNames []CookieName = []CookieName{
	Cookie_Token,
	Cookie_CSRF,
	Cookie_ImageProxy,
	Cookie_NovelFontType,
}

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

var CookieExpireDelete = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

func ClearCookie(c *fiber.Ctx, name CookieName) {
	cookie := fiber.Cookie{
		Name:  string(name),
		Value: "",
		Path:  "/",
		// expires in 30 days from now
		Expires:  CookieExpireDelete,
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteStrictMode,
	}
	c.Cookie(&cookie)
}

func ClearAllCookies(c *fiber.Ctx) {
	for _, name := range AllCookieNames {
		ClearCookie(c, name)
	}
}
