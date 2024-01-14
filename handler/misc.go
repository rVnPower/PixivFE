package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"codeberg.org/vnpower/pixivfe/models"
	"github.com/goccy/go-json"
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

func (p *PixivClient) GetNewestArtworks(worktype string, r18 string) ([]models.IllustShort, error) {
	var newWorks []models.IllustShort
	lastID := "0"

	for i := 0; i < 10; i++ {
		URL := fmt.Sprintf(ArtworkNewestURL, worktype, r18, lastID)

		response, err := p.PixivRequest(URL)
		if err != nil {
			return nil, err
		}

		var body struct {
			Illusts []models.IllustShort `json:"illusts"`
			LastID  string               `json:"lastId"`
		}

		err = json.Unmarshal([]byte(response), &body)
		if err != nil {
			return nil, err
		}
		newWorks = append(newWorks, body.Illusts...)

		lastID = body.LastID
	}

	return newWorks, nil
}

func (p *PixivClient) GetRanking(mode string, content string, date string, page string) (models.RankingResponse, error) {
	// Ranking data is formatted differently
	var pr models.RankingResponse

	if len(date) > 0 {
		date = "&date=" + date
	}

	url := fmt.Sprintf(ArtworkRankingURL, mode, content, date, page)

	s, err := p.TextRequest(url)

	if err != nil {
		return pr, err
	}

	err = json.Unmarshal([]byte(s), &pr)
	if err != nil {
		return pr, err
	}
	pr.PrevDate = strings.ReplaceAll(string(pr.PrevDateRaw[:]), "\"", "")
	pr.NextDate = strings.ReplaceAll(string(pr.NextDateRaw[:]), "\"", "")

	return pr, nil
}

func (p *PixivClient) GetSearch(artworkType string, name string, order string, age_settings string, page string) (*models.SearchResult, error) {
	URL := fmt.Sprintf(SearchArtworksURL, artworkType, name, order, age_settings, page)

	response, err := p.PixivRequest(URL)
	if err != nil {
		return nil, err
	}

	// IDK how to do better than this lol
	temp := strings.ReplaceAll(string(response), `"illust"`, `"works"`)
	temp = strings.ReplaceAll(temp, `"manga"`, `"works"`)
	temp = strings.ReplaceAll(temp, `"illustManga"`, `"works"`)

	var resultRaw struct {
		*models.SearchResult
		ArtworksRaw json.RawMessage `json:"works"`
	}
	var artworks models.SearchArtworks
	var result *models.SearchResult

	err = json.Unmarshal([]byte(temp), &resultRaw)
	if err != nil {
		return nil, err
	}

	result = resultRaw.SearchResult

	err = json.Unmarshal([]byte(resultRaw.ArtworksRaw), &artworks)
	if err != nil {
		return nil, err
	}

	result.Artworks = artworks

	return result, nil
}

func (p *PixivClient) GetDiscoveryArtwork(mode string, count int) ([]models.IllustShort, error) {
	var artworks []models.IllustShort

	for count > 0 {
		itemsForRequest := min(100, count)

		count -= itemsForRequest

		URL := fmt.Sprintf(ArtworkDiscoveryURL, mode, itemsForRequest)

		response, err := p.PixivRequest(URL)
		if err != nil {
			return nil, err
		}

		var thumbnail struct {
			Data json.RawMessage `json:"thumbnails"`
		}

		err = json.Unmarshal([]byte(response), &thumbnail)
		if err != nil {
			return nil, err
		}

		var body struct {
			Artworks []models.IllustShort `json:"illust"`
		}

		err = json.Unmarshal([]byte(thumbnail.Data), &body)
		if err != nil {
			return nil, err
		}

		artworks = append(artworks, body.Artworks...)
	}

	return artworks, nil
}

func (p *PixivClient) GetRankingLog(mode string, year, month int, image_proxy string) (template.HTML, error) {
	url := fmt.Sprintf("https://www.pixiv.net/ranking_log.php?mode=%s&date=%d%02d", mode, year, month)
	resp, err := http.Get(url)
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
					links = append(links, models.ProxyImage(a.Val, image_proxy))
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
