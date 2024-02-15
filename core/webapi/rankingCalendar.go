package core

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	session "codeberg.org/vnpower/pixivfe/v2/core/config"
	url "codeberg.org/vnpower/pixivfe/v2/core/http"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/net/html"
)

func get_weekday(n time.Weekday) int {
	switch n {
	case time.Sunday:
		return 1
	case time.Monday:
		return 2
	case time.Tuesday:
		return 3
	case time.Wednesday:
		return 4
	case time.Thursday:
		return 5
	case time.Friday:
		return 6
	case time.Saturday:
		return 7
	}
	return 0
}

func GetRankingCalendar(c *fiber.Ctx, mode string, year, month int) (template.HTML, error) {
	token := session.GetToken(c)
	URL := url.GetRankingCalendarURL(mode, year, month)

	req, _ := http.NewRequest("GET", URL, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0")
	req.Header.Add("Cookie", "PHPSESSID="+token)
	// req.AddCookie(&http.Cookie{
	// 	Name:  "PHPSESSID",
	// 	Value: token,
	// })

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Use the html package to parse the response body from the request
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}

	// Find and print all links on the web page
	var links []string
	var link func(*html.Node)
	link = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "img" {
			for _, a := range n.Attr {
				if a.Key == "data-src" {
					// adds a new link entry when the attribute matches
					links = append(links, session.ProxyImageUrlNoEscape(c, a.Val))
				}
			}
		}

		// traverses the HTML of the webpage from the first child node
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			link(c)
		}
	}
	link(doc)

	// now := time.Now()
	// yearNow := now.Year()
	// monthNow := now.Month()
	lastMonth := time.Date(year, time.Month(month), 0, 0, 0, 0, 0, time.UTC)
	thisMonth := time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC)

	renderString := ""
	for i := 0; i < get_weekday(lastMonth.Weekday()); i++ {
		renderString += "<div class=\"calendar-node calendar-node-empty\"></div>"
	}
	for i := 0; i < thisMonth.Day(); i++ {
		date := fmt.Sprintf("%d%02d%02d", year, month, i+1)
		if len(links) > i {
			renderString += fmt.Sprintf(`<a href="/ranking?mode=%s&date=%s"><div class="calendar-node" style="background-image: url(%s)"><span>%d</span></div></a>`, mode, date, links[i], i+1)
		} else {
			renderString += fmt.Sprintf(`<div class="calendar-node"><span>%d</span></div>`, i+1)
		}
	}
	return template.HTML(renderString), nil
}
