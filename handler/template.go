package handler

import (
	"fmt"
	"html/template"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

func GetRandomColor() string {
	// Some color shade I stole
	colors := []string{
		// Green
		"#C8847E",
		"#C8A87E",
		"#C8B87E",
		"#C8C67E",
		"#C7C87E",
		"#C2C87E",
		"#BDC87E",
		"#82C87E",
		"#82C87E",
		"#7EC8AF",
		"#7EAEC8",
		"#7EA6C8",
		"#7E99C8",
		"#7E87C8",
		"#897EC8",
		"#967EC8",
		"#AE7EC8",
		"#B57EC8",
		"#C87EA5",
	}

	// Randomly choose one and return
	return colors[rand.Intn(len(colors))]
}

func GetTemplateFunctions() template.FuncMap {
	return template.FuncMap{
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

		"proxyImage": func(url string, target string) string {
			if strings.Contains(url, "s.pximg.net") {
				// This subdomain didn't get proxied
				return url
			}

			regex := regexp.MustCompile(`.*?pximg\.net`)
			proxy := "https://" + target

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

		"randomColor": func() string {
			return GetRandomColor()
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
	}
}
