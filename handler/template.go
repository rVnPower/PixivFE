package handler

import (
	"fmt"
	"html/template"
	"math/rand"
	"regexp"
	"strconv"
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

func ParseEmojis(s string) template.HTML {
	emojiList := map[string]string{
		"normal":        "101",
		"surprise":      "102",
		"serious":       "103",
		"heaven":        "104",
		"happy":         "105",
		"excited":       "106",
		"sing":          "107",
		"cry":           "108",
		"normal2":       "201",
		"shame2":        "202",
		"love2":         "203",
		"interesting2":  "204",
		"blush2":        "205",
		"fire2":         "206",
		"angry2":        "207",
		"shine2":        "208",
		"panic2":        "209",
		"normal3":       "301",
		"satisfaction3": "302",
		"surprise3":     "303",
		"smile3":        "304",
		"shock3":        "305",
		"gaze3":         "306",
		"wink3":         "307",
		"happy3":        "308",
		"excited3":      "309",
		"love3":         "310",
		"normal4":       "401",
		"surprise4":     "402",
		"serious4":      "403",
		"love4":         "404",
		"shine4":        "405",
		"sweat4":        "406",
		"shame4":        "407",
		"sleep4":        "408",
		"heart":         "501",
		"teardrop":      "502",
		"star":          "503",
	}

	regex := regexp.MustCompile(`\(([^)]+)\)`)

	parsedString := regex.ReplaceAllStringFunc(s, func(s string) string {
		s = s[1 : len(s)-1] // Get the string inside
		id := emojiList[s]

		return fmt.Sprintf(`<img src="https://s.pximg.net/common/images/emoji/%s.png" alt="(%s)" class="emoji" />`, id, s)
	})
	return template.HTML(parsedString)
}

func GetTemplateFunctions() template.FuncMap {
	return template.FuncMap{
		"toInt": func(s string) int {
			n, _ := strconv.Atoi(s)
			return n
		},

		"parseEmojis": func(s string) template.HTML {
			return ParseEmojis(s)
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
		"reformatDate": func(s string) string {
			if len(s) != 8 {
				return s
			}
			return fmt.Sprintf("%s-%s-%s", s[4:], s[2:4], s[:2])
		},
	}
}
