package main

import (
	"html/template"
	"pixivfe/configs"
	"pixivfe/views"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	server := gin.Default()

	server.SetFuncMap(template.FuncMap{
		"inc": func(n int) int {
			// For rankings to increment a number by 1
			return n + 1
		},
		"add": func(a int, b int) int {
			return a + b
		},

		"dec": func(n int) int {
			return n - 1
		},

		"toInt": func(s string) int {
			n, _ := strconv.Atoi(s)
			return n
		},

		"proxyImage": func(url string) string {
			if strings.Contains(url, "s.pximg.net") {
				// This subdomain didn't get proxied
				return url
			}

			regex := regexp.MustCompile(`.*?pximg\.net`)
			proxy := "https://" + configs.Configs.ImageProxyServer

			return regex.ReplaceAllString(url, proxy)
		},

		"isEmpty": func(s string) bool {
			return len(s) < 1
		},

		"isEmphasize": func(s string) bool {
			switch s {
			case
				"R-18",
				"R-18G":
				return true
			}
			return false
		},
	})

	// Static files
	server.StaticFile("/favicon.ico", "./template/favicon.ico")
	server.Static("css/", "./template/css")
	server.Static("assets/", "./template/assets")

	// HTML templates, automatically loaded
	server.LoadHTMLGlob("template/*.html")

	// Routes/Views
	views.SetupRoutes(server)

	// Disable trusted proxies since we do not use any for now
	server.SetTrustedProxies(nil)

	return server
}

func main() {
	configs.Configs.ReadConfig()

	r := setupRouter()

	r.Run(":8080")
}
