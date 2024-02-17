package core

import (
	"log"
	"math/rand"
	"net/url"
	"strings"
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

func saveSession(sess *session.Session) error {
	if err := sess.Save(); err != nil {
		return err
	}

	return nil
}

func ProxyImageUrl(c *fiber.Ctx, s string) string {
	proxyOrigin := GetImageProxyPrefix(c)
	s = strings.ReplaceAll(s, `https:\/\/i.pximg.net`, proxyOrigin)
	// s = strings.ReplaceAll(s, `https:\/\/i.pximg.net`, "/proxy/i.pximg.net")
	s = strings.ReplaceAll(s, `https:\/\/s.pximg.net`, "/proxy/s.pximg.net")
	return s
}

func ProxyImageUrlNoEscape(c *fiber.Ctx, s string) string {
	proxyOrigin := GetImageProxyPrefix(c)
	s = strings.ReplaceAll(s, `https://i.pximg.net`, proxyOrigin)
	// s = strings.ReplaceAll(s, `https:\/\/i.pximg.net`, "/proxy/i.pximg.net")
	s = strings.ReplaceAll(s, `https://s.pximg.net`, "/proxy/s.pximg.net")
	return s
}

// footgun: if proxy server has prefix path (.Path != "/""), PixivFE will not work
func GetImageProxyOrigin(c *fiber.Ctx) string {
	url := GetImageProxy(c)
	return url.Scheme + "://" + url.Host
}

func GetImageProxyPrefix(c *fiber.Ctx) string {
	url := GetImageProxy(c)
	return url.Scheme + "://" + url.Host + url.Path
	// note: not sure if url.EscapedPath() is useful here. go's standard library is trash at handling URL (:// should be part of the scheme)
}

func GetImageProxy(c *fiber.Ctx) url.URL {
	sess, err := Store.Get(c)
	if err != nil {
		log.Println("Failed to get current session!")
		// fall through to default case
	} else {
		value := sess.Get("ImageProxy")
		if value_s, ok := value.(string); ok {
			proxyUrl, err := url.Parse(value_s)
			if err != nil {
				// fall through to default case
			} else {
				return *proxyUrl
			}
		}
	}
	return GlobalServerConfig.ProxyServer
}

func GetRandomDefaultToken() string {
	defaultToken := GlobalServerConfig.Token[rand.Intn(len(GlobalServerConfig.Token))]

	return defaultToken
}

func GetToken(c *fiber.Ctx) string {
	defaultToken := GlobalServerConfig.Token[rand.Intn(len(GlobalServerConfig.Token))]

	sess, err := Store.Get(c)
	if err != nil {
		log.Fatalln("Failed to get current session and its values! Falling back to server default!")
		return defaultToken
	}
	value := sess.Get("Token")
	if value != nil {
		return value.(string)
	}

	return defaultToken
}

func CheckToken(c *fiber.Ctx) string {
	sess, err := Store.Get(c)
	if err != nil {
		log.Fatalln("Failed to get current session and its values!")
		return ""
	}
	value := sess.Get("Token")
	if value != nil {
		return value.(string)
	}

	return ""
}

func GetCSRFToken(c *fiber.Ctx) string {
	sess, err := Store.Get(c)
	if err != nil {
		log.Fatalln("Failed to get current session and its values!")
		return ""
	}
	value := sess.Get("CSRF")
	if value != nil {
		return value.(string)
	}

	return ""
}

func SetSessionValue(c *fiber.Ctx, name, value string) error {
	sess, err := Store.Get(c)
	if err != nil {
		return err
	}

	sess.Set(name, value)

	if err = saveSession(sess); err != nil {
		log.Fatalln("Failed to save session storage!")
		return err
	}

	return nil
}

func RemoveSessionValue(c *fiber.Ctx, name string) error {
	sess, err := Store.Get(c)
	if err != nil {
		return err
	}

	sess.Delete(name)

	if err = saveSession(sess); err != nil {
		log.Fatalln("Failed to save session storage!")
		return err
	}

	return nil
}
