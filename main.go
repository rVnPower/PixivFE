package main

import (
	"fmt"
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
		"parseEmojis": func(s string) template.HTML {
			regex := regexp.MustCompile(`\(([^)]+)\)`)

			parsedString := regex.ReplaceAllStringFunc(s, func(s string) string {
				s = s[1 : len(s)-1] // Get the string inside
				var id string

				switch s {
				case "normal":
					id = "101"
				case "surprise":
					id = "102"
				case "serious":
					id = "103"
				case "heaven":
					id = "104"
				case "happy":
					id = "105"
				case "excited":
					id = "106"
				case "sing":
					id = "107"
				case "cry":
					id = "108"
				case "normal2":
					id = "201"
				case "shame2":
					id = "202"
				case "love2":
					id = "203"
				case "interesting2":
					id = "204"
				case "blush2":
					id = "205"
				case "fire2":
					id = "206"
				case "angry2":
					id = "207"
				case "shine2":
					id = "208"
				case "panic2":
					id = "209"
				case "normal3":
					id = "301"
				case "satisfaction3":
					id = "302"
				case "surprise3":
					id = "303"
				case "smile3":
					id = "304"
				case "shock3":
					id = "305"
				case "gaze3":
					id = "306"
				case "wink3":
					id = "307"
				case "happy3":
					id = "308"
				case "excited3":
					id = "309"
				case "love3":
					id = "310"
				case "normal4":
					id = "401"
				case "surprise4":
					id = "402"
				case "serious4":
					id = "403"
				case "love4":
					id = "404"
				case "shine4":
					id = "405"
				case "sweat4":
					id = "406"
				case "shame4":
					id = "407"
				case "sleep4":
					id = "408"
				case "heart":
					id = "501"
				case "teardrop":
					id = "502"
				case "star":
					id = "503"
				}
				return fmt.Sprintf(`<img src="https://s.pximg.net/common/images/emoji/%s.png" alt="(%s)" class="emoji" />`, id, s)
			})
			return template.HTML(parsedString)
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
