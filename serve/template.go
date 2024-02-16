package serve

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	core "codeberg.org/vnpower/pixivfe/v2/core/webapi"
	config "codeberg.org/vnpower/pixivfe/v2/core/config"
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

		return fmt.Sprintf(`<img src="/proxy/s.pximg.net/common/images/emoji/%s.png" alt="(%s)" class="emoji" />`, id, s)
	})
	return template.HTML(parsedString)
}

func ParsePixivRedirect(s string) template.HTML {
	regex := regexp.MustCompile(`\/jump\.php\?(http[^"]+)`)

	parsedString := regex.ReplaceAllStringFunc(s, func(s string) string {
		s = s[10:]
		return s
	})
	escaped, err := url.QueryUnescape(parsedString)
	if err != nil {
		return template.HTML(s)
	}
	return template.HTML(escaped)
}

func EscapeString(s string) string {
	escaped := url.QueryEscape(s)
	return escaped
}

func ParseTime(date time.Time) string {
	return date.Format("2006-01-02 15:04")
}

func CreatePaginator(base, ending string, current_page, max_page int) template.HTML {
	pageUrl := func(page int) string {
		return fmt.Sprintf(`%s%d%s`, base, page, ending)
	}

	const (
		peek  = 5          // this can be changed freely
		limit = peek*2 + 1 // tied to the algorithm below, do not change
	)
	hasMaxPage := max_page != -1
	count := 0
	pages := ""

	pages += `<div class="pagination-buttons">`
	{ // "jump to page" <form>
		hidden_section := ""
		urlParsed, err := url.Parse(base)
		if err != nil {
			panic(err)
		}
		for k, vs := range urlParsed.Query() {
			if k == "page" {
				continue
			}
			for _, v := range vs {
				hidden_section += fmt.Sprintf(`<input type="hidden" name="%s" value="%s"/>`, k, v)
			}
		}

		max_section := ""
		if hasMaxPage {
			max_section = fmt.Sprintf(`max="%d"`, max_page)
		}

		pages += fmt.Sprintf(`<form action="%s">%s<input name="page" type="number" required value="%d" min="%d" %s placeholder="Page№" title="Jump To Page Number"/></form>`, pageUrl(current_page), hidden_section, current_page, 1, max_section)
		pages += "<br />"
	}
	{
		// previous,first (two buttons)
		pages += `<span>`
		{
			pages += fmt.Sprintf(`<a href="%s" class="pagination-button">&laquo;</a>`, pageUrl(1))
			pages += fmt.Sprintf(`<a href="%s" class="pagination-button">&lsaquo;</a>`, pageUrl(max(1, current_page-1)))
		}
		pages += `</span>`

		for i := current_page - peek; (i <= max_page || max_page == -1) && count < limit; i++ {
			if i < 1 {
				continue
			}
			if i == current_page {
				pages += fmt.Sprintf(`<a href="%s" class="pagination-button" id="highlight">%d</a>`, pageUrl(i), i)
			} else {
				pages += fmt.Sprintf(`<a href="%s" class="pagination-button">%d</a>`, pageUrl(i), i)
			}
			count++
		}

		// next,last (two buttons)
		pages += `<span>`
		if hasMaxPage {
			pages += fmt.Sprintf(`<a href="%s" class="pagination-button">&rsaquo;</a>`, pageUrl(min(max_page, current_page+1)))
			pages += fmt.Sprintf(`<a href="%s" class="pagination-button">&raquo;</a>`, pageUrl(max_page))
		} else {
			pages += fmt.Sprintf(`<a href="%s" class="pagination-button">&rsaquo;</a>`, pageUrl(current_page+1))
			pages += fmt.Sprintf(`<a href="%s" class="pagination-button" class="disabled">&raquo;</a>`, pageUrl(max_page))
		}
		pages += `</span>`
	}
	pages += `</div>`

	return template.HTML(pages)
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

		"parsePixivRedirect": func(s string) template.HTML {
			return ParsePixivRedirect(s)
		},
		"escapeString": func(s string) string {
			return EscapeString(s)
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
		"parseTime": func(date time.Time) string {
			return ParseTime(date)
		},
		"createPaginator": func(base, ending string, current_page, max_page int) template.HTML {
			return CreatePaginator(base, ending, current_page, max_page)
		},
		"joinArtworkIds": func(artworks []core.ArtworkBrief) string {
			ids := []string{}
			for _, art := range artworks {
				ids = append(ids, art.ID)
			}
			return strings.Join(ids, ",")
		},
		"stripEmbed": func(s string) string {
			// this is stupid
			return s[:len(s)-6]
		},
		"getImageProxyProtocolAuthority": func() string {
			return "https://" + config.GlobalServerConfig.ProxyServerAuthority
		},
	}
}
